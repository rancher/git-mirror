package main

import (
  "encoding/json"
  "io/ioutil"
  "net"
  "time"

  log "github.com/Sirupsen/logrus"
)

const (
  StatsLogPeriod = 5 * time.Second
  WhoisFilename = "/var/log/nginx/whois.json"
)

type Stats struct {
  ipRequests    map[string]int
  linesParsed   int
  linesSkipped  int
  whoisRequests int
  whoisRecords  []*Whois
  t             *time.Ticker
}

func NewStats() *Stats {
  return &Stats{
    ipRequests: make(map[string]int),
    whoisRecords: LoadWhoisRecords(),
  }
}

func LoadWhoisRecords() []*Whois {
  list := &WhoisList{}
  if data, err := ioutil.ReadFile(WhoisFilename); err == nil {
    if len(data) > 0 {
      if err = json.Unmarshal(data, list); err != nil {
        log.Warn(err)
      }
    }
  }
  return list.WhoisRecords  
}

func (s *Stats) AddWhoisRecord(record *Whois) {
  s.whoisRecords = append(s.whoisRecords, record)
}

func (s *Stats) FindWhoisRecord(targetIp string) *Whois {
  ip := net.ParseIP(targetIp)
  for _, whois := range s.whoisRecords {
    for _, netBlock := range whois.NetBlocks {
      _, ipnet, err := net.ParseCIDR(netBlock["startAddress"] + "/" + netBlock["cidrLength"])
      if err != nil {
        log.Warn(err)
        continue
      }
      if ipnet.Contains(ip) {
        return whois
      }
    }
  }
  return nil
}

func (s *Stats) SaveWhoisRecords() {
  var data []byte
  var err  error
  if data, err = json.Marshal(&WhoisList{s.whoisRecords}); err != nil {
    log.Fatal(err)
  }
  if err = ioutil.WriteFile(WhoisFilename, data, 0666); err != nil {
    log.Fatalf("couldn't write: %v", err)
  }
}


func (s *Stats) StartLogging() {
  s.t = time.NewTicker(StatsLogPeriod)
  go func() {
    for _ = range s.t.C {
      s.Log()
    }
  }()
}

func (s *Stats) StopLogging() {
  s.t.Stop()
  s.Log()
}

func (s *Stats) Log() {
  log.WithFields(log.Fields{
    "lines_parsed": s.linesParsed,
    "lines_skipped": s.linesSkipped,
    "lines_analyzed": s.linesParsed - s.linesSkipped,
    "whois_req": s.whoisRequests,
    "unique_ip": len(s.ipRequests),
  }).Info("stats")
}
