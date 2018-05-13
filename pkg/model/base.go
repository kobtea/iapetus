package model

import pc "github.com/prometheus/prometheus/config"

type Cluster struct {
	Name  string `yaml:"name"`
	Nodes []Node `yaml:"nodes"`
	Rules []Rule `yaml:"rules"`
}

type Node struct {
	Name     string              `yaml:"name"`
	Url      string              `yaml:"url"`
	Relabels []*pc.RelabelConfig `yaml:"relabels"`
}

type Rule struct {
	Target  string           `yaml:"target"`
	Default bool             `yaml:"default"`
	Range   DurationCriteria `yaml:"range"`
	Time    TimeCriteria     `yaml:"time"`
	Start   TimeCriteria     `yaml:"start"`
	End     TimeCriteria     `yaml:"end"`
}
