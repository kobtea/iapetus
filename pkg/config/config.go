package config

import (
	"github.com/kobtea/iapetus/pkg/model"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Clusters []Cluster `yaml:"clusters"`
}

type Cluster struct {
	Name  string `yaml:"name"`
	Nodes []Node `yaml:"nodes"`
	Rules []Rule `yaml:"rules"`
}

type Node struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

type Rule struct {
	Target  string                 `yaml:"target"`
	Default bool                   `yaml:"default"`
	Range   model.DurationCriteria `yaml:"range"`
	Time    model.TimeCriteria     `yaml:"time"`
	Start   model.TimeCriteria     `yaml:"start"`
	End     model.TimeCriteria     `yaml:"end"`
}

func Parse(buf []byte) (*Config, error) {
	var d Config
	err := yaml.Unmarshal(buf, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
