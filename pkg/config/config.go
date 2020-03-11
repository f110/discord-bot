package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	EnablePlugins []string `yaml:"enable_plugins"`

	StorageHost string `yaml:"-"`
	Bucket      string `yaml:"-"`
	BucketHost  string `yaml:"-"`
}

func ReadConfig(p string) (*Config, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	if err := yaml.Unmarshal(b, conf); err != nil {
		return nil, err
	}

	return conf, nil
}
