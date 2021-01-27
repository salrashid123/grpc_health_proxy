module main

go 1.13

require (
	echo v0.0.0
	github.com/golang/protobuf v1.4.3 // indirect
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	google.golang.org/grpc v1.28.0
	google.golang.org/protobuf v1.25.0 // indirect
)

replace echo => ./src/echo
