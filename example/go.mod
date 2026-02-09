module main

go 1.24.0

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/salrashid123/grpc_health_proxy/example/src/echo v0.0.0
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1 // indirect
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20211223182754-3ac035c7e7cb // indirect
)

replace github.com/salrashid123/grpc_health_proxy/example/src/echo => ./src/echo
