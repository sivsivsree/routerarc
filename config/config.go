package config

import (
	"encoding/json"
	"github.com/sivsivsree/routerarc/data"
	"io/ioutil"
)

type ConfigFile struct {
	Router []struct {
		Service     string   `json:"service,omitempty"`
		Loadbalacer string   `json:"loadbalacer"`
		Upstream    []string `json:"upstream"`
		Servie      string   `json:"servie,omitempty"`
	} `json:"router"`

	Proxy []data.Proxy `json:"proxy"`
}

func (config ConfigFile) validateProxy() error {

	return nil
}

func GetConfig(filename string) (*ConfigFile, error) {
	plan, _ := ioutil.ReadFile(filename)
	var config ConfigFile
	if err := json.Unmarshal(plan, &config); err != nil {
		return nil, err
	}

	if err := config.validateProxy(); err != nil {
		return nil, err
	}

	return &config, nil
}
