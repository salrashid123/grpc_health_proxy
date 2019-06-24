# grpc_health_proxy

`grpc_health_proxy` is a webserver proxy for [gRPC Health Checking Protocol][hc].

This utility sarts up an HTTP/S server which responds back after making an RPC
call to an upstream server's gRPC healthcheck endpoint (`/grpc.health.v1.Health/Check`).

If the healthcheck passes, response back to the original http client will be `200`.  If the
gRPC HealthCheck failed, a `503` is returned.  If the service is not registered, a `404` is returned

Basically, this is an http proxy for the grpc healthcheck protocol.

  `client--->http-->grpc_heatlh_proxy-->gRPC HealthCheck-->gRPC Server`

This utlity uses similar flags, cancellation and timing snippets for the grpc call from [grpc-health-probe](https://github.com/grpc-ecosystem/grpc-health-probe). Use that tool as a specific [Liveness and Readiness Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/) for Kubernetes.  This utility can be used in the same mode but also as a generic HTTP interface (eg, as httpHealthCheck probe)

> This is not an official Google project and is unsupported by Google

**EXAMPLES**

Check the status of an upstream gRPC serviceName `echo.EchoService` listening on `:50051`:

For any mode, enable verbose logging with glog levels: append `--logtostderr=1 -v 10`

- HTTP to gRPC HealthCheck proxy:

`grpc_health_probe` will listen on `:8080` for HTTP healthcheck requests at path `/healthz`.

```text
$ grpc_health_proxy --http-listen-addr localhost:8080 \
                    --http-listen-path /healthz \
                    --grpcaddr localhost:50051 \
                    --service-name echo.EchoServer --logtostderr=1 -v 1
```

```text
curl http://localhost:8080/healthz
```

(also via `&serverName=echo.EchoServer` query parameter)


- HTTPS to gRPC HealthCheck proxy:

`grpc_health_probe` will listen on `:8080` for HTTPS healthcheck requests at path `/healthz`.

HTTPS listener will use keypairs [http_crt.pem, http_key.pem]

```text
$ grpc_health_proxy --http-listen-addr localhost:8080 \
                    --http-listen-path /heatlhz \
                    --grpcaddr localhost:50051 \
                    --https-listen-cert http_crt.pem \
                    --https-listen-key http_key.pem \
                    --grpc-service-name echo.EchoServer --logtostderr=1 -v 1
```

```text
curl --cacert CA_crt.pem  https://localhost:8080/healthz
```

- mTLS HTTPS to gRPC HealthCheck proxy:

`grpc_health_probe` will listen on `:8080` for HTTPS with mTLS healthcheck requests at path `/healthz`.

HTTPS listener will use keypairs [http_crt.pem, http_crt.pem] and verify client certificates issued by `CA_crt.pem`

```text
$ grpc_health_proxy --http-listen-addr localhost:8080 \
                    --http-listen-path /healthz \
                    --grpcaddr localhost:50051 \
                    --https-listen-cert http_crt.pem \
                    --https-listen-key http_key.pem \
                    --grpc-service-name echo.EchoServer \
                    --https-listen-verify \
                    --https-listen-ca=CA_crt.pem --logtostderr=1 -v 1
```

```text
curl --cacert CA_crt.pem --key client_key.pem --cert client_crt.pem  https://localhost:8080/healthz
```

- mTLS to gRPC server from proxy

Options to establish mTLS from the http proxy to gRPC server

- `client->http->grpc_health_proxy->mTLS->gRPC server`

in the example below, `grpc_client_crt.pem` and `grpc_client_key.pem` are the TLS client credentials to present to the gRPC server

```text
$ grpc_health_proxy --http-listen-addr localhost:8080 \
                    --http-listen-path=/healthz \
                    --grpcaddr localhost:50051 \
                    --grpctls \
                    --service-name echo.EchoServer \
                    --grpc-ca-cert=CA_crt.pem \
                    --grpc-client-cert=grpc_client_crt.pem \
                    --grpc-client-key=grpc_client_key.pem \
                    --grpc-sni-server-name=server.domain.com --logtostderr=1 -v 1
```


## Installation

Run this application stand alone or within a Docker image with TLS certificates mounted appropriately.

The [Dockerfile](Dockerfile) provided here run the proxy as an entrypoint and is available at:

  - ```docker.io/salrashid123/grpc_health_proxy```

## Required Options

| Option | Description |
|:------------|-------------|
| **`-http-listen-addr`** | host:port for the http(s) listener |
| **`-grpcaddr`** | upstream gRPC host:port the proxy will connect to |
| **`-service-name`** | gRPC service name to check  |

## gRPC Health Checking Protocol

gRPC server must implement the [gRPC Health Checking Protocol v1][hc]. This means you must to register the
`Health` service and implement the `rpc Check` that returns a `SERVING` status.

The sample gRPC server provided below for golang implements it at:

```golang
func (s *Server) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error)
```

## Health Checking TLS Servers

TLS options for the connection _from_ `grpc_health_proxy` _to_ the upstream gRPC server

| Option | Description |
|:------------|-------------|
| **`-grpctls`** | use TLS to access gRPC Server (default: false) |
| **`-grpc-ca-cert`** | path to file containing CA certificates (to override system root CAs) |
| **`-grpc-client-cert`** | client certificate for authenticating to the server |
| **`-grpc-client-key`** | private key for for authenticating to the server |
| **`-grpc-no-verify`** | use TLS, but do not verify the certificate presented by the server (INSECURE) (default: false) |
| **`-grpc-sni-server-name`** | override the hostname used to verify the server certificate |

## HTTP(s) Proxy

TLS options for the connection from an http client _to_ `grpc_health_proxy`.

Configuration options for HTTTP(s) listener supports TLS and mTLS

| Option | Description |
|:------------|-------------|
| **`-http-listen-addr`** | host:port for the http(s) listener |
| **`-http-listen-path`** | path for http healthcheck requests (defaut `/`|
| **`-https-listen-cert`** | server public certificate for https listner |
| **`-https-listen-key`** | server private key for https listner |
| **`-https-listen-verify`** | option to enable mTLS for HTTPS requests |
| **`-https-listen-ca`** | trust CA for mTLS |

----

[hc]: https://github.com/grpc/grpc/blob/master/doc/health-checking.md


### Example

The `example/` folder contains certificates and sample gRPC server application to test with.  
Add to `/etc/hosts`
`
127.0.0.1   server.domain.com http.domain.com
`

#### TLS Certificates

Sample TLS certificates for use with this sample under `example/app` folder:

- `CA_crt.pem`:  Root CA
- `grpc_server_crt.pem`:  TLS certificate for gRPC server
- `http_server_crt.pem`:  TLS certificate for http listner for `grpc_health_proxy`
- `client_crt.pem`: Client certificate to use while connecting via mTLS to `grpc_health_proxy`

#### grpc Server

To use, first prepare the gRPC server and then run `grpc_health_proxy`.  Use the curl command to invoke the http listener.  Copy the certificates file at `example/certs` to the folder where the `grpc_health_proxy`, the grpc Server and `curl` are run from.

```
cp -R example/app /tmp/
cd /tmp/example/app

export GOPTH=`pwd`
go get golang.org/x/net/context  \
   golang.org/x/net/http2  \
   google.golang.org/grpc  \
   google.golang.org/grpc/credentials \
   google.golang.org/grpc/health/grpc_health_v1
```

#### No TLS

`client->http->grpc_health_proxy->gRPC Server`

  - Run Proxy:
```
  grpc_health_proxy --http-listen-addr localhost:8080 --http-listen-path=/healthz --grpcaddr localhost:50051 --service-name echo.EchoServer  --logtostderr=1 -v 10
```

  - Run gRPC Server
```
  go run src/grpc_server.go --grpcport 0.0.0.0:50051 --insecure
```

  - Invoke http proxy
```
  curl -v http://http.domain.com:8080/healthz
```

#### TLS to Proxy

`client->https->grpc_health_proxy->gRPC Server`

  - Run Proxy:
```
  grpc_health_proxy --http-listen-addr localhost:8080 --http-listen-path=/healthz --grpcaddr localhost:50051 --https-listen-cert=http_server_crt.pem --https-listen-key=http_server_key.pem --service-name echo.EchoServer --logtostderr=1 -v 10
```

  - Run gRPC Server
```
  go run src/grpc_server.go --grpcport 0.0.0.0:50051 --insecure
```

  - Invoke http proxy
```
  curl -v --cacert example/certs/CA_crt.pem  https://http.domain.com:8080/healthz
```

#### mTLS to Proxy and gRPC service

`client->https->grpc_health_proxy->mTLS->gRPC Server`


  - Run Proxy:
```
  grpc_health_proxy --http-listen-addr localhost:8080 --http-listen-path=/healthz --grpcaddr localhost:50051 --https-listen-cert=http_server_crt.pem --https-listen-key=http_server_key.pem --service-name echo.EchoServer --https-listen-verify --https-listen-ca=CA_crt.pem --grpctls --grpc-client-cert=client_crt.pem --grpc-client-key=client_key.pem --grpc-ca-cert=CA_crt.pem --grpc-sni-server-name=server.domain.com --logtostderr=1 -v 10
```

  - Run gRPC Server
```
  go run src/grpc_server.go --grpcport 0.0.0.0:50051 --tlsCert=grpc_server_crt.pem --tlsKey=grpc_server_key.pem
```

  - Invoke http proxy
```
  curl -v --cacert CA_crt.pem  --key client_key.pem --cert client_crt.pem  https://http.domain.com:8080/healthz
```

Or as a docker container from the repo root to mount certs:

```
  docker run  -v `pwd`/example/certs:/certs/ -p 8080:8080 --net=host  -t salrashid123/grpc_health_proxy  --http-listen-addr grpc.domain.com:8080 --http-listen-path=/healthz --grpcaddr localhost:50051 --https-listen-cert=/certs/http_server_crt.pem --https-listen-key=/certs/http_server_key.pem --service-name echo.EchoServer --https-listen-verify --https-listen-ca=/certs/CA_crt.pem --grpctls --grpc-client-cert=/certs/client_crt.pem --grpc-client-key=/certs/client_key.pem --grpc-ca-cert=/certs/CA_crt.pem --grpc-sni-server-name=server.domain.com --logtostderr=1 -v 10
```

### CLI Exit Codes

#### Serving: 0
```
$ ./grpc_health_proxy --runcli --grpcaddr localhost:50051 --service-name echo.EchoServer  --logtostderr=1
echo.EchoServer SERVING

$ echo $?
0
```

#### Connection Failure: 1
```
$ ./grpc_health_proxy --runcli --grpcaddr localhost:50051 --service-name echo.EchoServer  --logtostderr=1 
timeout: failed to connect service localhost:50051 within 1s
HealtCheck Probe Error: StatusConnectionFailure
$ echo $?
1
```

#### Unknown Service: 3
```
$ ./grpc_health_proxy --runcli --grpcaddr localhost:50051 --service-name foo  --logtostderr=1
error Service Not Found rpc error: code = NotFound desc = unknown service
HealtCheck Probe Error: StatusServiceNotFound
$ echo $?
3
```