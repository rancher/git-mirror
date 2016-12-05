package main

import (
  "time"

  log "github.com/Sirupsen/logrus"
)

const (
  StatsLogPeriod = 1 * time.Second
)

type Stats struct {
  ipRequests    map[string]int
  linesParsed   int
  whoisRequests int
  t             *time.Ticker
}

func NewStats() *Stats {
  return &Stats{
    ipRequests: make(map[string]int),
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
  log.Infof("lines: %d\twhois: %d\tunique_ip: %d", s.linesParsed, s.whoisRequests, len(s.ipRequests))
}
