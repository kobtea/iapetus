# Iapetus

[![CircleCI](https://circleci.com/gh/kobtea/iapetus.svg?style=svg)](https://circleci.com/gh/kobtea/iapetus)
[![Go Report Card](https://goreportcard.com/badge/github.com/kobtea/iapetus)](https://goreportcard.com/report/github.com/kobtea/iapetus)

This project is unstable yet.
So breaking change may be happen no notice.


## Overview

Iapetus is a reverse proxy for [Prometheus](https://prometheus.io/) that aim to handle long term time series with low cost.
Like the database of [RRDTool](https://oss.oetiker.ch/rrdtool/), it get recent metrics with high-resolution, and old metrics with low-resolution.

The architecture using iapetus is simple.
Prepare prometheus nodes that have same time series data but only resolution and retention are different.
Iapetus dispatch a request to a prometheus node that match proxy rules.


## Install

### Binary

Go to https://github.com/kobtea/iapetus/releases

### Building from source

```bash
$ go get -d github.com/kobtea/iapetus
$ cd $GOPATH/src/github.com/kobtea/iapetus
$ make build
```


## Usage

```bash
$ iapetus --help
usage: iapetus --config=CONFIG [<flags>]

Flags:
  --help                         Show context-sensitive help (also try --help-long and --help-man).
  --config=CONFIG                iapetus config file path.
  --listen.addr=":19090"         address to listen.
  --listen.prefix=LISTEN.PREFIX  path prefix of this endpoint. remove this prefix when dispatch to a backend.
  --log.level=LOG.LEVEL          log level (debug, info, warn, error)
  --version                      Show application version.
```

configuration format is below.

```yml
# config.yml
log:
  # debug, info, warn, error. default is info.
  level: info
# Multi clusters are not support yet. So iapetus use 1st cluster setting.
clusters:
  - name: <string>
    # list prometheus as node
    nodes:
      - name: <string>
        url: <string>
        relabels:
          # support relabelling rules at prometheus
          [ - <relabel_config> ... ]
    # proxy rules
    # each rule are pair of `target: <node_name>` and some rule.
    # support rules are below.
    # - default: <bool>, use when no match other rules
    # - start: <op duration>, compare `start` at request parameter
    # - end: <op duration>, compare `end` at request parameter
    # - range: <op duration>, range is between `start` and `end` at request parameter
    # - required_labels: [ <label_name>: <label_value> ... ], find labels from `query` or `match[]` parameter(s). If a request satisfy this rule, Iapetus send not matched metrics but whole query send to the target. It is mean that Iapetus does not calculate values.
    rules:
      [ - <rules>, ...]
```

sample

```yml
log:
  level: info
clusters:
  - name: cluster1
    nodes:
      - name: primary
        url: http://localhost:9090
      - name: archive
        url: http://localhost:9091
        relabels:
          - source_labels: [__name__]
            target_label: __name__
            replacement: ${1}_avg
    rules:
      - target: primary
        default: true
      - target: archive
        range: "> 1d"
      - target: archive
        start: "< now-1d"
```


## Roadmap

### Metrics aggregator

Exporting aggregated metrics from a high-resolution prometheus.
It is for a low-resolution prometheus.

### Cluster

Handling multi prometheus clusters.
Proxy rule is TBD.
It may be a same rule as node, (or not).


## License

MIT
