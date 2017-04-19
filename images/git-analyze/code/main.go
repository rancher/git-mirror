package main

import (
	"bufio"
	"flag"
	"os"
	"regexp"
	"sync"
	"time"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

const (
	WhoisRequestPerSecond = 128
	WhoisRequestPeriod    = time.Second / time.Duration(WhoisRequestPerSecond)
)

func main() {
	var file_path = flag.String("filepath", "/var/log/nginx/access.log", "Log files to analyze, wildcard allowed between quotes.")
	var period = flag.String("period", "24h", "period of time (past to now) to analyze")
	flag.Parse()

	files , err := filepath.Glob(*file_path)
	if err != nil {
		log.Fatal(err)
	}

	//        log_format main '[$time_local] $http_host $remote_addr $http_x_forwarded_for '
	//                        '"$request" $status $body_bytes_sent "$http_referer" '
	//                        '"$http_user_agent" $request_time $upstream_response_time';
	logline, err := regexp.Compile("^\\[([^\\]]+)\\] ([^ ]+) ([^ ]+) ([^ ]+) \"([^\"]*)\" ([^ ]+) ([^ ]+) \"([^\"]*)\" \"([^\"]*)\" ([^ ]+) ([^ ]+)")
	if err != nil {
		log.Fatal(err)
	}
	// GET /rancher-catalog.git/info/refs?service=git-upload-pack HTTP/1.1
	uploadpack, err := regexp.Compile("^GET /(rancher)-catalog(.git)?/info/refs\\?service=git-upload-pack")
	if err != nil {
		log.Fatal(err)
	}

	stats := NewStats()
	stats.StartLogging()

	var cutoffTime time.Time
	if *period != "" {
		d, err := time.ParseDuration(*period)
		if err != nil {
			log.Warnf("Invalid duration: %s", *period)
		} else {
			cutoffTime = time.Now().Add(-d).Round(time.Second)
			log.Infof("Cutoff Time: %v", cutoffTime)
		}
	}

	// Getting stats for every file
	for _, f := range files {
		log.Info("Analyzing ", f)
		file, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}

		// 64Kb buffer should be big enough
		scanner := bufio.NewScanner(file)
		var wg sync.WaitGroup

		for scanner.Scan() {
			line := scanner.Text()

			submatches := logline.FindStringSubmatch(string(line))
			stats.linesParsed++

			if len(submatches) != 12 {
				log.Warn(string(line))
				continue
			}

			logTime, err := time.Parse("2/Jan/2006:15:04:05 -0700", submatches[1])
			if err != nil {
				log.Warn(err)
				stats.linesSkipped++
				continue
			}

			// ensure the log line is in the requested period
			if !cutoffTime.IsZero() && logTime.Before(cutoffTime) {
				stats.linesSkipped++
				continue
			}

			// ensure the request is a catalog git-upload-pack
			if !uploadpack.MatchString(submatches[5]) {
				stats.linesSkipped++
				continue
			}

			// do a whois lookup
			ip := submatches[4]
			if _, ok := stats.ipRequests[ip]; !ok {
				if w := stats.FindWhoisRecord(ip); w == nil {
					wg.Add(1)
					go func() {
						w, err := whois(ip, &wg)
						stats.whoisRequests++
						if err == nil {
							stats.AddWhoisRecord(w)
						}
					}()
					time.Sleep(WhoisRequestPeriod)
				}
			}
			stats.ipRequests[ip] = append(stats.ipRequests[ip], logTime)
		}
		file.Close()
		if err != nil {
			log.Fatal(err)
		}
		wg.Wait()
	}
	stats.SaveWhoisRecords()
	stats.StopLogging()
	stats.LogResult()
}
