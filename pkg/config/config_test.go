package config

import (
	"github.com/kobtea/iapetus/pkg/model"
	"reflect"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		c *Config
		e int
	}{
		{
			&Config{},
			1,
		},
		{
			&Config{
				Clusters: []model.Cluster{{}},
			},
			3,
		},
		{
			&Config{
				Clusters: []model.Cluster{{
					Name: "foo",
					Nodes: []model.Node{{
						Name: "bar",
						Url:  "baz",
					}},
					Rules: []model.Rule{{
						Default: true,
					}},
				}},
			},
			0,
		},
	}
	for _, test := range tests {
		if e := Validate(test.c); len(e) != test.e {
			t.Errorf("expect # of error %v, but got %v", test.e, len(e))
		}
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		b []byte
		c *Config
		e error
	}{
		{
			[]byte(""),
			&Config{},
			nil,
		},
		{
			[]byte(`
clusters:
  - name: foo
`),
			&Config{Clusters: []model.Cluster{{
				Name: "foo",
			}}},
			nil,
		},
	}
	for _, test := range tests {
		c, e := Parse(test.b)
		if !reflect.DeepEqual(c, test.c) {
			t.Errorf("expect %v, but got %v", test.c, c)
		}
		if !reflect.DeepEqual(e, test.e) {
			t.Errorf("expect %v, but got %v", test.e, e)
		}
	}
}
