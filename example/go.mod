module main

go 1.24.0

require (
	github.com/golang/glog v1.2.5
	github.com/salrashid123/grpc_health_proxy/example/src/echo v0.0.0
	golang.org/x/net v0.49.0
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11 // indirect
)

require (
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260203192932-546029d2fa20 // indirect
)

replace github.com/salrashid123/grpc_health_proxy/example/src/echo => ./src/echo
