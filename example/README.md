## Sample gRPC Client/Server


### Insecure

```bash
cd example/
# Server
docker run --net=host -p 50051:50051 \
  -t docker.io/salrashid123/grpc_app /grpc_server \
  --grpcport :50051 --insecure

# Client
docker run --net=host \
  -t docker.io/salrashid123/grpc_app /grpc_client \
  --host localhost:50051 --insecure
```

### Server TLS

```bash
cd example/
# Server
go run src/grpc_server.go  --grpcport :50051  --tlsCert=certs/grpc_server_crt.pem       --tlsKey=certs/grpc_server_key.pem

# docker run --net=host -p 50051:50051 \
#    -t docker.io/salrashid123/grpc_app /grpc_server \
#       --grpcport :50051  --tlsCert=certs/grpc_server_crt.pem \
#       --tlsKey=/certs/grpc_server_key.pem
```

# Client

```bash
cd example/
go run src/grpc_client.go --host localhost:50051 --tlsCert certs/CA_crt.pem \
      -skipHealthCheck --servername grpc.domain.com

# docker run --net=host -p 50051:50051 \
#    -t docker.io/salrashid123/grpc_app /grpc_client \
#       --host localhost:50051 --tlsCert certs/CA_crt.pem \
#       -skipHealthCheck --servername grpc.domain.com
```

