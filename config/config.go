package config

import (
	"encoding/json"
	"github.com/sivsivsree/routerarc/data"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Configurations struct {
	Router []data.Router `json:"router";yaml:"router"`
	Proxy  []data.Proxy  `json:"proxy";yaml:"proxy"`
}

func (config Configurations) validateProxy() error {
	return nil
}

func (config Configurations) ProxyServiceCount() int {
	return len(config.Proxy)

}

func (config Configurations) validateRouter() error {
	// Todo: make sure the static and upstream does not exist each other
	return nil
}

func (config Configurations) RouterServiceCount() int {
	return len(config.Router)

}

func GetConfig(filename string, jsonfile bool) (*Configurations, error) {

	plan, readErr := ioutil.ReadFile(filename)
	if readErr != nil {
		return nil, readErr
	}

	var config Configurations
	if jsonfile {
		if err := json.Unmarshal(plan, &config); err != nil {
			return nil, err
		}
	} else {
		if err := yaml.Unmarshal(plan, &config); err != nil {
			return nil, err
		}
	}

	if err := config.validateProxy(); err != nil {
		return nil, err
	}

	return &config, nil
}
