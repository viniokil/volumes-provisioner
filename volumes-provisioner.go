package main

import (
	"fmt"
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

func parseDirConfig() []dirProperties {
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

// logParseFormat takes a string fortmat and returns the Logrus log formtat.
func logParseFormat(fortmat string) (log.Formatter, error) {
	switch strings.ToLower(fortmat) {
	case "json":
		return &log.JSONFormatter{}, nil
	case "logfmt":
		return &log.TextFormatter{DisableColors: true, FullTimestamp: true}, nil
	case "text":
		return &log.TextFormatter{FullTimestamp: true}, nil
	}

	var f log.Formatter
	return f, fmt.Errorf("not a valid logrus format: %q", fortmat)
}

// logConfig receives values from environment variables and configures logs
func logConfig() {
	level, envVarExist := os.LookupEnv("LOG_LEVEL")
	if !envVarExist {
		level = "info"
	}

	logLevel, err := log.ParseLevel(level)
	if err != nil {
		log.Error(err)
		logLevel = log.DebugLevel
	}
	log.SetLevel(logLevel)

	format, envVarExist := os.LookupEnv("LOG_FORMAT")
	if !envVarExist {
		format = "text"
	}
	logFormat, err := logParseFormat(format)
	if err != nil {
		log.Error(err)
		logFormat = &log.TextFormatter{FullTimestamp: true}
	}
	log.SetFormatter(logFormat)

}

func main() {

	logConfig()

	for _, dir := range parseDirConfig() {
		dir.chmod()
		dir.chown()
	}
}
