package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	etcdclient "github.com/coreos/etcd/client"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

type client struct {
	config     *config
	repos      []*repository
	pollTicker *time.Ticker
	kapi       etcdclient.KeysAPI
}

func (c *client) testKapi() error {
	testKey := fmt.Sprintf("/test%d", rand.Int63())
	if _, err := c.kapi.Set(context.Background(), testKey, "test", nil); err != nil {
		return err
	}
	_, err := c.kapi.Delete(context.Background(), testKey, nil)
	return err
}

func newClient(cfg *config) *client {
	c := &client{
		config: cfg,
	}

	if len(cfg.etcdEndpoints) > 0 {
		etcdConfig := etcdclient.Config{
			Endpoints:               cfg.etcdEndpoints,
			Transport:               etcdclient.DefaultTransport,
			HeaderTimeoutPerRequest: 3 * time.Second,
		}
		etcdClient, err := etcdclient.New(etcdConfig)
		if err == nil {
			c.kapi = etcdclient.NewKeysAPI(etcdClient)
			err = c.testKapi()
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, repoURL := range cfg.Repositories {
		repo := newRepository(repoURL, cfg.Dir)
		c.repos = append(c.repos, repo)
	}
	c.pollTicker = time.NewTicker(c.config.PollPeriod)
	return c
}

func main() {
	cfg := loadConfig()
	client := newClient(cfg)

	r := mux.NewRouter()
	if client.kapi != nil {
		log.Info("Starting HTTP server to receive GitHub events")
		go client.watchEvents()
		r.HandleFunc("/postreceive", client.PostReceiveHandler)
	}
	r.HandleFunc("/repos/{repo}/commits", client.RepoRefHandler)
	r.HandleFunc("/repos/{repo}/commits/{branch}", client.RepoBranchRefHandler)

	go client.poll()
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(cfg.GithubListenAddress, nil))
}

func (c *client) poll() {
	log.WithFields(log.Fields{"Period": c.config.PollPeriod}).Info("Starting poll ticker")
	for _ = range c.pollTicker.C {
		for _, repo := range c.repos {
			go repo.fetch("poll")
		}
	}
}

var repoKeyChroot = "/git-mirror"

// watchEvents watches for github push events published to etcd
func (c *client) watchEvents() {
	log.Info("Starting event watcher")
	w := c.kapi.Watcher(repoKeyChroot, &etcdclient.WatcherOptions{
		Recursive: true,
	})

	for {
		if resp, err := w.Next(context.Background()); err != nil {
			log.WithFields(log.Fields{
				"message": err.Error(),
			}).Warn("Error receiving watch event")
		} else {
			log.WithFields(log.Fields{
				"Action":        resp.Action,
				"Key":           resp.Node.Key,
				"CreatedIndex":  resp.Node.CreatedIndex,
				"ModifiedIndex": resp.Node.ModifiedIndex,
			}).Info("Received watch event")

			if resp.Action == "set" {
				repoName := strings.TrimPrefix(resp.Node.Key, repoKeyChroot+"/")
				if repo, err2 := c.getRepoByName(repoName); err2 != nil {
					log.WithFields(log.Fields{"Reason": "event", "Repo": repoName}).Error(err2)
				} else {
					repo.fetch("event")
				}
			}
		}
	}
}

// writeEvent writes the GH push event to etcd
func (c *client) writeEvent(e githubPushEvent) {
	key := fmt.Sprintf("%s/%s", repoKeyChroot, e.Repo.Name)
	val := fmt.Sprintf("%d", e.Repo.ID)

	log.WithFields(log.Fields{
		"repo": e.Repo.Name,
		"id":   e.Repo.ID,
	}).Debug("Writing event")

	if _, err := c.kapi.Set(context.Background(), key, val, nil); err != nil {
		log.WithFields(log.Fields{
			"repo": e.Repo.Name,
			"id":   e.Repo.ID,
		}).Error("Error writing event")
	}
}

func (c *client) getRepoByName(name string) (*repository, error) {
	for _, repo := range c.repos {
		if repo.name == name {
			return repo, nil
		}
	}
	return nil, errors.New("Repo not being mirrored")
}

func (c *client) RepoRefHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repo"]

	// defer log.WithFields(log.Fields{
	// 	"Repo": repoName,
	// }).Debugf("RepoRefHandler")

	if repo, err := c.getRepoByName(repoName); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, err.Error())
	} else {
		fmt.Fprintf(w, repo.getRefs())
	}
}

func (c *client) RepoBranchRefHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repo"]
	branch := vars["branch"]
	oldHash := strings.Trim(r.Header.Get("If-None-Match"), "\"")

	// defer log.WithFields(log.Fields{
	// 	"Branch": branch,
	// 	"Hash":   oldHash,
	// 	"Repo":   repoName,
	// }).Debugf("RepoBranchRefHandler")

	if repo, err := c.getRepoByName(repoName); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, err.Error())
	} else if currentHash, exists := repo.getHeadRef(branch); !exists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Ref not found")
	} else if oldHash == currentHash {
		w.WriteHeader(http.StatusNotModified)
	} else {
		fmt.Fprintf(w, currentHash)
	}
}

func (c *client) PostReceiveHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("PostReceiveHandler")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	event, err := parsePushEvent(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		fmt.Fprintf(w, err.Error())
	} else if event.Repo.Name == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not found")
	} else {
		log.WithFields(log.Fields{
			"repo": event.Repo.Name,
			"id":   event.Repo.ID,
		}).Info("Received GitHub event")
		go c.writeEvent(event)
		fmt.Fprintf(w, "OK")
	}
}
