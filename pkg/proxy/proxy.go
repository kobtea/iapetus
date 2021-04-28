package proxy

import (
	"fmt"
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
	var err error
	director := func(request *http.Request) {
		if len(request.URL.Scheme) == 0 {
			request.URL.Scheme = "http"
		}
		reqUrl := *request.URL
		in, e := dispatcher.NewInput(request)
		if e != nil {
			err = e
			return
		}
		node := d.FindNode(in)

		request.ParseForm()
		if v, ok := request.Form["query"]; ok {
			// update query
			in.Query, e = relabel.Process(in.Query, node.Relabels)
			if e != nil {
				request.Header.Set(headerRequestError, e.Error())
				err = e
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
				if e != nil {
					err = e
					return
				}
				q.Add("match[]", in.Matchers[i])
			}

			reqUrl.RawQuery = q.Encode()
		}

		nodeUrl, e := url.Parse(node.Url)
		if e != nil {
			err = e
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

		req, e := http.NewRequest(request.Method, reqUrl.String(), request.Body)
		err = e
		req.Header = request.Header
		level.Debug(logger).Log("request", fmt.Sprintf("%s://%s%s", request.URL.Scheme, request.Host, request.RequestURI))
		level.Debug(logger).Log("backend", fmt.Sprintf("%s://%s%s", reqUrl.Scheme, reqUrl.Host, reqUrl.RequestURI()))
		level.Info(logger).Log("target", node.Name, "query", in.Query, "match[]", fmt.Sprintf("%+v", in.Matchers), "origin", request.URL.RawQuery)
		*request = *req
	}
	if err != nil {
		return nil, err
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
