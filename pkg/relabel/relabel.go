package relabel

import (
	"fmt"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/relabel"
	"github.com/prometheus/prometheus/promql/parser"
	"sort"
)

func Matchers2Labels(ms []*labels.Matcher) labels.Labels {
	var ls labels.Labels
	for _, m := range ms {
		ls = append(ls, labels.Label{Name: m.Name, Value: m.Value})
	}
	return ls
}

func MergeLabels(ls labels.Labels, dst []*labels.Matcher) {
	for _, m := range dst {
		if v := ls.Get(m.Name); v != "" {
			m.Value = v
		}
	}
}

func Process(query string, configs []*relabel.Config) (string, error) {
	expr, err := parser.ParseExpr(query)
	if err != nil {
		return "", err
	}
	parser.Inspect(expr, func(node parser.Node, nodes []parser.Node) error {
		switch n := node.(type) {
		case *parser.VectorSelector:
			ls := Matchers2Labels(n.LabelMatchers)
			relabeledLs := relabel.Process(ls, configs...)
			if relabeledLs == nil || labels.Equal(ls, relabeledLs) {
				// did not relabel, keep original node
				return nil
			}
			sort.Sort(ls)
			if name, isDup := ls.HasDuplicateLabelNames(); isDup {
				e := fmt.Errorf("duplicate label names (%s) can cause unexpected query result after relabeling."+
					"consider to avoid using same label name in a query, or to avoid relabeling", name)
				err = e
				return e
			}
			ls = relabeledLs
			MergeLabels(ls, n.LabelMatchers)
			// relabel a metric name not a __name__ label
			// e.g. foo{bar="baz"} => relabeled_foo{bar="baz"}
			if v := ls.Get(labels.MetricName); v != "" && n.Name != "" {
				n.Name = v
			}
			// https://github.com/prometheus/prometheus/blob/v2.23.0/promql/parser/printer.go#L161
			// but ms.Value has already been relabeled, that differ from original n.Name
			for _, ms := range n.LabelMatchers {
				if ms.Name == labels.MetricName && ms.Type == labels.MatchEqual {
					n.Name = ms.Value
				}
			}
		}
		return nil
	})
	return expr.String(), err
}
