package main

import (
  "flag"
  "io/ioutil"
  "os"
  "strings"
  "time"

  log "github.com/Sirupsen/logrus"
  "gopkg.in/yaml.v2"
//  "github.com/coreos/etcd/client"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
  return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
  for _, v := range strings.Split(value, ",") {
    *i = append(*i, v)
  }
  return nil
}

type Config struct {
  configFile         string
  etcdEndpoints      arrayFlags

  Dir                 string
  GithubListenAddress string
  PollPeriod          time.Duration
  Repositories        arrayFlags
}

func FlagToEnv(prefix, name string) string {
  return prefix + "_" + strings.ToUpper(strings.Replace(name, "-", "_", -1))
}

func setFlagsFromEnv(fs *flag.FlagSet) error {
  var err error
  fs.VisitAll(func(f *flag.Flag) {
    key := FlagToEnv("MIRROR", f.Name)
    val := os.Getenv(key)
    if val != "" {
      err = fs.Set(f.Name, val)
      log.Infof("recognized and used environment variable %s=%s", key, val)
    }
  })
  return err
}

func LoadConfig() *Config {
  cfg := Config{}

  fs := flag.NewFlagSet("config", flag.ContinueOnError)
  fs.StringVar(&cfg.configFile, "config-file", "", "location of the YAML configuration file")
  fs.Var(&cfg.etcdEndpoints, "etcd-endpoints", "comma-delimited list of etcd client endpoints, may be specified multiple times")
  fs.StringVar(&cfg.Dir, "data-dir", "/data", "location to store git repositories")
  fs.StringVar(&cfg.GithubListenAddress, "github-listen-addr", ":4141", "Listen on this address for GitHub push events")
  fs.DurationVar(&cfg.PollPeriod, "poll-period", 5 * time.Minute, "Poll each git repository periodically")
  fs.Var(&cfg.Repositories, "repo", "comma-delimited list of git repos to mirror, may be specified multiple times")
  fs.Parse(os.Args[1:])

  setFlagsFromEnv(fs)

  if cfg.configFile != "" {
    file, err := os.Open(cfg.configFile)
    if err != nil {
      log.Fatal("Couldn't load config file: " + err.Error())
    }

    var data []byte
    if data, err = ioutil.ReadAll(file); err != nil {
      log.Fatal("Couldn't read config file: " + err.Error())
    }

    if err = yaml.Unmarshal(data, &cfg); err != nil {
      log.Fatal("Couldn't unmarshal config file: " + err.Error())
    }    
  }

  log.WithFields(log.Fields{
    "etcd-endpoints": cfg.etcdEndpoints,
    "data-dir":    cfg.Dir,
    "github-listen-addr": cfg.GithubListenAddress,
    "poll-period":  cfg.PollPeriod,
    "repo":     cfg.Repositories,
  }).Info("Loaded configuration")

  return &cfg
}
