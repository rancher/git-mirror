package main

import (
  "fmt"
  "net/http"
  "io/ioutil"

  log "github.com/Sirupsen/logrus"
)

const (
  gitStoragePath = "/var/git"
)

func main() {
  // TODO: configurable repos for mirroring

  // TODO: run a fetch periodically in case webhook fails

  // TODO: fetch mirror of each repo idempotently
  // git clone --mirror https://github.com/llparse/infra-catalog
  http.HandleFunc("/postreceive", postReceiveHandler)
  log.Fatal(http.ListenAndServe(":4141", nil))
}

func postReceiveHandler(w http.ResponseWriter, r *http.Request) {
  data, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf(string(data))
  // TODO: send message to all mirror hosts (LB will only tell 1 host about event)
  
  // TODO: get repo from JSON and run `git fetch -p origin`
  fmt.Fprintf(w, "OK")
}