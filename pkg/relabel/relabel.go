package relabel

import (
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/relabel"
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

func MergeLabelSet(ls model.LabelSet, dst []*labels.Matcher) {
	for _, m := range dst {
		if v, ok := ls[model.LabelName(m.Name)]; ok {
			m.Value = string(v)
		}
	}
}

func Process(query string, configs []*relabel.Config) (string, error) {
	expr, err := promql.ParseExpr(query)
	if err != nil {
		return "", err
	}
	promql.Inspect(expr, func(node promql.Node, nodes []promql.Node) error {
		switch n := node.(type) {
		case *promql.VectorSelector:
			ls := Matchers2LabelSet(n.LabelMatchers)
			ls2 := Matchers2LabelSet(n.LabelMatchers)
			pl.Process(ls2, configs...)
			MergeLabelSet(ls2, n.LabelMatchers)
			if v, ok := ls["__name__"]; ok && string(v) == n.Name {
				if v2, ok2 := ls2["__name__"]; ok2 {
					n.Name = string(v2)
				}
			}
			// https://github.com/prometheus/prometheus/blob/v2.2.1/promql/printer.go#L218-L221
			for _, ms := range n.LabelMatchers {
				if ms.Name == "__name__" && ms.Type == labels.MatchEqual {
					n.Name = ms.Value
				}
			}
		case *promql.MatrixSelector:
			ls := Matchers2LabelSet(n.LabelMatchers)
			ls2 := Matchers2LabelSet(n.LabelMatchers)
			pl.Process(ls2, configs...)
			MergeLabelSet(ls2, n.LabelMatchers)
			if v, ok := ls["__name__"]; ok && string(v) == n.Name {
				if v2, ok2 := ls2["__name__"]; ok2 {
					n.Name = string(v2)
				}
			}
		}
		return nil
	})
	return expr.String(), nil
}
