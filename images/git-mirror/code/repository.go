package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
)

type repository struct {
	sync.Mutex
	url       string
	name      string
	targetDir string
}

func newRepository(url string, baseDir string) *repository {
	var r repository
	r.url = url
	r.name = r.nameFromURL()
	r.targetDir = filepath.Join(baseDir, r.name)
	r.mirror()
	return &r
}

func (r *repository) mirror() {
	r.Lock()
	defer r.Unlock()

	if pathExists(r.targetDir) {
		log.WithFields(log.Fields{"Repo": r.name}).Info("Already exists")
		return
	}

	log.WithFields(log.Fields{"Repo": r.name}).Info("Cloning")
	cmd := exec.Command("git", "clone", "--mirror", r.url, r.targetDir)

	if err := cmd.Run(); err != nil {
		log.Fatal("Error creating mirror: " + err.Error())
	}

	log.WithFields(log.Fields{"Repo": r.name}).Info("Cloned")
}

func (r *repository) fetch(reason string) {
	r.Lock()
	defer r.Unlock()

	log.WithFields(log.Fields{"Reason": reason, "Repo": r.name}).Debug("Fetching")
	cmd := exec.Command("git", "-C", r.targetDir, "fetch", "-p", "origin")

	if err := cmd.Run(); err != nil {
		log.Fatal("Error fetching origin: " + err.Error())
	}

	log.WithFields(log.Fields{"Reason": reason, "Repo": r.name}).Debug("Fetched")
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func (r *repository) nameFromURL() string {
	parts := strings.Split(r.url, "/")
	name := parts[len(parts)-1]
	return strings.TrimSuffix(name, ".git")
}
