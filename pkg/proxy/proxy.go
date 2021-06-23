package proxy

import (
	"fmt"
	"io"
	"github.com/go-kit/kit/log/level"
	"github.com/kobtea/iapetus/pkg/config"
	"github.com/kobtea/iapetus/pkg/dispatcher"
	"github.com/kobtea/iapetus/pkg/relabel"
	"github.com/kobtea/iapetus/pkg/util"
	phttputil "github.com/prometheus/prometheus/util/httputil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"path"
	"regexp"
	"strings"
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

func NewProxyHandler(config config.Config) (http.Handler, error) {
	logger := util.NewLogger(config.Log.Level)
	cluster := config.Clusters[0] // TODO: support multi clusters
	d := dispatcher.NewDispatcher(cluster)
	director := func(request *http.Request) {
		level.Debug(logger).Log("request", fmt.Sprintf("%s://%s%s", request.URL.Scheme, request.Host, request.RequestURI))
		if len(request.URL.Scheme) == 0 {
			request.URL.Scheme = "http"
		}
		reqUrl := *request.URL
		in, err := dispatcher.NewInput(request)
		if err != nil {
			request.Header.Set(headerRequestError, err.Error())
			level.Warn(logger).Log("msg", err.Error())
			return
		}
		node := d.FindNode(in)

		request.ParseForm()
		if v, ok := request.Form["query"]; ok {
			// update query
			in.Query, err = relabel.Process(in.Query, node.Relabels)
			if err != nil {
				request.Header.Set(headerRequestError, err.Error())
				level.Warn(logger).Log("msg", err.Error())
				return
			}

			if in.Query != v[0] {
				q := reqUrl.Query()
				q.Set("query", in.Query)
				reqUrl.RawQuery = q.Encode()
			}
		}
		if _, ok := request.Form["match[]"]; ok {
			q := reqUrl.Query()
			q.Del("match[]")

			for i := range in.Matchers {
				// update query
				in.Matchers[i], err = relabel.Process(in.Matchers[i], node.Relabels)
				if err != nil {
					request.Header.Set(headerRequestError, err.Error())
					level.Warn(logger).Log("msg", err.Error())
					return
				}
				q.Add("match[]", in.Matchers[i])
			}

			reqUrl.RawQuery = q.Encode()
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

		var body io.Reader
		body = request.Body
		if request.Method == "POST" {
			body = strings.NewReader(request.Form.Encode())
		}
		req, err := http.NewRequest(request.Method, reqUrl.String(), body)
		if err != nil {
			level.Error(logger).Log("msg", err.Error())
			return
		}
		req.Header = request.Header
		level.Debug(logger).Log("backend", fmt.Sprintf("%s://%s%s", reqUrl.Scheme, reqUrl.Host, reqUrl.RequestURI()))
		level.Info(logger).Log("target", node.Name, "query", in.Query, "match[]", fmt.Sprintf("%+v", in.Matchers), "origin", request.URL.RawQuery)
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
