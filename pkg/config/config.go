package config

import (
	"github.com/kobtea/iapetus/pkg/model"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Clusters []model.Cluster `yaml:"clusters"`
}

func Parse(buf []byte) (*Config, error) {
	var d Config
	err := yaml.Unmarshal(buf, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
