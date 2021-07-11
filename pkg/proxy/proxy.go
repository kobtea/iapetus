package proxy

import (
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/kobtea/iapetus/pkg/config"
	"github.com/kobtea/iapetus/pkg/dispatcher"
	"github.com/kobtea/iapetus/pkg/relabel"
	"github.com/kobtea/iapetus/pkg/util"
	phttputil "github.com/prometheus/prometheus/util/httputil"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"
)

const headerRequestError = "X-Iapetus-Request-Error"

type transport struct {
	parent     http.RoundTripper
	corsOrigin *regexp.Regexp
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	if errMsg, ok := r.Header[headerRequestError]; ok {
		s := fmt.Sprintf(`{"status":"error","errorType":"bad_data","error":"%s"}`, strings.Join(errMsg, ","))
		rw := httptest.NewRecorder()
		rw.Header().Set("Content-Type", "application/json")
		phttputil.SetCORS(rw, t.corsOrigin, r)
		rw.WriteHeader(http.StatusBadRequest)
		if _, err := rw.WriteString(s); err != nil {
			return nil, err
		}
		return rw.Result(), nil
	}
	return t.parent.RoundTrip(r)
}

func formatValues(values url.Values) string {
	var ss []string
	for k, v := range values {
		ss = append(ss, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(ss, ",")
}

func NewProxyHandler(config config.Config) (http.Handler, error) {
	logger := util.NewLogger(config.Log.Level)
	cluster := config.Clusters[0] // TODO: support multi clusters
	d := dispatcher.NewDispatcher(cluster)
	director := func(request *http.Request) {
		if len(request.URL.Scheme) == 0 {
			request.URL.Scheme = "http"
		}
		reqUrl := *request.URL
		if err := request.ParseForm(); err != nil {
			level.Error(logger).Log("msg", err.Error())
			return
		}
		values := request.Form
		origValues := formatValues(values)
		level.Debug(logger).Log("request", fmt.Sprintf("%s://%s%s", request.URL.Scheme, request.Host, request.RequestURI), "values", values.Encode())

		in, err := dispatcher.NewInput(values)
		if err != nil {
			request.Header.Set(headerRequestError, err.Error())
			level.Warn(logger).Log("msg", err.Error())
			return
		}
		node := d.FindNode(in)

		if origQuery := values.Get("query"); origQuery != "" {
			// update query
			in.Query, err = relabel.Process(in.Query, node.Relabels)
			if err != nil {
				request.Header.Set(headerRequestError, err.Error())
				level.Warn(logger).Log("msg", err.Error())
				return
			}

			if in.Query != origQuery {
				values.Set("query", in.Query)
			}
		}
		if _, ok := values["match[]"]; ok {
			values.Del("match[]")
			for i := range in.Matchers {
				// update query
				in.Matchers[i], err = relabel.Process(in.Matchers[i], node.Relabels)
				if err != nil {
					request.Header.Set(headerRequestError, err.Error())
					level.Warn(logger).Log("msg", err.Error())
					return
				}
				values.Add("match[]", in.Matchers[i])
			}
		}
		if origStep := values.Get("step"); origStep != "" {
			step, err := util.ParseDuration(origStep)
			if err == nil && step < time.Duration(node.MinStep) {
				values.Set("step", node.MinStep.String())
			}
		}

		nodeUrl, err := url.Parse(node.Url)
		if err != nil {
			level.Error(logger).Log("msg", err.Error())
			return
		}
		reqUrl.Scheme = nodeUrl.Scheme
		if len(reqUrl.Scheme) == 0 {
			reqUrl.Scheme = "http"
		}
		reqUrl.Host = nodeUrl.Host
		reqUrl.Path = strings.TrimPrefix(reqUrl.Path, config.Listen.Prefix)
		if len(nodeUrl.Path) > 0 {
			reqUrl.Path = path.Join(nodeUrl.Path, reqUrl.Path)
		}

		if request.Method == http.MethodGet {
			reqUrl.RawQuery = values.Encode()
		}
		var body io.Reader
		if request.Method == http.MethodPost {
			body = strings.NewReader(values.Encode())
		}

		req, err := http.NewRequest(request.Method, reqUrl.String(), body)
		if err != nil {
			level.Error(logger).Log("msg", err.Error())
			return
		}
		req.Header = request.Header
		level.Debug(logger).Log("backend", fmt.Sprintf("%s://%s%s", reqUrl.Scheme, reqUrl.Host, reqUrl.RequestURI()), "values", values.Encode())
		level.Info(logger).Log("target", node.Name, "query", in.Query, "match[]", fmt.Sprintf("%+v", in.Matchers), "origin", origValues)
		*request = *req
	}
	proxy := &httputil.ReverseProxy{
		Director: director,
		ErrorLog: util.NewStdLogger(level.Error(logger)),
		Transport: &transport{
			parent:     http.DefaultTransport,
			corsOrigin: regexp.MustCompile("^(?:.*)$"),
		},
	}
	return proxy, nil
}
