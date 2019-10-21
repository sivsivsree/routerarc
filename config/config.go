package config

import (
	"encoding/json"
	"github.com/sivsivsree/routerarc/data"
	"io/ioutil"
)

type Configurations struct {
	Router []struct {
		Service     string   `json:"service,omitempty"`
		Loadbalacer string   `json:"loadbalacer"`
		Upstream    []string `json:"upstream"`
		Servie      string   `json:"servie,omitempty"`
	} `json:"router"`

	Proxy []data.Proxy `json:"proxy"`
}

func (config Configurations) validateProxy() error {
	return nil
}

func (config Configurations) ProxyServiceCount() int {
	return len(config.Proxy)

}

func GetConfig(filename string) (*Configurations, error) {

	plan, readErr := ioutil.ReadFile(filename)
	if readErr != nil {
		return nil, readErr
	}

	var config Configurations
	if err := json.Unmarshal(plan, &config); err != nil {
		return nil, err
	}

	if err := config.validateProxy(); err != nil {
		return nil, err
	}

	return &config, nil
}
