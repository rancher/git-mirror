package main

import (
  "errors"
  "fmt"
  "io/ioutil"
  "net/http"
  "time"

  log "github.com/Sirupsen/logrus"
)

type Client struct {
  config     *Config
  repos      []*Repository
  pollTicker *time.Ticker
}

func NewClient(config *Config) *Client {
  c := &Client{
    config: config,
  }
  for _, repoUrl := range config.Repositories {
    repo := NewRepository(repoUrl, config.Dir)
    c.repos = append(c.repos, repo)
  }
  c.pollTicker = time.NewTicker(c.config.PollPeriod)
  return c
}

func main() {
  //log.SetFormatter(&log.JSONFormatter{})

  config := LoadConfig()
  client := NewClient(config)

  go client.poll()

  http.Handle("/postreceive", client)
  log.Fatal(http.ListenAndServe(config.GithubListenAddress, nil))
}

func (c *Client) poll() {
  for _ = range c.pollTicker.C {
    for _, repo := range c.repos {
      go repo.Fetch("poll")
    }
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
  } else if event.Repo == nil || event.Repo.Name == "" {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "Not found")
  } else {
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
