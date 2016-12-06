package main

import (
  "bufio"
  "flag"
  "io"
  "os"
  "regexp"
  "sync"
  "time"

  log "github.com/Sirupsen/logrus"
)

const (
  WhoisRequestPerSecond = 128
  WhoisRequestPeriod = time.Second / time.Duration(WhoisRequestPerSecond)
)

func main() {
  var filepath = flag.String("filepath", "/var/log/nginx/access.log", "the log file to analyze")
  var period = flag.String("period", "24h", "period of time (past to now) to analyze")
  flag.Parse()

  log.Info("Analyzing ", *filepath)
  file, err := os.Open(*filepath)
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

  var cutoffTime     time.Time
  if *period != "" {
    d, err := time.ParseDuration(*period)
    if err != nil {
      log.Warnf("Invalid duration: %s", *period)
    } else {
      cutoffTime = time.Now().Add(-d).Round(time.Second)
      log.Infof("Cutoff Time: %v", cutoffTime)
    }
  }

  reader := bufio.NewReader(file)
  var wg sync.WaitGroup

  for {
    line, isPrefix, err := reader.ReadLine()
    if err == io.EOF {
      break
    }
    if err != nil {
      log.Fatal(err)
      break
    }
    if isPrefix {
      log.Fatalf("Line (len=%d) too long for buffer", len(line))
      break
    }

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
  wg.Wait()
  stats.SaveWhoisRecords()
  stats.StopLogging()
  stats.LogResult()
}
