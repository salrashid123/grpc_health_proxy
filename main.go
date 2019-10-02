// Copyright 2018 Google LLC
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
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/golang/glog"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"golang.org/x/net/http2"
)

type ProbeConfig struct {
	flGrpcServerAddr string
	flRunCli         bool
	flHTTPListenAddr string
	flHTTPListenPath string
	flServiceName   string
	flUserAgent     string
	flConnTimeout   time.Duration
	flRPCTimeout    time.Duration
	flGrpcTLS       bool
	flGrpcTLSNoVerify   bool
	flGrpcTLSCACert     string
	flGrpcTLSClientCert string
	flGrpcTLSClientKey  string
	flGrpcSNIServerName string
	flHTTPSTLSServerCert string
	flHTTPSTLSServerKey string
	flHTTPSTLSVerifyCA string
	flHTTPSTLSVerifyClient bool
}

var (
	cfg = &ProbeConfig{}
	opts = []grpc.DialOption{}
)

type GrpcProbeError struct {
	Code int
	Message string
}

func NewGrpcProbeError(code int, message string) *GrpcProbeError {
    return &GrpcProbeError{
		Code: code,
		Message: message,
	}
}
func (e *GrpcProbeError) Error() string {
    return e.Message
}
const (
	StatusConnectionFailure = 1
	StatusRPCFailure = 2
	StatusServiceNotFound = 3
	StatusUnimplemented = 4
	StatusUnhealthy = 5
)

func init() {
	flag.StringVar(&cfg.flGrpcServerAddr, "grpcaddr", "", "(required) tcp host:port to connect")
	flag.StringVar(&cfg.flServiceName, "service-name", "", "service name to check.  If specified, server will ignore ?serviceName= request parameter")
	flag.StringVar(&cfg.flUserAgent, "user-agent", "grpc_health_proxy", "user-agent header value of health check requests")
	flag.BoolVar(&cfg.flRunCli, "runcli", false, "execute healthCheck via CLI; will not start webserver")
	// settings for HTTPS lisenter
	flag.StringVar(&cfg.flHTTPListenAddr, "http-listen-addr", "localhost:8080", "(required) http host:port to listen (default: localhost:8080")
	flag.StringVar(&cfg.flHTTPListenPath, "http-listen-path", "/", "path to listen for healthcheck traffic (default '/')")
	flag.StringVar(&cfg.flHTTPSTLSServerCert, "https-listen-cert", "", "TLS Server certificate to for HTTP listner")
	flag.StringVar(&cfg.flHTTPSTLSServerKey, "https-listen-key", "", "TLS Server certificate key to for HTTP listner")
	flag.StringVar(&cfg.flHTTPSTLSVerifyCA, "https-listen-ca", "", "Use CA to verify client requests against CA")
	flag.BoolVar(&cfg.flHTTPSTLSVerifyClient, "https-listen-verify", false, "Verify client certificate provided to the HTTP listner")
	// timeouts
	flag.DurationVar(&cfg.flConnTimeout, "connect-timeout", time.Second, "timeout for establishing connection")
	flag.DurationVar(&cfg.flRPCTimeout, "rpc-timeout", time.Second, "timeout for health check rpc")
	// tls settings
	flag.BoolVar(&cfg.flGrpcTLS, "grpctls", false, "use TLS for upstream gRPC(default: false, INSECURE plaintext transport)")
	flag.BoolVar(&cfg.flGrpcTLSNoVerify, "grpc-tls-no-verify", false, "(with -tls) don't verify the certificate (INSECURE) presented by the server (default: false)")
	flag.StringVar(&cfg.flGrpcTLSCACert, "grpc-ca-cert", "", "(with -tls, optional) file containing trusted certificates for verifying server")
	flag.StringVar(&cfg.flGrpcTLSClientCert, "grpc-client-cert", "", "(with -grpctls, optional) client certificate for authenticating to the server (requires -tls-client-key)")
	flag.StringVar(&cfg.flGrpcTLSClientKey, "grpc-client-key", "", "(with -grpctls) client private key for authenticating to the server (requires -tls-client-cert)")
	flag.StringVar(&cfg.flGrpcSNIServerName, "grpc-sni-server-name", "", "(with -grpctls) override the hostname used to verify the gRPC server certificate")

	flag.Parse()

	argError := func(s string, v ...interface{}) {
		//flag.PrintDefaults()
		glog.Errorf("Invalid Argument error: "+s, v...)
		os.Exit(-1)
	}

	if cfg.flGrpcServerAddr == "" {		
		argError("-grpcaddr not specified")		
	}
	if !cfg.flRunCli && cfg.flHTTPListenAddr == "" {
		argError("-http-listen-addr not specified")
	}	
	if cfg.flConnTimeout <= 0 {
		argError("-connect-timeout must be greater than zero (specified: %v)", cfg.flConnTimeout)
	}
	if cfg.flRPCTimeout <= 0 {
		argError("-rpc-timeout must be greater than zero (specified: %v)", cfg.flRPCTimeout)
	}
	if !cfg.flGrpcTLS && cfg.flGrpcTLSNoVerify {
		argError("specified -grpc-tls-no-verify without specifying -grpctls")
	}
	if !cfg.flGrpcTLS && cfg.flGrpcTLSCACert != "" {
		argError("specified -grpc-ca-cert without specifying -grpctls")
	}
	if !cfg.flGrpcTLS && cfg.flGrpcTLSClientCert != "" {
		argError("specified -grpc-client-cert without specifying -grpctls")
	}
	if !cfg.flGrpcTLS && cfg.flGrpcSNIServerName != "" {
		argError("specified -grpc-sni-server-name without specifying -grpctls")
	}
	if cfg.flGrpcTLSClientCert != "" && cfg.flGrpcTLSClientKey == "" {
		argError("specified -grpc-client-cert without specifying -grpc-client-key")
	}
	if cfg.flGrpcTLSClientCert == "" && cfg.flGrpcTLSClientKey != "" {
		argError("specified -grpc-client-key without specifying -grpc-client-cert")
	}
	if cfg.flGrpcTLSNoVerify && cfg.flGrpcTLSCACert != "" {
		argError("cannot specify -grpc-ca-cert with -grpc-tls-no-verify (CA cert would not be used)")
	}
	if cfg.flGrpcTLSNoVerify && cfg.flGrpcSNIServerName != "" {
		argError("cannot specify -grpc-sni-server-name with -grpc-tls-no-verify (server name would not be used)")
	}
	if ( (cfg.flHTTPSTLSServerCert == "" && cfg.flHTTPSTLSServerKey != "") || (cfg.flHTTPSTLSServerCert != "" && cfg.flHTTPSTLSServerKey == "") ) {
		argError("must specify both -https-listen-cert and -https-listen-key")
	}
	if cfg.flHTTPSTLSVerifyCA == "" && cfg.flHTTPSTLSVerifyClient {
		argError("cannot specify -https-listen-ca if https-listen-verify is set (you need a trust CA for client certificate https auth)")
	}
	
	glog.V(10).Infof("parsed options:")
	glog.V(10).Infof("> addr=%s conn_timeout=%s rpc_timeout=%s", cfg.flGrpcServerAddr, cfg.flConnTimeout, cfg.flRPCTimeout)
	glog.V(10).Infof("> grpctls=%v", cfg.flGrpcTLS)
	glog.V(10).Infof(" http-listen-addr=%s ", cfg.flHTTPListenAddr)
	glog.V(10).Infof(" http-listen-path=%s ", cfg.flHTTPListenPath)

	glog.V(10).Infof(" https-listen-cert=%s ", cfg.flHTTPSTLSServerCert)
	glog.V(10).Infof(" https-listen-key=%s ", cfg.flHTTPSTLSServerKey)
	glog.V(10).Infof(" https-listen-verify=%v ", cfg.flHTTPSTLSVerifyClient)
	glog.V(10).Infof(" https-listen-ca=%s ", cfg.flHTTPSTLSVerifyCA)
	glog.V(10).Infof("  > grpc-tls-no-verify=%v ", cfg.flGrpcTLSNoVerify)
	glog.V(10).Infof("  > grpc-ca-cert=%s", cfg.flGrpcTLSCACert)
	glog.V(10).Infof("  > grpc-client-cert=%s", cfg.flGrpcTLSClientCert)
	glog.V(10).Infof("  > grpc-client-key=%s", cfg.flGrpcTLSClientKey)
	glog.V(10).Infof("  > grpc-sni-server-name=%s", cfg.flGrpcSNIServerName)
}

func buildGrpcCredentials() (credentials.TransportCredentials, error) {
	var tlsCfg tls.Config

	if cfg.flGrpcTLSClientCert != "" && cfg.flGrpcTLSClientKey != "" {
		keyPair, err := tls.LoadX509KeyPair(cfg.flGrpcTLSClientCert, cfg.flGrpcTLSClientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load tls client cert/key pair. error=%v", err)
		}
		tlsCfg.Certificates = []tls.Certificate{keyPair}
	}

	if cfg.flGrpcTLSNoVerify {
		tlsCfg.InsecureSkipVerify = true
	} else if cfg.flGrpcTLSCACert != "" {
		rootCAs := x509.NewCertPool()
		pem, err := ioutil.ReadFile(cfg.flGrpcTLSCACert)
		if err != nil {
			return nil, fmt.Errorf("failed to load root CA certificates from file (%s) error=%v", cfg.flGrpcTLSCACert, err)
		}
		if !rootCAs.AppendCertsFromPEM(pem) {
			return nil, fmt.Errorf("no root CA certs parsed from file %s", cfg.flGrpcTLSCACert)
		}
		tlsCfg.RootCAs = rootCAs
	}
	if cfg.flGrpcSNIServerName != "" {
		tlsCfg.ServerName = cfg.flGrpcSNIServerName
	}
	return credentials.NewTLS(&tlsCfg), nil
}

func checkService(ctx context.Context, serviceName string) (healthpb.HealthCheckResponse_ServingStatus, error) {

	glog.V(10).Infof("establishing connection")
	connStart := time.Now()
	dialCtx, dialCancel := context.WithTimeout(ctx, cfg.flConnTimeout)
	defer dialCancel()
	conn, err := grpc.DialContext(dialCtx, cfg.flGrpcServerAddr, opts...)
	if err != nil {
		if err == context.DeadlineExceeded {
			glog.Warningf("timeout: failed to connect service %s within %s", cfg.flGrpcServerAddr, cfg.flConnTimeout)
		} else {
			glog.Warningf("error: failed to connect service at %s: %+v", cfg.flGrpcServerAddr, err)
		}
		return healthpb.HealthCheckResponse_UNKNOWN, NewGrpcProbeError(StatusConnectionFailure, "StatusConnectionFailure")
	}
	connDuration := time.Since(connStart)
	defer conn.Close()
	glog.V(10).Infof("connection established %v", connDuration)

	rpcStart := time.Now()
	rpcCtx, rpcCancel := context.WithTimeout(ctx, cfg.flRPCTimeout)
	defer rpcCancel()

	glog.V(10).Infoln("Running HealthCheck for service: ", serviceName)

	resp, err := healthpb.NewHealthClient(conn).Check(rpcCtx, &healthpb.HealthCheckRequest{Service: serviceName})
	if err != nil {
		// first handle and return gRPC-level errors
		if stat, ok := status.FromError(err); ok && stat.Code() == codes.Unimplemented {
			glog.Warningf("error: this server does not implement the grpc health protocol (grpc.health.v1.Health)")
			return healthpb.HealthCheckResponse_UNKNOWN, NewGrpcProbeError(StatusUnimplemented, "StatusUnimplemented")
		} else if stat, ok := status.FromError(err); ok && stat.Code() == codes.DeadlineExceeded {
			glog.Warningf("error timeout: health rpc did not complete within ", cfg.flRPCTimeout)
			return healthpb.HealthCheckResponse_UNKNOWN, NewGrpcProbeError(StatusRPCFailure, "StatusRPCFailure")
		} else if stat, ok := status.FromError(err); ok && stat.Code() == codes.NotFound {
			// wrap a grpC NOT_FOUND as grpcProbeError.
			// https://github.com/grpc/grpc/blob/master/doc/health-checking.md
			// if the service name is not registerered, the server returns a NOT_FOUND GPRPC status. 
			// the Check for a not found should "return nil, status.Error(codes.NotFound, "unknown service")"
			glog.Warningf("error Service Not Found %v", err )
			return healthpb.HealthCheckResponse_SERVICE_UNKNOWN, NewGrpcProbeError(StatusServiceNotFound, "StatusServiceNotFound")
		} else {
			glog.Warningf("error: health rpc failed: ", err)
		}		
	}
	rpcDuration := time.Since(rpcStart)
	// otherwise, retrurn gRPC-HC status
	glog.V(10).Infof("time elapsed: connect=%s rpc=%s", connDuration, rpcDuration)

	return resp.GetStatus(), nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {

	var serviceName string
	if (cfg.flServiceName != "") {
		serviceName = cfg.flServiceName
	}	
	keys, ok := r.URL.Query()["serviceName"]
	if ok && len(keys[0]) > 0 {
		serviceName = keys[0]
	}

	resp, err := checkService(r.Context(), serviceName)
	// first handle errors derived from gRPC-codes
	if err != nil {
        if pe, ok := err.(*GrpcProbeError); ok {
			glog.Errorf("HealtCheck Probe Error: %v", pe.Error())
			switch pe.Code {
			case StatusConnectionFailure:
				http.Error(w, err.Error(), http.StatusBadGateway)
			case StatusRPCFailure:
				http.Error(w, err.Error(), http.StatusBadGateway)
			case StatusUnimplemented:
				http.Error(w, err.Error(), http.StatusNotImplemented)				
			case StatusServiceNotFound:
				http.Error(w, fmt.Sprintf("%s ServiceNotFound", cfg.flServiceName), http.StatusNotFound)									
			default:
				http.Error(w, err.Error(), http.StatusBadGateway)
			}
			return
		}
	}

	// then grpc-hc codes
	glog.Infof("%s %v",  cfg.flServiceName,  resp.String())
	switch resp {
	case healthpb.HealthCheckResponse_SERVING:
		fmt.Fprintf(w, "%s %v", cfg.flServiceName, resp)
	case healthpb.HealthCheckResponse_NOT_SERVING:
		http.Error(w, fmt.Sprintf("%s %v", cfg.flServiceName, resp.String()), http.StatusBadGateway)	
	case healthpb.HealthCheckResponse_UNKNOWN:
		http.Error(w, fmt.Sprintf("%s %v", cfg.flServiceName, resp.String()), http.StatusBadGateway)	
	case healthpb.HealthCheckResponse_SERVICE_UNKNOWN:
		http.Error(w, fmt.Sprintf("%s %v", cfg.flServiceName, resp.String()), http.StatusNotFound)	
	}
}

func main() {

	opts = append(opts, grpc.WithUserAgent(cfg.flUserAgent))
	opts = append(opts, grpc.WithBlock())		
	if cfg.flGrpcTLS {
		creds, err := buildGrpcCredentials()
		if err != nil {
			glog.Fatalf("failed to initialize tls credentials. error=%v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	
	if (cfg.flRunCli) {
		resp, err := checkService(context.Background(), cfg.flServiceName)
		if err != nil {
			if pe, ok := err.(*GrpcProbeError); ok {
				glog.Errorf("HealtCheck Probe Error: %v", pe.Error())
				switch pe.Code {
				case StatusConnectionFailure:
					os.Exit(StatusConnectionFailure)
				case StatusRPCFailure:
					os.Exit(StatusRPCFailure)
				case StatusUnimplemented:
					os.Exit(StatusUnimplemented)
				case StatusServiceNotFound:
					os.Exit(StatusServiceNotFound)
				default:
					os.Exit(StatusUnhealthy)
				}
			}
		}
		if (resp != healthpb.HealthCheckResponse_SERVING) {
			glog.Errorf("HealtCheck Probe Error: service %s failed with reason: %v",  cfg.flServiceName,  resp.String())
			os.Exit(StatusUnhealthy)
		} else {
			glog.Infof("%s %v",  cfg.flServiceName,  resp.String())
		}
	} else {

		tlsConfig := &tls.Config{}
		if (cfg.flHTTPSTLSVerifyClient) {
			caCert, err := ioutil.ReadFile(cfg.flHTTPSTLSVerifyCA)
			if err != nil {
				glog.Fatal(err)
			}
			caCertPool := x509.NewCertPool()
			if !caCertPool.AppendCertsFromPEM(caCert){
				glog.Fatal("Unable to add https server root CA certs")
			}

			tlsConfig = &tls.Config{
				ClientCAs: caCertPool,
				ClientAuth: tls.RequireAndVerifyClientCert,
			}
		}
		tlsConfig.BuildNameToCertificate()

		srv := &http.Server{
			Addr: cfg.flHTTPListenAddr,
			TLSConfig: tlsConfig,
		}
		http2.ConfigureServer(srv, &http2.Server{})
		http.HandleFunc(cfg.flHTTPListenPath, healthHandler)
		
		var err error
		if (cfg.flHTTPSTLSServerCert != "" && cfg.flHTTPSTLSServerKey != "" ) {
			err = srv.ListenAndServeTLS(cfg.flHTTPSTLSServerCert, cfg.flHTTPSTLSServerKey)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil {
			glog.Fatalf("ListenAndServe Error: ", err)
		}
	}
}
