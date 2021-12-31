# grpc_health_proxy

`grpc_health_proxy` is a webserver proxy for [gRPC Health Checking Protocol][hc].

This utility starts up an HTTP/S server which responds back after making an RPC
call to an upstream server's gRPC healthcheck endpoint (`/grpc.health.v1.Health/Check`).

If the healthcheck passes, response back to the original http client will be `200`.  If the
gRPC HealthCheck failed, a `503` is returned.  If the service is not registered, a `404` is returned

Basically, this is an http proxy for the grpc healthcheck protocol.

  `client--->http-->grpc_heatlh_proxy-->gRPC HealthCheck-->gRPC Server`

This utility uses similar flags, cancellation and timing snippets for the grpc call from [grpc-health-probe](https://github.com/grpc-ecosystem/grpc-health-probe). Use that tool as a specific [Liveness and Readiness Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/) for Kubernetes.  This utility can be used in the same cli mode but also as a generic HTTP interface (eg, as httpHealthCheck probe).  For more information on the CLI mode without http listener, see the section at the end.

> This is not an official Google project and is unsupported by Google

**EXAMPLES**

Check the status of an upstream gRPC serviceName `echo.EchoService` listening on `:50051`:

For any mode, enable verbose logging with glog levels: append `--logtostderr=1 -v 10`

### HTTP to gRPC HealthCheck proxy:

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

---

### HTTPS to gRPC HealthCheck proxy:

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
curl \
  --cacert CA_crt.pem \
  --resolve 'http.domain.com:8080:127.0.0.1' \
  https://host.domain.com:8080/healthz
```

---

### mTLS HTTPS to gRPC HealthCheck proxy:

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
curl \
  --cacert CA_crt.pem \
  --key client_key.pem \
  --cert client_crt.pem \
  --resolve 'http.domain.com:8080:127.0.0.1'
  https://host.domain.com:8080/healthz
```

---

### mTLS to gRPC server from proxy

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

The [Dockerfile](Dockerfile) provided here for the proxy but you are _strongly_ encouraged to deploy your own
docker image of the same:

  - ```docker.io/salrashid123/grpc_health_proxy```
    **NOTE:** the default docker image listens on containerPort `:8080`


>> Note, the application an healthcheck docker image now is on `docker.io`.  Please generate and host your own images as necessary:

* `gcr.io/cloud-solutions-images/grpc_health_proxy` =-> `docker.io/salrashid123/grpc_health_proxy`

To compile the proxy directly, run

```
go build -o grpc_health_proxy main.go
```

or download a binary from the Release page.

The proxy version also correspond to docker image tags.
-  `docker.io/salrashid123/grpc_health_proxy:1.0.0` `sha256:bba655892eedd2a59a0197f0949faad24f49546a6e548489be545c56776abbf9`)

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
```
  cd example/
  grpc_health_proxy \
    --http-listen-addr localhost:8080 \
    --http-listen-path=/healthz \
    --grpcaddr localhost:50051 \
    --service-name echo.EchoServer \
    --logtostderr=1 -v 10
```

  - Run gRPC Server
```
  go run src/grpc_server.go \
    --grpcport 0.0.0.0:50051 \
    --insecure
```

  - Invoke http proxy
```
  curl -v \
    --resolve 'http.domain.com:8080:127.0.0.1' \
    http://http.domain.com:8080/healthz
```

---

#### TLS to Proxy

`client->https->grpc_health_proxy->gRPC Server`

  - Run Proxy:
```bash
  cd example/
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
```
  go run src/grpc_server.go --grpcport 0.0.0.0:50051 --insecure
```

  - Invoke http proxy
```
  curl -v \
    --cacert certs/CA_crt.pem  \
    --resolve 'http.domain.com:8080:127.0.0.1' \
    https://http.domain.com:8080/healthz
```

---

#### mTLS to Proxy and gRPC service

`client->https->grpc_health_proxy->mTLS->gRPC Server`

Note that for convenience, we are reusing the same client and CA certificate during various stages here:

  - Run Proxy:
```bash
  cd example/
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
    --grpc-client-cert=certs/client_crt.pem \
    --grpc-client-key=certs/client_key.pem \
    --grpc-ca-cert=certs/CA_crt.pem \
    --grpc-sni-server-name=grpc.domain.com \
    --logtostderr=1 -v 10
```

  - Run gRPC Server
```
  go run src/grpc_server.go \
    --grpcport 0.0.0.0:50051 \
    --tlsCert=certs/grpc_server_crt.pem \
    --tlsKey=certs/grpc_server_key.pem
```

  - Invoke http proxy
```
  curl -v \
   --resolve 'http.domain.com:8080:127.0.0.1' \
   --cacert certs/CA_crt.pem \
   --key certs/client_key.pem \
   --cert certs/client_crt.pem \
   https://http.domain.com:8080/healthz
```

Or as a docker container from the repo root to mount certs:

```
  docker run  -v `pwd`/certs:/certs/ \
    -p 8080:8080 \
    --net=host  \
    -t docker.io/salrashid123/grpc_health_proxy \
    --http-listen-addr localhost:8080 \
    --http-listen-path=/healthz \
    --grpcaddr localhost:50051 \
    --https-listen-cert=/certs/http_server_crt.pem \
    --https-listen-key=/certs/http_server_key.pem \
    --service-name echo.EchoServer \
    --https-listen-verify \
    --https-listen-ca=/certs/CA_crt.pem \
    --grpctls \
    --grpc-client-cert=/certs/client_crt.pem \
    --grpc-client-key=/certs/client_key.pem \
    --grpc-ca-cert=/certs/CA_crt.pem \
    --grpc-sni-server-name=grpc.domain.com \
    --logtostderr=1 -v 10
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

The docker image used for the gRPC Server is taken from `eample/` folder in the same repo

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

```
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