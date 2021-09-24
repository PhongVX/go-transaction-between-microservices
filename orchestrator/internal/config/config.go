package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)


func ReadConfig(path string) (*Config, error) {
	var c *Config
	conf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(conf, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

