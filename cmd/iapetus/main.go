package main

import (
	"fmt"
	"github.com/kobtea/iapetus/pkg/config"
	"github.com/kobtea/iapetus/pkg/proxy"
	"github.com/kobtea/iapetus/pkg/util"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"net/http"
)

var (
	configFile   = kingpin.Flag("config", "iapetus config file path.").Required().String()
	listenAddr   = kingpin.Flag("listen.addr", "address to listen.").Default(":19090").String()
	listenPrefix = kingpin.Flag("listen.prefix", "path prefix of this endpoint. remove this prefix when dispatch to a backend.").String()
	logLevel     = kingpin.Flag("log.level", "log level (debug, info, warn, error)").String()
)

func main() {
	kingpin.Version(version.Print("iapetus"))
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

	if len(*listenAddr) > 0 {
		c.Listen.Addr = *listenAddr
	}
	if len(*listenPrefix) > 0 {
		c.Listen.Prefix = *listenPrefix
	}
	if len(*logLevel) > 0 {
		c.Log.Level = *logLevel
		util.SetLogFilterLevel(c.Log.Level)
	}

	if err := config.Validate(c); err != nil {
		for _, e := range err {
			fmt.Println(e.Error())
		}
		return
	}

	handler, err := proxy.NewProxyHandler(*c)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	server := http.Server{
		Addr:     c.Listen.Addr,
		Handler:  handler,
		ErrorLog: util.GetStdErrorLogger(),
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}
}
