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
		in, e := dispatcher.NewInput(request)
		if e != nil {
			err = e
			return
		}
		node := d.FindNode(in)
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
