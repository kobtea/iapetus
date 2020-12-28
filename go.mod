module github.com/kobtea/iapetus

go 1.15

require (
	github.com/go-kit/kit v0.10.0
	github.com/prometheus/common v0.15.0
	// https://github.com/prometheus/prometheus/issues/7663
	github.com/prometheus/prometheus v1.8.2-0.20201126101154-26d89b4b0776 // v2.23.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/yaml.v2 v2.3.0
)
