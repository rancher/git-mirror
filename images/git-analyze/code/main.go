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

    ParseLine:
    for i, submatch := range logline.FindStringSubmatch(string(line)) {
      switch i {
      case 0:
        continue
      case 1:
        if !cutoffTime.IsZero() {
          logTime, err := time.Parse("2/Jan/2006:15:04:05 -0700", submatch)
          if err != nil {
            log.Warn(err)
          } else if logTime.Before(cutoffTime) {
            stats.linesSkipped++
            break ParseLine
          }
        }
      case 4:
        ip := submatch
        // TODO: fetch from storage
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
        stats.ipRequests[submatch] += 1
      default:
      }
    }
    stats.linesParsed++
  }
  wg.Wait()
  stats.SaveWhoisRecords()
  stats.StopLogging()

  log.Infof("%d unique addresses", len(stats.ipRequests))
}
