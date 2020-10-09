package config

import (
	"io/ioutil"
	"os"
	"testing"
)

const jsonConfig = "{\"router\":[{\"service\":\"/auth\",\"loadbalacer\":\"round-robin\",\"upstream\":[\"http://localhost:8081\",\"http://localhost:8082\"]},{\"servie\":\"/retrival\",\"upstream\":[\"http://localhost:8084\",\"http://localhost:8085\"],\"loadbalacer\":\"round-robin\"}],\"proxy\":[{\"name\":\"backend\",\"port\":\"8000\",\"to\":[\"http://service1.ae\",\"http://service2.example.com\"],\"loadbalacer\":\"round-robin\"},{\"name\":\"frontend\",\"port\":\"9000\",\"to\":[\"https://jsonplaceholder.typicode.com\",\"http://example.com\"],\"loadbalacer\":\"round-robin\"}]}"
const yamlConfig = `---
router:
- service: "/auth"
  loadbalacer: round-robin
  upstream:
  - http://localhost:8081
  - http://localhost:8082
- servie: "/retrival"
  upstream:
  - http://localhost:8084
  - http://localhost:8085
  loadbalacer: round-robin
proxy:
- name: backend
  port: '8000'
  to:
  - http://service1.ae
  - http://service2.example.com
  loadbalacer: round-robin
- name: frontend
  port: '9000'
  to:
  - https://jsonplaceholder.typicode.com
  - http://example.com
  loadbalacer: round-robin
`

func TestGetConfig(t *testing.T) {
	err := writeConfigFiles()
	if err != nil {
		t.Fatal("unable to write temporary tests files")
	}

	if _, err := GetConfig("fake_config.json", false); err == nil {
		t.Error(err)
	}

	validJsonConfig, err := GetConfig("config.json", true)
	if err != nil {
		t.Error(err)
	}
	if validJsonConfig.ProxyServiceCount() != 2 {
		t.Fatal("proxy counts are not determined correctly")
	}
	if validJsonConfig.RouterServiceCount() != 2 {
		t.Fatal("router counts are not determined correctly")
	}
	if err := validJsonConfig.validateProxy(); err != nil {
		t.Error(err)
	}
	if err := validJsonConfig.validateRouter(); err != nil {
		t.Error(err)
	}

	validYamlConfig, err := GetConfig("config.yml", false)
	if err != nil {
		t.Error(err)
	}
	if validYamlConfig.ProxyServiceCount() != 2 {
		t.Fatal("proxy counts are not determined correctly")
	}
	if validYamlConfig.RouterServiceCount() != 2 {
		t.Fatal("router counts are not determined correctly")
	}
	if err := validYamlConfig.validateProxy(); err != nil {
		t.Error(err)
	}
	if err := validYamlConfig.validateRouter(); err != nil {
		t.Error(err)
	}

	err = deleteConfigFiles()
	if err != nil {
		t.Fatal("unable to remove temporary tests files")
	}
}

func writeConfigFiles() error {
	err := ioutil.WriteFile("config.json", []byte(jsonConfig), 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("config.yml", []byte(yamlConfig), 0644)
	if err != nil {
		return err
	}

	return nil
}

func deleteConfigFiles() error {
	err := os.Remove("config.json")
	if err != nil {
		return err
	}

	err = os.Remove("config.yml")
	if err != nil {
		return err
	}

	return nil
}
