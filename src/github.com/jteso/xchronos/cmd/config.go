package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jteso/xchronos/config"
)

func parseJobs(location string) []*config.Config {
	allConfFiles := []*config.Config{}
	var confFile *config.Config

	files, _ := ioutil.ReadDir(location)
	for _, f := range files {
		if f.IsDir() == false && strings.HasSuffix(f.Name(), ".batch") {
			content, err := ioutil.ReadFile(f.Name())
			if err != nil {
				break
			}
			confFile, _ = config.ParseConfig(string(content))
			allConfFiles = append(allConfFiles, confFile)
		}
	}
	return allConfFiles
}

func lookupXChronosConfig() *config.Config {
	locations := []string{".", "~/.xchronos/etc/xchronos.conf"}
	var conf *config.Config

	for _, loc := range locations {
		conf, _ = getXChronosConfig(loc)
		if conf != nil {
			return conf
		}
	}
	panic("No config file found")
}

func getXChronosConfig(path string) (*config.Config, error) {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".conf") {
			content, err := ioutil.ReadFile(f.Name())
			if err != nil {
				return nil, err
			}
			return config.ParseConfig(string(content))
		}
	}
	return nil, fmt.Errorf("No config file found")
}
