module app-parity

go 1.12

require (
	github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc // indirect
	github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf // indirect
	github.com/bilibili/kratos v0.2.0
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.1
	github.com/pkg/errors v0.8.1
	github.com/prometheus/common v0.0.0-20181126121408-4724e9255275
	golang.org/x/net v0.0.0-20190628185345-da137c7871d7
	google.golang.org/genproto v0.0.0-20190701230453-710ae3a149df
	google.golang.org/grpc v1.22.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6 // indirect
)

replace github.com/bilibili/kratos => github.com/C-isCoder/kratos v0.2.1-0.20190722040647-89ee902bb4f4
