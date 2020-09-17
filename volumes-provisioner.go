package main

import (
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type dirConfig struct {
	uid  int
	gid  int
	mode int
	path string
}

func chmod() {
	if err := os.Chmod("some-filename", 0644); err != nil {
		log.Fatal(err)
	}
}

func chown() {
	if err := os.Chown("some-filename", 1000, 1000); err != nil {
		log.Fatal(err)
	}
}

func dirsProperties() []string {
	var provisionDirectories string = (os.Getenv("PROVISION_DIRECTORIES"))
	var dirsConfig []string = strings.Split(provisionDirectories, ";")
	return dirsConfig
}

func main() {
	for _, dir := range dirsProperties() {
		var properties []string = strings.Split(dir, ":")

		uid, _ := strconv.Atoi(properties[0])
		gid, _ := strconv.Atoi(properties[1])
		modeUint, _ := strconv.ParseUint(properties[2], 0, 32)
		mode := os.FileMode(modeUint)
		path := properties[3]

		log.Info("uid: ", uid, ", gid: ", gid, ", mode: ", mode, ", path: ", path)

		if err := os.Chmod(path, mode); err != nil {
			log.Fatal(err)
		}
		if err := os.Chown(path, uid, gid); err != nil {
			log.Fatal(err)
		}
	}

	// chmod()
	// chown()

	log.Info(dirsProperties())
}
