package relabel

import (
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"
	pl "github.com/prometheus/prometheus/relabel"
)

func Matchers2LabelSet(ms []*labels.Matcher) model.LabelSet {
	ls := make(model.LabelSet)
	for _, m := range ms {
		ls[model.LabelName(m.Name)] = model.LabelValue(m.Value)
	}
	return ls
}

func MargeLabelSet(ls model.LabelSet, dst []*labels.Matcher) {
	for _, m := range dst {
		if v, ok := ls[model.LabelName(m.Name)]; ok {
			m.Value = string(v)
		}
	}
}

func Process(query string, configs []*config.RelabelConfig) (string, error) {
	expr, err := promql.ParseExpr(query)
	if err != nil {
		return "", err
	}
	promql.Inspect(expr, func(node promql.Node, nodes []promql.Node) bool {
		switch n := node.(type) {
		case *promql.VectorSelector:
			ls := Matchers2LabelSet(n.LabelMatchers)
			ls2 := pl.Process(ls, configs...)
			MargeLabelSet(ls2, n.LabelMatchers)
			if v, ok := ls2["__name__"]; ok {
				n.Name = string(v)
			}
		case *promql.MatrixSelector:
			ls := Matchers2LabelSet(n.LabelMatchers)
			ls2 := pl.Process(ls, configs...)
			MargeLabelSet(ls2, n.LabelMatchers)
			if v, ok := ls2["__name__"]; ok {
				n.Name = string(v)
			}
		}
		return true
	})
	return expr.String(), nil
}
