package main

import (
  "io/ioutil"
  "os"
  "time"

  log "github.com/Sirupsen/logrus"
  "gopkg.in/yaml.v2"
)

type Config struct {
  StorageDir    string
  ServerAddress string
  PollInterval  time.Duration
  Repositories  []string
}

func LoadConfig(configPath string) *Config {
  file, err := os.Open(configPath)
  if err != nil {
    log.Fatal("Couldn't load config file: " + err.Error())
  }

  var data []byte
  if data, err = ioutil.ReadAll(file); err != nil {
    log.Fatal("Couldn't read config file: " + err.Error())
  }

  config := Config{}
  if err = yaml.Unmarshal(data, &config); err != nil {
    log.Fatal("Couldn't unmarshal config file: " + err.Error())
  }

  log.WithFields(log.Fields{
    "StorageDir":    config.StorageDir,
    "ServerAddress": config.ServerAddress,
    "PollInterval":  config.PollInterval,
    "RepoCount":     len(config.Repositories),
  }).Info("Loaded configuration")

  return &config
}
