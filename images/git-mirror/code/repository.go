package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
)

type repository struct {
	sync.Mutex
	url       string
	name      string
	targetDir string
	refs      map[string]string
	reflist   string
}

func newRepository(url string, baseDir string) *repository {
	var r repository
	r.url = url
	r.name = r.nameFromURL()
	r.targetDir = filepath.Join(baseDir, r.name)
	r.refs = make(map[string]string)
	r.mirror()
	return &r
}

func (r *repository) mirror() {
	r.Lock()
	defer r.Unlock()
	defer r.parsePackedRefs()

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
	defer r.parsePackedRefs()

	log.WithFields(log.Fields{"Reason": reason, "Repo": r.name}).Debug("Fetching")
	cmd := exec.Command("git", "-C", r.targetDir, "fetch", "-p", "origin")

	if err := cmd.Run(); err != nil {
		log.Fatal("Error fetching origin: " + err.Error())
	}

	log.WithFields(log.Fields{"Reason": reason, "Repo": r.name}).Debug("Fetched")
}

func (r *repository) parsePackedRefs() {
	packedRefPath := strings.Join([]string{r.targetDir, "packed-refs"}, "/")

	var advance int
	var token []byte
	if data, err := ioutil.ReadFile(packedRefPath); err == nil {
		newRefs := make(map[string]string)
		var keys []string
		for {
			advance, token, err = bufio.ScanLines(data, false)
			if advance == 0 {
				break
			}
			data = data[advance:]
			if err != nil {
				log.Warn("Error scanning packed-refs: %+v", err)
				break
			}
			line := string(token)
			if line[0:1] == "#" {
				continue
			}
			ref := strings.Trim(line[40:], " ")
			hash := line[:40]
			newRefs[ref] = hash
			keys = append(keys, ref)
		}
		sort.Strings(keys)

		newReflist := ""
		for _, key := range keys {
			newReflist = newReflist + fmt.Sprintf("%s\t%s\n", newRefs[key], key)
		}

		r.refs = newRefs
		r.reflist = newReflist
	} else {
		log.Debugf("error scanning file: %s", err.Error())
	}

	log.WithFields(log.Fields{"Path": packedRefPath, "Repo": r.name}).Debug("Parsed packed-refs")
}

func (r *repository) getHeadRef(branch string) (string, bool) {
	r.Lock()
	defer r.Unlock()

	val, exists := r.refs[strings.Join([]string{"refs/heads", branch}, "/")]
	return val, exists
}

func (r *repository) getRefs() string {
	r.Lock()
	defer r.Unlock()

	return r.reflist
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
