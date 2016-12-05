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
  WhoisRequestPerSecond = 256
  WhoisRequestPeriod = time.Second / time.Duration(WhoisRequestPerSecond)
)

func main() {
  var filepath = flag.String("filepath", "/var/log/nginx/access.log", "the log file to analyze")
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

    for i, submatch := range logline.FindStringSubmatch(string(line)) {
      switch i {
      case 0:
        continue
      case 4:
        // TODO: fetch from storage
        if _, ok := stats.ipRequests[submatch]; !ok {
          wg.Add(1)
          // TODO: persist to storage
          go whois(submatch, &wg)
          stats.whoisRequests++
          time.Sleep(WhoisRequestPeriod)
        }
        stats.ipRequests[submatch] += 1
      default:
      }
    }
    stats.linesParsed++
  }
  wg.Wait()
  stats.StopLogging()

  log.Infof("%d unique addresses", len(stats.ipRequests))
}
