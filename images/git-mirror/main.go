package main

import (
  "flag"
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
    repo := NewRepository(repoUrl, config.StorageDir)
    c.repos = append(c.repos, repo)
  }
  c.pollTicker = time.NewTicker(c.config.PollInterval)
  return c
}

func main() {
  //log.SetFormatter(&log.JSONFormatter{})

  var configPath string
  flag.StringVar(&configPath, "config-file", "config.yaml", "location of the YAML configuration file")
  flag.Parse()

  config := LoadConfig(configPath)
  client := NewClient(config)

  go client.poll()

  // TODO: fetch mirror of each repo idempotently
  // git clone --mirror https://github.com/llparse/infra-catalog
  http.HandleFunc("/postreceive", postReceiveHandler)
  log.Fatal(http.ListenAndServe(":4141", nil))
}

func (c *Client) poll() {
  for _ = range c.pollTicker.C {
    for _, repo := range c.repos {
      go repo.Fetch()
    }
  }
}

func postReceiveHandler(w http.ResponseWriter, r *http.Request) {
  data, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Fatal(err)
  }

  log.Printf(string(data))
  // TODO: send message to all mirror hosts (LB will only tell 1 host about event)
  
  // TODO: get repo from JSON and run `git fetch -p origin`
  fmt.Fprintf(w, "OK")
}
