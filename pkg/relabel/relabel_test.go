package relabel

import (
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/relabel"
	"reflect"
	"testing"
)

func TestMatchers2LabelSet(t *testing.T) {
	tests := []struct {
		m  []*labels.Matcher
		ls model.LabelSet
	}{
		{
			[]*labels.Matcher{},
			model.LabelSet{},
		},
		{
			[]*labels.Matcher{
				{
					Type:  labels.MatchEqual,
					Name:  "foo",
					Value: "bar",
				},
			},
			model.LabelSet{"foo": "bar"},
		},
		{
			[]*labels.Matcher{
				{
					Type:  labels.MatchRegexp,
					Name:  "foo",
					Value: "bar.*baz",
				},
			},
			model.LabelSet{"foo": "bar.*baz"},
		},
	}
	for _, test := range tests {
		res := Matchers2LabelSet(test.m)
		if !reflect.DeepEqual(res, test.ls) {
			t.Errorf("expect %v, but got %v", test.ls, res)
		}
	}
}

func TestMergeLabelSet(t *testing.T) {
	m := []*labels.Matcher{{Type: labels.MatchEqual, Name: "foo", Value: "bar"}}
	tests := []struct {
		ls model.LabelSet
		m  []*labels.Matcher
	}{
		{
			model.LabelSet{},
			m,
		},
		{
			model.LabelSet{"aaa": "bbb"},
			append(m),
		},
		{
			model.LabelSet{"foo": "modified"},
			[]*labels.Matcher{{Type: labels.MatchEqual, Name: "foo", Value: "modified"}},
		},
	}
	for _, test := range tests {
		dst := make([]*labels.Matcher, len(m))
		copy(dst, m)
		MergeLabelSet(test.ls, dst)
		if !reflect.DeepEqual(dst, test.m) {
			t.Errorf("expect %v, but got %v", test.m, dst)
		}
	}
}

func TestProcess(t *testing.T) {
	c := []*relabel.Config{
		{
			SourceLabels: model.LabelNames{"__name__"},
			Separator:    ";",
			Regex:        relabel.MustNewRegexp("(.*)"),
			TargetLabel:  "__name__",
			Replacement:  "${1}_avg",
			Action:       relabel.Replace,
		},
		{
			SourceLabels: model.LabelNames{"__name__"},
			Separator:    ";",
			Regex:        relabel.MustNewRegexp("_avg"),
			TargetLabel:  "__name__",
			Replacement:  "$1",
			Action:       relabel.Drop,
		},
	}
	tests := []struct {
		in  string
		out string
		err error
	}{
		{
			"foo",
			"foo_avg",
			nil,
		},
		{
			`{foo="bar"}`,
			`{foo="bar"}`,
			nil,
		},
		{
			`foo{bar="baz"}`,
			`foo_avg{bar="baz"}`,
			nil,
		},
		{
			`{__name__="foo"}`,
			`foo_avg`,
			nil,
		},
		{
			`{__name__=~"foo"}`,
			`{__name__=~"foo_avg"}`,
			nil,
		},
	}
	for _, test := range tests {
		res, err := Process(test.in, c)
		if !reflect.DeepEqual(res, test.out) {
			t.Errorf("expect %v, but got %v", test.out, res)
		}
		if !reflect.DeepEqual(err, test.err) {
			t.Errorf("expect %v, but got %v", test.err, err)
		}
	}
}
