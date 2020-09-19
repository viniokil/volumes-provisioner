package main

import (
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type dirProperties struct {
	uid  int
	gid  int
	mode os.FileMode
	path string
}

func (d dirProperties) chmod() {
	if err := os.Chmod(d.path, d.mode); err != nil {
		log.Error(err)
	} else {
		log.Info("Chmod dir ", d.path, " to mode ", d.mode)
	}
}

func (d dirProperties) chown() {
	if err := os.Chown(d.path, d.uid, d.gid); err != nil {
		log.Error(err)
	} else {
		log.Info("Chown dir ", d.path, " to uid ", d.uid, " and gid ", d.gid)
	}
}

func parseConfig() []dirProperties {
	var dirsConfig []dirProperties

	envVarVvalue, envVarExist := os.LookupEnv("PROVISION_DIRECTORIES")
	if !envVarExist {
		return nil
	}

	var provisionDirConfigs []string = strings.Split(envVarVvalue, ";")

	for _, dir := range provisionDirConfigs {
		var properties []string = strings.Split(dir, ":")

		log.Debug("Dir properties:", properties)

		uid, _ := strconv.Atoi(properties[0])
		gid, _ := strconv.Atoi(properties[1])
		modeUint, _ := strconv.ParseUint(properties[2], 8, 32)
		mode := os.FileMode(modeUint)
		path := properties[3]

		log.Debug("Parsed values", "uid: ", uid, ", gid: ", gid, ", mode: ", mode, ", path: ", path)

		var config = dirProperties{uid, gid, mode, path}

		dirsConfig = append(dirsConfig, config)
	}

	return dirsConfig
}

func main() {

	log.SetLevel(log.DebugLevel)
	// log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	// log.SetFormatter(&log.JSONFormatter{})

	for _, dir := range parseConfig() {
		dir.chmod()
		dir.chown()
	}
}
