package main

import (
	"fmt"
	"github.com/kobtea/iapetus/pkg/config"
	"github.com/kobtea/iapetus/pkg/proxy"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"net/http"
)

var (
	configFile = kingpin.Flag("config", "iapetus config file path.").Required().String()
	addr       = kingpin.Flag("addr", "address to listen.").Default(":19090").String()
	logLevel   = kingpin.Flag("log.level", "log level (debug, info, warn, error)").String()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	buf, err := ioutil.ReadFile(*configFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c, err := config.Parse(buf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// override log level
	if len(*logLevel) > 0 {
		c.Log.Level = *logLevel
	}

	handler, err := proxy.NewProxyHandler(*c)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	server := http.Server{
		Addr:    *addr,
		Handler: handler,
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}
}
