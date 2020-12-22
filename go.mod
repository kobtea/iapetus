module github.com/kobtea/iapetus

go 1.15

require (
	github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc
	github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf
	github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973
	github.com/fsnotify/fsnotify v1.4.7
	github.com/go-kit/kit v0.8.0
	github.com/go-logfmt/logfmt v0.4.0
	github.com/golang/protobuf v1.2.0
	github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515
	github.com/matttproud/golang_protobuf_extensions v1.0.1
	github.com/miekg/dns v1.0.4
	github.com/oklog/ulid v1.3.1
	github.com/opentracing/opentracing-go v1.0.2
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.2
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90
	github.com/prometheus/common v0.0.0-20181119215939-b36ad289a3ea
	github.com/prometheus/procfs v0.0.0-20190219184716-e4d4a2206da0
	github.com/prometheus/tsdb v0.4.0
	golang.org/x/crypto v0.0.0-20190222235706-ffb98f73852f
	golang.org/x/net v0.0.0-20190213061140-3a22650c66bd
	golang.org/x/sync v0.0.0-20181221193216-37e7f081c4d4
	golang.org/x/sys v0.0.0-20190225065934-cc5685c2db12
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/yaml.v2 v2.2.2
)

replace gopkg.in/fsnotify.v1 v1.4.7 => github.com/fsnotify/fsnotify v1.4.7
