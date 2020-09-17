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

func parseConfig() []dirProperties {
	var dirsConfig []dirProperties

	envVarVvalue, envVarExist := os.LookupEnv("PROVISION_DIRECTORIES")
	if envVarExist {
		var provisionDirConfigs []string = strings.Split(envVarVvalue, ";")
		for _, dir := range provisionDirConfigs {
			var properties []string = strings.Split(dir, ":")

			log.Debug("Dir properties:", properties)

			uid, _ := strconv.Atoi(properties[0])
			gid, _ := strconv.Atoi(properties[1])
			modeUint, _ := strconv.ParseUint(properties[2], 0, 32)
			mode := os.FileMode(modeUint)
			path := properties[3]

			log.Debug("Parsed values", "uid: ", uid, ", gid: ", gid, ", mode: ", mode, ", path: ", path)

			var config = dirProperties{uid, gid, mode, path}

			dirsConfig = append(dirsConfig, config)

		}
	}

	return dirsConfig
}

func main() {
	for _, dir := range parseConfig() {
		if err := os.Chmod(dir.path, dir.mode); err != nil {
			log.Error(err)
		} else {
			log.Info("Chmod dir ", dir.path, " to mode ", dir.mode)
		}

		if err := os.Chown(dir.path, dir.uid, dir.gid); err != nil {
			log.Error(err)
		} else {
			log.Info("Chown dir ", dir.path, " to mode ", dir.mode)
		}
	}
}
