package main

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	StatsLogPeriod = 15 * time.Second
	WhoisFilename  = "/var/log/nginx/whois.json"
)

type Stats struct {
	ipRequests    map[string][]time.Time
	whoisRecords  []*Whois
	linesParsed   int
	linesSkipped  int
	whoisRequests int
	totalInstalls int
	t             *time.Ticker
}

func NewStats() *Stats {
	return &Stats{
		ipRequests:   make(map[string][]time.Time),
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
	var err error
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
			s.LogProgress()
		}
	}()
}

func (s *Stats) StopLogging() {
	s.t.Stop()
}

func (s *Stats) LogProgress() {
	log.WithFields(log.Fields{
		"lines_parsed":   s.linesParsed,
		"lines_skipped":  s.linesSkipped,
		"whois_requests": s.whoisRequests,
		"unique_ip":      len(s.ipRequests),
	}).Info("progress")
}

func (s *Stats) LogResult() {
	s.calculateResult()
	log.WithFields(log.Fields{
		"lines_analyzed": s.linesParsed - s.linesSkipped,
		"total_installs": s.totalInstalls,
		"unique_ip":      len(s.ipRequests),
	}).Info("result")
}

func (s *Stats) calculateResult() {
	for _, requestTimes := range s.ipRequests {
		// assume 1 install if we only have 1 sample point
		if len(requestTimes) <= 1 {
			s.totalInstalls += 1
			continue
		}

		// calculate window of time in which we received requests
		earliest := time.Now()
		var latest time.Time
		for _, requestTime := range requestTimes {
			if requestTime.Before(earliest) {
				earliest = requestTime
			}
			if requestTime.After(latest) {
				latest = requestTime
			}
		}
		window := latest.Sub(earliest)

		theoreticalTicks := round(window.Minutes() / 5.0)
		if theoreticalTicks == 0 {
			s.totalInstalls += 1
		} else {
			requestPerTick := float64(len(requestTimes)-1) / theoreticalTicks
			if requestPerTick < 0 {
				log.Warnf("%d %f %v", requestPerTick, theoreticalTicks, window)
				continue
			}
			s.totalInstalls += int(round(requestPerTick))
		}
	}
}

func round(value float64) float64 {
	return math.Floor(value + .5)
}
