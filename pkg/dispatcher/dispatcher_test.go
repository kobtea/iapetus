package dispatcher

import (
	"github.com/kobtea/iapetus/pkg/model"
	pm "github.com/prometheus/common/model"
	"reflect"
	"testing"
	"time"
)

type testCase struct {
	input  Input
	target *model.Node
}

func TestDispatcher_FindNode(t *testing.T) {
	clusters := []struct {
		cluster model.Cluster
		tests   []testCase
	}{
		{
			cluster: model.Cluster{
				Name: "invalid empty nodes",
			},
			tests: []testCase{
				{
					input:  Input{},
					target: nil,
				},
			},
		},
		{
			cluster: model.Cluster{
				Name: "default",
				Nodes: []model.Node{
					{Name: "one"},
					{Name: "two"},
				},
				Rules: []model.Rule{
					{Target: "one", Default: true},
				},
			},
			tests: []testCase{
				{
					input:  Input{},
					target: &model.Node{Name: "one"},
				},
			},
		},
		{
			cluster: model.Cluster{
				Name: "range",
				Nodes: []model.Node{
					{Name: "one"},
					{Name: "two"},
				},
				Rules: []model.Rule{
					{Target: "one", Default: true},
					{Target: "two", Range: model.DurationCriteria{Op: ">", Duration: time.Hour}},
				},
			},
			tests: []testCase{
				{
					input: Input{
						Query: `foo`,
						start: time.Unix(0, 0),
						end:   time.Unix(3601, 0),
					},
					target: &model.Node{Name: "two"},
				},
				{
					input: Input{
						Query: `foo`,
						start: time.Unix(0, 0),
						end:   time.Unix(3600, 0),
					},
					target: &model.Node{Name: "one"},
				},
			},
		},
		{
			cluster: model.Cluster{
				Name: "time",
				Nodes: []model.Node{
					{Name: "one"},
					{Name: "two"},
				},
				Rules: []model.Rule{
					{Target: "one", Default: true},
					{Target: "two", Time: model.TimeCriteria{Op: ">", Time: time.Unix(1000, 0)} },
				},
			},
			tests: []testCase{
				{
					input: Input{
						Query: `foo`,
						time: time.Unix(2000, 0),
					},
					target: &model.Node{Name: "two"},
				},
				{
					input: Input{
						Query: `foo`,
						time: time.Unix(500, 0),
					},
					target: &model.Node{Name: "one"},
				},
			},
		},
		{
			cluster: model.Cluster{
				Name: "start",
				Nodes: []model.Node{
					{Name: "one"},
					{Name: "two"},
				},
				Rules: []model.Rule{
					{Target: "one", Default: true},
					{Target: "two", Start: model.TimeCriteria{Op: ">", Time: time.Unix(1000, 0)} },
				},
			},
			tests: []testCase{
				{
					input: Input{
						Query: `foo`,
						start: time.Unix(2000, 0),
					},
					target: &model.Node{Name: "two"},
				},
				{
					input: Input{
						Query: `foo`,
						start: time.Unix(500, 0),
					},
					target: &model.Node{Name: "one"},
				},
			},
		},
		{
			cluster: model.Cluster{
				Name: "end",
				Nodes: []model.Node{
					{Name: "one"},
					{Name: "two"},
				},
				Rules: []model.Rule{
					{Target: "one", Default: true},
					{Target: "two", End: model.TimeCriteria{Op: ">", Time: time.Unix(1000, 0)} },
				},
			},
			tests: []testCase{
				{
					input: Input{
						Query: `foo`,
						end: time.Unix(2000, 0),
					},
					target: &model.Node{Name: "two"},
				},
				{
					input: Input{
						Query: `foo`,
						end: time.Unix(500, 0),
					},
					target: &model.Node{Name: "one"},
				},
			},
		},
		{
			cluster: model.Cluster{
				Name: "required_labels",
				Nodes: []model.Node{
					{Name: "one"},
					{Name: "two"},
				},
				Rules: []model.Rule{
					{Target: "one", Default: true},
					{Target: "two", RequiredLabels: pm.LabelSet{"job": "some_job"}},
				},
			},
			tests: []testCase{
				{
					input: Input{
						Query: `foo{job="some_job"} + 100`,
					},
					target: &model.Node{Name: "two"},
				},
				{
					input: Input{
						Query: `foo + 100`,
					},
					target: &model.Node{Name: "one"},
				},
			},
		},
	}

	for _, cluster := range clusters {
		dp := NewDispatcher(cluster.cluster)
		for _, test := range cluster.tests{
			target := dp.FindNode(test.input)
			if !reflect.DeepEqual(target, test.target) {
				t.Errorf("expect %v, but got %v", test.target, target)
			}
		}
	}
}
