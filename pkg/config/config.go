package config

import (
	"errors"
	"fmt"
	"github.com/kobtea/iapetus/pkg/model"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
	Listen struct {
		Addr   string `yaml:"addr"`
		Prefix string `yaml:"prefix"`
	} `yaml:"listen"`
	Clusters []model.Cluster `yaml:"clusters"`
}

func Validate(c *Config) []error {
	var res []error
	if len(c.Clusters) == 0 {
		res = append(res, errors.New("config error: require `clusters` at least 1"))
	}
	for _, cluster := range c.Clusters {
		if len(cluster.Name) == 0 {
			res = append(res, errors.New("config error: require `name` at cluster"))
		}
		if len(cluster.Nodes) == 0 {
			res = append(res, fmt.Errorf("config error: require `node` at least 1 at %s", cluster.Name))
		}
		haveDefault := false
		for _, rule := range cluster.Rules {
			if rule.Default {
				haveDefault = true
			}
		}
		if !haveDefault {
			res = append(res, fmt.Errorf("config error: require default target at %s", cluster.Name))
		}
	}
	return res
}

func Parse(buf []byte) (*Config, error) {
	var d Config
	err := yaml.Unmarshal(buf, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
