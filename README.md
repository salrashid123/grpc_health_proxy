# grpc_health_proxy

`grpc_health_proxy` is a webserver proxy for [gRPC Health Checking Protocol][hc].

This utility starts up an HTTP/S server which responds back after making an RPC
call to an upstream server's gRPC healthcheck endpoint (`/grpc.health.v1.Health/Check`).

If the healthcheck passes, response back to the original http client will be `200`.  If the
gRPC HealthCheck failed, a `503` is returned.  If the service is not registered, a `404` is returned

Basically, this is an http proxy for the grpc healthcheck protocol.

  `client--->TLS-->grpc_heatlh_proxy *gRPC HealthCheck*-->TLS-->gRPC Server`

This utility uses similar flags, cancellation and timing snippets for the grpc call from [grpc-health-probe](https://github.com/grpc-ecosystem/grpc-health-probe). Use that tool as a specific [Liveness and Readiness Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/) for Kubernetes.  This utility can be used in the same cli mode but also as a generic HTTP interface (eg, as httpHealthCheck probe).  For more information on the CLI mode without http listener, see the section at the end.

> This is not an official Google project and is unsupported by Google

### Build

You can either build from source

```bash
go build -o grpc_health_proxy main.go
```

or use one of the binaries in `Releases` section or the docker image

## Quickstart

The following in the `example/` folder checks the status of an upstream gRPC serviceName `echo.EchoService` listening on `:50051`:

`client->http->grpc_health_proxy->gRPC Server`

  - Run gRPC Server

```bash
go run src/grpc_server.go \
    --grpcport 0.0.0.0:50051 \
    --insecure
```

- Run Proxy:

you can either build the proxy from source or use the binary in releases page:

```bash
grpc_health_proxy \
    --http-listen-addr localhost:8080 \
    --http-listen-path=/healthz \
    --grpcaddr localhost:50051 \
    --service-name echo.EchoServer \
    --logtostderr=1 -v 10
```

  - Invoke http proxy

```bash
  curl -v \
    --resolve 'http.domain.com:8080:127.0.0.1' \
    http://http.domain.com:8080/healthz
```

---

## Installation

Run this application stand alone or within a Docker image with TLS certificates mounted appropriately.

#### as docker

The [Dockerfile](Dockerfile) and Bazel build directives are provided here for you to build your own image.

  - [docker.io/salrashid123/grpc_health_proxy](https://hub.docker.com/repository/docker/salrashid123/grpc_health_proxy/general)

    **NOTE:** the default docker image listens on containerPort `:8080`

#### from source

To compile the proxy directly, run

```bash
go build -o grpc_health_proxy main.go
```

#### from binary

Download a binary from the Release page.

The proxy version also correspond to docker image version tags (eg `docker.io/salrashid123/grpc_health_proxy:1.1.0`)

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
| **`-grpc-tls-no-verify`** | use TLS, but do not verify the certificate presented by the server (INSECURE) (default: false) |
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

#### TLS Certificates

Sample TLS certificates for use with this sample under `example/certs` folder:

- `CA_crt.pem`:  Root CA
- `grpc_server_crt.pem`:  TLS certificate for gRPC server
- `http_server_crt.pem`:  TLS certificate for http listner for `grpc_health_proxy`
- `client_crt.pem`: Client certificate to use while connecting via mTLS to `grpc_health_proxy`

#### grpc Server

To use, first prepare the gRPC server and then run `grpc_health_proxy`.  Use the curl command to invoke the http listener.  Copy the certificates file at `example/certs` to the folder where the `grpc_health_proxy`, the grpc Server and `curl` are run from.

---

#### No TLS

`client->http->grpc_health_proxy->gRPC Server`

  - Run Proxy:

```bash
grpc_health_proxy \
    --http-listen-addr localhost:8080 \
    --http-listen-path=/healthz \
    --grpcaddr localhost:50051 \
    --service-name echo.EchoServer \
    --logtostderr=1 -v 10


I0329 21:05:37.969302 1128850 main.go:229] Running HealthCheck for service:  echo.EchoServer
I0329 21:05:37.970145 1128850 main.go:253] time elapsed: connect=763.1µs rpc=848.974µs
I0329 21:05:37.970289 1128850 main.go:291] echo.EchoServer SERVING
```

  - Run gRPC Server

```bash
go run src/grpc_server.go \
    --grpcport 0.0.0.0:50051 \
    --insecure
```

  - Invoke http proxy

```bash
$ curl -v     --resolve 'http.domain.com:8080:127.0.0.1'     http://http.domain.com:8080/healthz

echo.EchoServer SERVING
```

#### TLS to Proxy

`client->TLS->grpc_health_proxy->gRPC Server`

  - Run Proxy:

```bash
grpc_health_proxy \
    --http-listen-addr localhost:8080 \
    --http-listen-path=/healthz \
    --grpcaddr localhost:50051 \
    --https-listen-cert=certs/http_server_crt.pem \
    --https-listen-key=certs/http_server_key.pem \
    --service-name echo.EchoServer \
    --logtostderr=1 -v 10
```

  - Run gRPC Server

```bash
go run src/grpc_server.go --grpcport 0.0.0.0:50051 --insecure
```

  - Invoke http proxy

```bash
curl -v \
    --cacert certs/CA_crt.pem  \
    --resolve 'http.domain.com:8080:127.0.0.1' \
    https://http.domain.com:8080/healthz
```

#### TLS to Proxy and TLS gRPC service

`client->TLS->grpc_health_proxy->TLS->gRPC Server`

Note that for convenience, we are reusing the same client and CA certificate during various stages here:

  - Run Proxy:

```bash
grpc_health_proxy \
    --http-listen-addr localhost:8080 \
    --http-listen-path=/healthz \
    --grpcaddr localhost:50051 \
    --https-listen-cert=certs/http_server_crt.pem \
    --https-listen-key=certs/http_server_key.pem \
    --service-name echo.EchoServer \
    --https-listen-ca=certs/CA_crt.pem \
    --grpctls \
    --grpc-ca-cert=certs/CA_crt.pem \
    --grpc-sni-server-name=grpc.domain.com \
    --logtostderr=1 -v 10
```

  - Run gRPC Server

```bash
go run src/grpc_server.go \
    --grpcport 0.0.0.0:50051 \
    --tlsCert=certs/grpc_server_crt.pem \
    --tlsKey=certs/grpc_server_key.pem
```

  - Invoke http proxy

```bash
curl -v \
   --resolve 'http.domain.com:8080:127.0.0.1' \
   --cacert certs/CA_crt.pem \
   https://http.domain.com:8080/healthz
```

#### mTLS to Proxy and mTLS to gRPC service

`client->mTLS->grpc_health_proxy->mTLS->gRPC Server`

Note that for convenience, we are reusing the same client and CA certificate during various stages here:

  - Run Proxy:

```bash
grpc_health_proxy \
    --http-listen-addr localhost:8080 \
    --http-listen-path=/healthz \
    --grpcaddr localhost:50051 \
    --https-listen-cert=certs/http_server_crt.pem \
    --https-listen-key=certs/http_server_key.pem \
    --service-name echo.EchoServer \
    --https-listen-verify \
    --https-listen-ca=certs/CA_crt.pem \
    --grpctls \
    --grpc-client-cert=certs/proxy_client_crt.pem \
    --grpc-client-key=certs/proxy_client_key.pem \
    --grpc-ca-cert=certs/CA_crt.pem \
    --grpc-sni-server-name=grpc.domain.com \
    --logtostderr=1 -v 10
```

  - Run gRPC Server

```bash
go run src/grpc_server.go \
    --grpcport 0.0.0.0:50051 \
    --backendMTLS \
    --mtlsBackendCA=certs/CA_crt.pem \
    --tlsCert=certs/grpc_server_crt.pem \
    --tlsKey=certs/grpc_server_key.pem
```

  - Invoke http proxy

```bash
curl -v \
   --resolve 'http.domain.com:8080:127.0.0.1' \
   --cacert certs/CA_crt.pem \
   --key certs/client_key.pem \
   --cert certs/client_crt.pem \
   https://http.domain.com:8080/healthz
```

### Kubernetes Pod Healthcheck

You can use this utility as a proxy for service healthchecks.

This is useful for external services that utilize HTTP but need to verify a gRPC services health status.

In the kubernetes deployment below, an http request to the healthcheck serving port (`:8080`) will reflect
the status of the gRPC service listening on port `:50051`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp-deployment
  labels:
    type: myapp-deployment-label
spec:
  replicas: 1
  selector:
    matchLabels:
      type: myapp
  template:
    metadata:
      labels:
        type: myapp
    spec:
      containers:
      - name: hc-proxy
        image: docker.io/salrashid123/grpc_health_proxy
        args: [
          "--http-listen-addr=0.0.0.0:8080",
          "--grpcaddr=localhost:50051",
          "--service-name=echo.EchoServer",
          "--logtostderr=1",
          "-v=1"
        ]
        ports:
        - containerPort: 8080
      - name: grpc-app
        image: docker.io/salrashid123/grpc_only_backend
        args: [
          "/grpc_server",
          "--grpcport=0.0.0.0:50051",
          "--insecure"
        ]
        ports:
        - containerPort: 50051
```

The docker image used for the gRPC Server is taken from `example/` folder in the same repo

---

### CLI Exit Codes

You can run tis utiity is cli mode directly similar to the `grpc_health_probe` cited above.  In this cli mode, you can a grpc healthcheck service without resorting to curl, etc.

There are several exit codes this utility returns

- 0: Serving

```bash
$ ./grpc_health_proxy \
   --runcli \
   --grpcaddr localhost:50051 \
   --service-name echo.EchoServer  \
   --logtostderr=1

echo.EchoServer SERVING

$ echo $?
0
```

- 5: Unhealthy

```bash
$ ./grpc_health_proxy \
   --runcli \
   --grpcaddr localhost:50051 \
   --service-name echo.EchoServer  \
   --logtostderr=1

echo.EchoServer UNHEALTHY

$ echo $?
5
```

- 1: Connection Failure

```bash
$ ./grpc_health_proxy \
   --runcli \
   --grpcaddr localhost:50051 \
   --service-name echo.EchoServer \
   --logtostderr=1 

timeout: failed to connect service localhost:50051 within 1s
HealtCheck Probe Error: StatusConnectionFailure

$ echo $?
1
```

- 3: Unknown Service

```bash
$ ./grpc_health_proxy \
   --runcli \
   --grpcaddr localhost:50051 \
   --service-name foo  \
   --logtostderr=1

error Service Not Found rpc error: code = NotFound desc = unknown service
HealtCheck Probe Error: StatusServiceNotFound

$ echo $?
3
```

#### Verify Release Binary

If you download a binary from the "Releases" page, you can verify the signature with GPG:

```bash
gpg --keyserver keyserver.ubuntu.com --recv-keys 5D8EA7261718FE5728BA937C97341836616BF511

## to verify the checksum file for a given release:
wget https://github.com/salrashid123/grpc_health_proxy/releases/download/v1.1.0/grpc_health_proxy_1.1.0_checksums.txt
wget https://github.com/salrashid123/grpc_health_proxy/releases/download/v1.1.0/grpc_health_proxy_1.1.0_checksums.txt.sig

gpg --verify grpc_health_proxy_1.1.0_checksums.txt.sig grpc_health_proxy_1.1.0_checksums.txt
```

#### Verify Container Image Signature

The images are also signed using my github address (`salrashid123@gmail`).  If you really want to, you can verify each signature usign `cosign`:

```bash
## for tag/version  index.docker.io/salrashid123/grpc_health_proxy:1.1.0:
IMAGE="index.docker.io/salrashid123/grpc_health_proxy@sha256:4e65829ecf0a523c4d3a16ac98c19d3476c867dc5bd78f0d4ef429907bc18572"

## i signed it directly, keyless:
# $ cosign sign $IMAGE

## which you can verify:
$ cosign verify --certificate-identity=salrashid123@gmail.com  --certificate-oidc-issuer=https://github.com/login/oauth $IMAGE | jq '.'

## search and get 
# $ rekor-cli search --rekor_server https://rekor.sigstore.dev  --email salrashid123@gmail.com
# $ rekor-cli get --rekor_server https://rekor.sigstore.dev  --log-index $LogIndex  --format=json | jq '.'
```

These images were built using bazel so you should get the same container hash (i.e., deterministic builds)