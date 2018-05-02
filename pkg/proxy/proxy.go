package proxy

import (
	"github.com/kobtea/iapetus/pkg/config"
	"github.com/kobtea/iapetus/pkg/dispatcher"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxyHandler(config config.Config) (http.Handler, error) {
	cluster := config.Clusters[0] // TODO: support multi clusters
	d := dispatcher.NewDispatcher(cluster)
	var err error
	director := func(request *http.Request) {
		reqUrl := *request.URL
		target := d.FindNode()
		targetUrl, e := url.Parse(target.Url)
		if e != nil {
			err = e
		}
		reqUrl.Scheme = targetUrl.Scheme
		if len(reqUrl.Scheme) == 0 {
			reqUrl.Scheme = "http"
		}
		reqUrl.Host = targetUrl.Host
		req, e := http.NewRequest(request.Method, reqUrl.String(), request.Body)
		err = e
		req.Header = request.Header
		*request = *req
	}
	if err != nil {
		return nil, err
	}
	proxy := &httputil.ReverseProxy{Director: director}
	return proxy, nil
}
