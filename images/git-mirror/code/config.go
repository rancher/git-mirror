package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	for _, v := range strings.Split(value, ",") {
		*i = append(*i, strings.TrimSpace(v))
	}
	return nil
}

type config struct {
	configFile    string
	etcdEndpoints arrayFlags
	debug         bool

	Dir                 string
	GithubListenAddress string
	PollPeriod          time.Duration
	Repositories        arrayFlags
}

func flagToEnv(prefix, name string) string {
	return prefix + "_" + strings.ToUpper(strings.Replace(name, "-", "_", -1))
}

func setFlagsFromEnv(fs *flag.FlagSet) error {
	var err error
	fs.VisitAll(func(f *flag.Flag) {
		key := flagToEnv("MIRROR", f.Name)
		val := os.Getenv(key)
		if val != "" {
			err = fs.Set(f.Name, val)
			log.Infof("recognized and used environment variable %s=%s", key, val)
		}
	})
	return err
}

func loadConfig() *config {
	cfg := config{}

	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	fs.StringVar(&cfg.configFile, "config-file", "", "Location of the YAML configuration file")
	fs.Var(&cfg.etcdEndpoints, "etcd-endpoints", "Comma-delimited list of etcd client endpoints, may be specified multiple times")
	fs.StringVar(&cfg.Dir, "data-dir", "/data", "Location to store git repositories")
	fs.StringVar(&cfg.GithubListenAddress, "github-listen-addr", ":4141", "Listen on this address for GitHub push events")
	fs.DurationVar(&cfg.PollPeriod, "poll-period", 5*time.Minute, "Poll each git repository periodically")
	fs.Var(&cfg.Repositories, "repo", "Comma-delimited list of git repos to mirror, may be specified multiple times")
	fs.BoolVar(&cfg.debug, "debug", false, "Enable debug-level logging")
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
		"etcd-endpoints":     cfg.etcdEndpoints,
		"data-dir":           cfg.Dir,
		"github-listen-addr": cfg.GithubListenAddress,
		"poll-period":        cfg.PollPeriod,
		"repo":               cfg.Repositories,
		"debug":              cfg.debug,
	}).Info("Loaded configuration")

	if cfg.debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	//log.SetFormatter(&log.JSONFormatter{})

	return &cfg
}
