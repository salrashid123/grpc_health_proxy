// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Sample gRPC Server application demonstrating
// gRPC healthchecks
//
// Usage:
// add to /etc/hosts:  127.0.0.1 server.domain.com
// 1. Start gRPC Server.
//    go run src/grpc_server.go --grpcport 0.0.0.0:50051 --insecure
// 3. Optionally add TLS settings:
//    go run src/grpc_server.go --grpcport 0.0.0.0:50051 --tlsCert=pub.pem --tlsKey=key.pem
// Add --unhealthyproability flag while running to simulate random healthcheck failure
// probability:  eg, --unhealthyproability 50 will return grpC "HealthCheckResponse_NOT_SERVING"

package main

import (
	"echo"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	tlsCert             = flag.String("tlsCert", "", "tls Certificate")
	tlsKey              = flag.String("tlsKey", "", "tls Key")
	grpcport            = flag.String("grpcport", "", "grpcport")
	insecure            = flag.Bool("insecure", false, "startup without TLS")
	unhealthyproability = flag.Int("unhealthyproability", 0, "percentage chance the service is unhealthy (0->100)")
)

const (
	address string = ":50051"
)

type Server struct {
	mu sync.Mutex
	// statusMap stores the serving status of the services this Server monitors.
	statusMap map[string]healthpb.HealthCheckResponse_ServingStatus
}

// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{
		statusMap: make(map[string]healthpb.HealthCheckResponse_ServingStatus),
	}
}

func (s *Server) SayHello(ctx context.Context, in *echo.EchoRequest) (*echo.EchoReply, error) {

	log.Println("Got rpc: --> ", in.Name)

	var h, err = os.Hostname()
	if err != nil {
		log.Fatalf("Unable to get hostname %v", err)
	}
	return &echo.EchoReply{Message: "Hello " + in.Name + "  from hostname " + h}, nil
}

func (s *Server) SayHelloStream(in *echo.EchoRequest, stream echo.EchoServer_SayHelloStreamServer) error {

	log.Println("Got stream:  -->  ")
	stream.Send(&echo.EchoReply{Message: "Hello " + in.Name})
	stream.Send(&echo.EchoReply{Message: "Hello " + in.Name})

	return nil
}

func (s *Server) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if in.Service == "" {
		// return overall status
		return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
	}

	r := rand.Intn(100)
	if r <= *unhealthyproability {
		s.statusMap["echo.EchoServer"] = healthpb.HealthCheckResponse_NOT_SERVING
	} else {
		s.statusMap["echo.EchoServer"] = healthpb.HealthCheckResponse_SERVING
	}
	status, ok := s.statusMap[in.Service]
	if !ok {
		// https://github.com/grpc/grpc/blob/master/doc/health-checking.md
		// "If the service name is not registered, the server returns a NOT_FOUND GRPC status."
		return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_UNKNOWN}, grpc.Errorf(codes.NotFound, "unknown service")
	}
	return &healthpb.HealthCheckResponse{Status: status}, nil
}

func (s *Server) Watch(in *healthpb.HealthCheckRequest, srv healthpb.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}

func main() {

	flag.Parse()

	if *grpcport == "" {
		fmt.Fprintln(os.Stderr, "missing -grpcport flag (:50051)")
		flag.Usage()
		os.Exit(2)
	}

	lis, err := net.Listen("tcp", *grpcport)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	sopts := []grpc.ServerOption{grpc.MaxConcurrentStreams(10)}
	if *insecure == false {
		if *tlsCert == "" || *tlsKey == "" {
			log.Fatalf("Must set --tlsCert and tlsKey if --insecure flags is used")
		}
		ce, err := credentials.NewServerTLSFromFile(*tlsCert, *tlsKey)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		sopts = append(sopts, grpc.Creds(ce))
	}

	s := grpc.NewServer(sopts...)
	srv := NewServer()
	healthpb.RegisterHealthServer(s, srv)
	echo.RegisterEchoServerServer(s, srv)
	reflection.Register(s)

	s.Serve(lis)

}
