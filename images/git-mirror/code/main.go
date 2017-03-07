package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

type Client struct {
	config     *Config
	repos      []*Repository
	pollTicker *time.Ticker
	kapi       client.KeysAPI
}

func (c *Client) testKapi() error {
	testKey := fmt.Sprintf("/test%d", rand.Int63())
	if _, err := c.kapi.Set(context.Background(), testKey, "test", nil); err != nil {
		return err
	} else {
		_, err := c.kapi.Delete(context.Background(), testKey, nil)
		return err
	}
}

func NewClient(cfg *Config) *Client {
	c := &Client{
		config: cfg,
	}

	etcdConfig := client.Config{
		Endpoints:               cfg.etcdEndpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: 3 * time.Second,
	}
	etcdClient, err := client.New(etcdConfig)
	if err == nil {
		c.kapi = client.NewKeysAPI(etcdClient)
		err = c.testKapi()
	}
	if err != nil {
		log.Warn(errors.New(err.Error() + ". Mirror operating in poll-only mode."))
		c.kapi = nil
	}

	for _, repoUrl := range cfg.Repositories {
		repo := NewRepository(repoUrl, cfg.Dir)
		c.repos = append(c.repos, repo)
	}
	c.pollTicker = time.NewTicker(c.config.PollPeriod)
	return c
}

func main() {
	//log.SetFormatter(&log.JSONFormatter{})
	cfg := LoadConfig()
	client := NewClient(cfg)

	if client.kapi != nil {
		log.Info("Starting HTTP server to receive GitHub events")
		go func() {
			http.Handle("/postreceive", client)
			log.Fatal(http.ListenAndServe(cfg.GithubListenAddress, nil))
		}()
		go client.watchEvents()
	}

	client.poll()
}

func (c *Client) poll() {
	for _ = range c.pollTicker.C {
		for _, repo := range c.repos {
			go repo.Fetch("poll")
		}
	}
}

func (c *Client) watchEvents() {
	log.Info("Starting event watcher")

}

func (c *Client) writeEvent(e GHPushEvent) {
	key := fmt.Sprintf("/repo/%s", e.Repo.Name)
	val := fmt.Sprintf("%d", e.Repo.Id)
	if resp, err := c.kapi.Set(context.Background(), key, val, nil); err != nil {
		log.Warnf("Error writing event: %v", err)
	} else {
		log.Infof("Wrote event: %+v", resp)
	}
}

func (c *Client) GetRepoByName(name string) (*Repository, error) {
	for _, repo := range c.repos {
		if repo.name == name {
			return repo, nil
		}
	}
	return nil, errors.New("Repo not being mirrored")
}

func (c *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	event, err := ParsePushEvent(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		fmt.Fprintf(w, err.Error())
	} else if event.Repo.Name == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not found")
	} else {
		c.writeEvent(event)
		// TODO: send message to all mirror hosts (LB will only tell 1 host about event)
		repo, err := c.GetRepoByName(event.Repo.Name)
		if err != nil {
			log.WithFields(log.Fields{"Reason": "event", "Repo": event.Repo.Name}).Error(err)
			return
		}
		go repo.Fetch("event")
		fmt.Fprintf(w, "OK")
	}
}
