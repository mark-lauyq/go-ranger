package ranger_config

import (
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestNewRemoteConfigReader(t *testing.T) {
	configReader := GetConfigReader("http://www.google.com")

	if configReader.GetConfigPath() != "http://www.google.com" {
		t.Error("invalid url set to configReader")
	}
}

func TestNewLocalConfigReader(t *testing.T) {
	configReader := GetConfigReader("file://./test_config.yaml")

	if configReader.GetConfigPath() != "./test_config.yaml" {
		t.Error("invalid url set to configReader")
	}
}

func TestParseLocalConfig(t *testing.T) {
	configReader := GetConfigReader("file://./test_config.yaml")
	data, err := configReader.ReadConfig()

	config := &Config{}
	yaml.Unmarshal(data, config)

	if err != nil {
		t.Error("unable to read config")
	}

	if config.AppName != "test" {
		t.Error("config parsed incorrectly")
	}
}

func TestParseRemoteConfig(t *testing.T) {
	//t.Error("@todo TestParseLocalConfig")
}

func TestParseConfig_InvalidConfig(t *testing.T) {
	//t.Error("@todo TestParseConfig_InvalidConfig")
}
