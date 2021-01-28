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

// Sample gRPC Client application demonstrating
// gRPC Client request and API call
//
// Usage:
// 1. Start gRPC Server.
// 2  Run client and connect
//    go run src/grpc_client.go --host localhost:50051 --insecure
// 3. Optionally add TLS settings on client and server:
//    go run src/grpc_client.go --host localhost:50051 --tlsCert certs/CA_crt.pem -skipHealthCheck --servername grpc.domain.com

package main

import (
	"crypto/tls"
	"crypto/x509"
	"echo"
	"flag"
	"io/ioutil"
	"time"

	log "github.com/golang/glog"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const ()

var (
	conn *grpc.ClientConn
)

func main() {

	address := flag.String("host", "localhost:50051", "host:port of gRPC server")
	insecure := flag.Bool("insecure", false, "connect without TLS")
	skipHealthCheck := flag.Bool("skipHealthCheck", false, "Skip Initial Healthcheck")
	tlsCert := flag.String("tlsCert", "", "tls Certificate")
	serverName := flag.String("servername", "grpc.domain.com", "CACert for server")
	repeat := flag.Int("repeat", 1, "Number of Unary Requests to send")

	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "INFO")
	flag.Parse()

	var err error
	var conn *grpc.ClientConn
	if *insecure == true {
		conn, err = grpc.Dial(*address, grpc.WithInsecure())
	} else {

		var tlsCfg tls.Config
		rootCAs := x509.NewCertPool()
		pem, err := ioutil.ReadFile(*tlsCert)
		if err != nil {
			log.Fatalf("failed to load root CA certificates  error=%v", err)
		}
		if !rootCAs.AppendCertsFromPEM(pem) {
			log.Fatalf("no root CA certs parsed from file ")
		}
		tlsCfg.RootCAs = rootCAs
		tlsCfg.ServerName = *serverName

		ce := credentials.NewTLS(&tlsCfg)
		conn, err = grpc.Dial(*address, grpc.WithTransportCredentials(ce))
	}
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := echo.NewEchoServerClient(conn)
	ctx := context.Background()

	// how to perform healthcheck request manually:
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if !*skipHealthCheck {
		resp, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{Service: "echo.EchoServer"})
		if err != nil {
			log.Fatalf("HealthCheck failed %+v", err)
		}

		if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
			log.Fatalf("service not in serving state: ", resp.GetStatus().String())
		}
		log.Infof("RPC HealthChekStatus: %v\n", resp.GetStatus())
	}
	// now make a gRPC call
	for i := 0; i < *repeat; i++ {
		r, err := c.SayHello(ctx, &echo.EchoRequest{Name: "unary RPC msg "})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Infof("RPC Response: %v\n", r)
	}
}
