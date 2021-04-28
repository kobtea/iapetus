package proxy

import (
	"github.com/kobtea/iapetus/pkg/config"
	"github.com/kobtea/iapetus/pkg/model"
	pm "github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/relabel"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewProxyHandler(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	d, err := model.NewDurationCriteria("> 1d")
	if err != nil {
		t.Fatal(err)
	}
	conf := config.Config{
		Clusters: []model.Cluster{{
			Name: "one",
			Nodes: []model.Node{{
				Name: "one",
				Url:  srv.URL,
			}, {
				Name: "two",
				Url:  srv.URL,
				Relabels: []*relabel.Config{{
					SourceLabels: pm.LabelNames{"__name__"},
					Separator:    ";",
					Regex:        relabel.MustNewRegexp("(foo.*)"),
					TargetLabel:  "__name__",
					Replacement:  "${1}_avg",
					Action:       relabel.Replace,
				}},
			}},
			Rules: []model.Rule{{
				Target:  "one",
				Default: true,
			}, {
				Target: "two",
				Range:  d,
			}},
		},
		},
	}
	proxy, err := NewProxyHandler(conf)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		q          url.Values
		statusCode int
	}{{
		url.Values{"query": {"foo"}},
		http.StatusOK,
	}, {
		// invalid input
		url.Values{"query": {"foo"}, "time": {"invalid_time"}},
		http.StatusBadRequest,
	}, {
		// valid relabeling
		url.Values{"query": {`{__name__="foo"}`}, "start": {"2021-01-01T00:00:00.000Z"}, "end": {"2021-01-02T01:00:00.000Z"}},
		http.StatusOK,
	}, {
		// invalid relabeling using same label keys
		url.Values{"query": {`{__name__="foo", __name__=~".+"}`}, "start": {"2021-01-01T00:00:00.000Z"}, "end": {"2021-01-02T01:00:00.000Z"}},
		http.StatusBadRequest,
	}}
	for _, test := range tests {
		req, err := http.NewRequest("GET", "/api/v1/query_range", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.URL.RawQuery = test.q.Encode()
		rec := httptest.NewRecorder()
		proxy.ServeHTTP(rec, req)
		if rec.Code != test.statusCode {
			t.Errorf("expect %d, but got %d", test.statusCode, rec.Code)
		}
	}

}
