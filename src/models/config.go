package models

import (
	"gopkg.in/yaml.v2"
	"os"

	"github.com/lubiedo/yav/src/utils"
)

type Config struct {
	Port    string `yaml:"port"`
	Addr    string `yaml:"addr"`
	Verbose bool   `yaml:"verbose"`

	UseHTTPS bool   `yaml:"use-https",omitempty`
	CertPath string `yaml:"cert-path,omitempty"`
	KeyPath  string `yaml:"key-path,omitempty"`

	LogFile string `yaml:"log-file,omitempty"`
	Log     *utils.Log

	TplVars     map[string]interface{} `yaml:"vars,omitempty"`
	TplErroPage string                 `yaml:"error-page,omitempty"`

	WriteTimeOut int `yaml:"write-timeout,omitempty"`
	IdleTimeOut  int `yaml:"idle-timeout,omitempty"`
}

func (c *Config) LoadConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		c.Log.Fatal("Error loading configuration from \"%s\"",
			path)
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		c.Log.Fatal("Parsing configuration from \"%s\" failed",
			path)
	}
}
