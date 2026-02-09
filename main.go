// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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

	"net/http"
	"os"
	"time"

	"log/slog"

	"github.com/gorilla/mux"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ProbeConfig struct {
	flGrpcServerAddr        string
	flRunCli                bool
	flHTTPListenAddr        string
	flMetricsHTTPListenAddr string
	flMetricsHTTPPath       string
	flHTTPListenPath        string
	flServiceName           string
	flUserAgent             string
	flConnTimeout           time.Duration
	flRPCTimeout            time.Duration
	flGrpcTLS               bool
	flGrpcTLSNoVerify       bool
	flGrpcTLSCACert         string
	flGrpcTLSClientCert     string
	flGrpcTLSClientKey      string
	flGrpcSNIServerName     string
	flHTTPSTLSServerCert    string
	flHTTPSTLSServerKey     string
	flHTTPSTLSVerifyCA      string
	flHTTPSTLSVerifyClient  bool
	flLogTarget             string
	flJSONLog               bool
	flDebug                 bool
}

var (
	cfg  = &ProbeConfig{}
	opts = []grpc.DialOption{}

	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "grpc_health_check_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})

	serviceDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "grpc_health_check_service_duration_seconds",
		Help: "Duration of HTTP requests per service.",
	}, []string{"service_name"})

	grpcReqs = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_health_check_service_requests",
			Help: "backend status, partitioned by status code and service_name.",
		},
		[]string{"code", "service_name"},
	)

	logger *slog.Logger
)

type GrpcProbeError struct {
	Code    int
	Message string
}

func NewGrpcProbeError(code int, message string) *GrpcProbeError {
	return &GrpcProbeError{
		Code:    code,
		Message: message,
	}
}
func (e *GrpcProbeError) Error() string {
	return e.Message
}

const (
	StatusConnectionFailure = 1
	StatusRPCFailure        = 2
	StatusServiceNotFound   = 3
	StatusUnimplemented     = 4
	StatusUnhealthy         = 5
)

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}

func init() {
	flag.StringVar(&cfg.flGrpcServerAddr, "grpcaddr", "", "(required) tcp host:port to connect")
	flag.StringVar(&cfg.flServiceName, "service-name", "", "service name to check.  If specified, server will ignore ?serviceName= request parameter")
	flag.StringVar(&cfg.flUserAgent, "user-agent", "grpc_health_proxy", "user-agent header value of health check requests")
	flag.BoolVar(&cfg.flRunCli, "runcli", false, "execute healthCheck via CLI; will not start webserver")
	// settings for HTTPS listener
	flag.StringVar(&cfg.flHTTPListenAddr, "http-listen-addr", "localhost:8080", "(required) http host:port to listen (default: localhost:8080")
	flag.StringVar(&cfg.flMetricsHTTPListenAddr, "metrics-http-listen-addr", "localhost:9000", "http host:port for metrics endpoint (default: localhost:9000")
	flag.StringVar(&cfg.flMetricsHTTPPath, "metrics-http-path", "/metrics", "http path metrics endpoint (default:  /metrics")
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

	flag.StringVar(&cfg.flLogTarget, "logTarget", "", "log to file target (default stdout)")
	flag.BoolVar(&cfg.flJSONLog, "jsonLog", false, "enable json logging")
	flag.BoolVar(&cfg.flDebug, "debug", false, "enable debug logging")

	flag.Parse()

	mlogTarget := os.Stdout // default
	if cfg.flLogTarget != "" {
		var err error
		mlogTarget, err = os.OpenFile(cfg.flLogTarget, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			slog.Error("Failed to open log file", "err", err)
			os.Exit(-1)
		}
	}

	logLevel := slog.LevelInfo
	if cfg.flDebug {
		logLevel = slog.LevelDebug
	}
	if cfg.flJSONLog {
		logger = slog.New(slog.NewJSONHandler(mlogTarget, &slog.HandlerOptions{
			Level: logLevel,
		}))
	} else {
		logger = slog.New(slog.NewTextHandler(mlogTarget, &slog.HandlerOptions{
			Level: logLevel,
		}))
	}

	argError := func(s string, v ...interface{}) {
		//flag.PrintDefaults()
		logger.Error("Invalid Argument error: "+s, v...)
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
	if (cfg.flHTTPSTLSServerCert == "" && cfg.flHTTPSTLSServerKey != "") || (cfg.flHTTPSTLSServerCert != "" && cfg.flHTTPSTLSServerKey == "") {
		argError("must specify both -https-listen-cert and -https-listen-key")
	}
	if cfg.flHTTPSTLSVerifyCA == "" && cfg.flHTTPSTLSVerifyClient {
		argError("cannot specify -https-listen-ca if https-listen-verify is set (you need a trust CA for client certificate https auth)")
	}

	logger.Info("parsed options:")
	logger.Info(">", slog.String("addr", cfg.flGrpcServerAddr), slog.Duration("conn_timeout", cfg.flConnTimeout), slog.Duration("rpc_timeout", cfg.flRPCTimeout))
	logger.Info(">", slog.Bool("grpctls", cfg.flGrpcTLS))
	logger.Info(">", slog.String("http-listen-addr", cfg.flHTTPListenAddr))
	logger.Info(">", slog.String("http-listen-path", cfg.flHTTPListenPath))

	logger.Info(">", slog.String("https-listen-cert", cfg.flHTTPSTLSServerCert))
	logger.Info(">", slog.String("https-listen-key", cfg.flHTTPSTLSServerKey))
	logger.Info(">", slog.Bool("https-listen-verify", cfg.flHTTPSTLSVerifyClient))
	logger.Info(">", slog.String("https-listen-ca", cfg.flHTTPSTLSVerifyCA))
	logger.Info(">", slog.Bool("grpc-tls-no-verify", cfg.flGrpcTLSNoVerify))
	logger.Info(">", slog.String("grpc-ca-cert", cfg.flGrpcTLSCACert))
	logger.Info(">", slog.String("grpc-client-cert", cfg.flGrpcTLSClientCert))
	logger.Info(">", slog.String("grpc-client-key", cfg.flGrpcTLSClientKey))
	logger.Info(">", slog.String("grpc-sni-server-name", cfg.flGrpcSNIServerName))
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
		pem, err := os.ReadFile(cfg.flGrpcTLSCACert)
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

	timer := prometheus.NewTimer(serviceDuration.WithLabelValues(serviceName))
	defer timer.ObserveDuration()

	logger.Info("establishing connection")
	connStart := time.Now()
	dialCtx, dialCancel := context.WithTimeout(ctx, cfg.flConnTimeout)
	defer dialCancel()
	conn, err := grpc.DialContext(dialCtx, cfg.flGrpcServerAddr, opts...)
	if err != nil {
		if err == context.DeadlineExceeded {
			logger.Warn("timeout: failed to connect service %s within %s", cfg.flGrpcServerAddr, cfg.flConnTimeout)
		} else {
			logger.Warn("error: failed to connect service at %s: %+v", cfg.flGrpcServerAddr, err)
		}
		return healthpb.HealthCheckResponse_UNKNOWN, NewGrpcProbeError(StatusConnectionFailure, "StatusConnectionFailure")
	}
	connDuration := time.Since(connStart)
	defer conn.Close()
	logger.Info("connection established", slog.Duration("duration", connDuration))

	rpcStart := time.Now()
	rpcCtx, rpcCancel := context.WithTimeout(ctx, cfg.flRPCTimeout)
	defer rpcCancel()

	logger.Info("Running HealthCheck for service:", slog.String("service_name", serviceName))

	resp, err := healthpb.NewHealthClient(conn).Check(rpcCtx, &healthpb.HealthCheckRequest{Service: serviceName})
	if err != nil {
		// first handle and return gRPC-level errors
		if stat, ok := status.FromError(err); ok && stat.Code() == codes.Unimplemented {
			defer grpcReqs.WithLabelValues(codes.Unimplemented.String(), serviceName).Inc()
			logger.Warn("error: this server does not implement the grpc health protocol (grpc.health.v1.Health)")
			return healthpb.HealthCheckResponse_UNKNOWN, NewGrpcProbeError(StatusUnimplemented, "StatusUnimplemented")
		} else if stat, ok := status.FromError(err); ok && stat.Code() == codes.DeadlineExceeded {
			defer grpcReqs.WithLabelValues(codes.DeadlineExceeded.String(), serviceName).Inc()
			logger.Warn("error timeout: health rpc did not complete within ", cfg.flRPCTimeout)
			return healthpb.HealthCheckResponse_UNKNOWN, NewGrpcProbeError(StatusRPCFailure, "StatusRPCFailure")
		} else if stat, ok := status.FromError(err); ok && stat.Code() == codes.NotFound {
			defer grpcReqs.WithLabelValues(codes.NotFound.String(), serviceName).Inc()
			// wrap a grpC NOT_FOUND as grpcProbeError.
			// https://github.com/grpc/grpc/blob/master/doc/health-checking.md
			// if the service name is not registerered, the server returns a NOT_FOUND GPRPC status.
			// the Check for a not found should "return nil, status.Error(codes.NotFound, "unknown service")"
			logger.Warn("error Service Not Found ", slog.String("", err.Error()))
			return healthpb.HealthCheckResponse_SERVICE_UNKNOWN, NewGrpcProbeError(StatusServiceNotFound, "StatusServiceNotFound")
		} else {
			defer grpcReqs.WithLabelValues(codes.Unknown.String(), serviceName).Inc()
			logger.Warn("error: health rpc failed: ", slog.String("", err.Error()))
		}
	} else {
		defer grpcReqs.WithLabelValues(resp.GetStatus().String(), serviceName).Inc()
	}
	rpcDuration := time.Since(rpcStart)
	// otherwise, retrurn gRPC-HC status
	logger.Info("time elapsed", slog.Duration("connect", connDuration), slog.Duration("rpc", rpcDuration))

	return resp.GetStatus(), nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {

	var serviceName string
	if cfg.flServiceName != "" {
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
			logger.Error("HealtCheck Probe Error:", slog.String("", pe.Error()))
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
	logger.Info("check ", slog.String("service_name", cfg.flServiceName), slog.String("response", resp.String()))
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
	if cfg.flGrpcTLS {
		creds, err := buildGrpcCredentials()
		if err != nil {
			logger.Error("failed to initialize tls credentials", slog.String("", err.Error()))
			os.Exit(-1)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if cfg.flRunCli {
		resp, err := checkService(context.Background(), cfg.flServiceName)
		if err != nil {
			if pe, ok := err.(*GrpcProbeError); ok {
				logger.Error("HealtCheck Probe Error: ", slog.String("", pe.Error()))
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
		if resp != healthpb.HealthCheckResponse_SERVING {
			logger.Error("HealtCheck Probe Error: service %s failed with reason: %v", cfg.flServiceName, resp.String())
			os.Exit(StatusUnhealthy)
		} else {
			logger.Info("%s %v", cfg.flServiceName, resp.String())
		}
	} else {

		tlsConfig := &tls.Config{}
		if cfg.flHTTPSTLSVerifyClient {
			caCert, err := os.ReadFile(cfg.flHTTPSTLSVerifyCA)
			if err != nil {
				logger.Error("Error reading ca", slog.String("", err.Error()))
				os.Exit(-1)
			}
			caCertPool := x509.NewCertPool()
			if !caCertPool.AppendCertsFromPEM(caCert) {
				logger.Error("Unable to add https server root CA certs")
				os.Exit(-1)
			}

			tlsConfig = &tls.Config{
				ClientCAs:  caCertPool,
				ClientAuth: tls.RequireAndVerifyClientCert,
			}
		}

		r := mux.NewRouter()
		r.Use(prometheusMiddleware)
		r.Path(cfg.flHTTPListenPath).HandlerFunc(healthHandler)

		go func() {
			http.Handle(cfg.flMetricsHTTPPath, promhttp.Handler())
			err := http.ListenAndServe(cfg.flMetricsHTTPListenAddr, nil)
			if err != nil {
				logger.Error("Error starting server ", slog.String("", err.Error()))
				os.Exit(-1)
			}
		}()

		srv := &http.Server{
			Addr:      cfg.flHTTPListenAddr,
			TLSConfig: tlsConfig,
			Handler:   r,
		}

		var err error
		if cfg.flHTTPSTLSServerCert != "" && cfg.flHTTPSTLSServerKey != "" {
			err = srv.ListenAndServeTLS(cfg.flHTTPSTLSServerCert, cfg.flHTTPSTLSServerKey)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil {
			logger.Error("ListenAndServe Error:", slog.String("", err.Error()))
			os.Exit(-1)
		}
	}
}
