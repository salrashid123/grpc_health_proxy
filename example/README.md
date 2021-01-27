## Sample gRPC Client/Server


### Insecure

```bash
# Server
docker run --net=host -p 50051:50051 \
  -t gcr.io/cloud-solutions-images/grpc_app /grpc_server \
  --grpcport :50051 --insecure

# Client
docker run --net=host \
  -t gcr.io/cloud-solutions-images/grpc_app /grpc_client \
  --host localhost:50051 --insecure
```

### Server TLS

```bash
# Server
docker run --net=host -p 50051:50051 \
   -t gcr.io/cloud-solutions-images/grpc_app /grpc_server \
      --grpcport :50051  --tlsCert=certs/grpc_server_crt.pem \
      --tlsKey=/certs/grpc_server_key.pem

# Client
docker run --net=host -p 50051:50051 \
   -t gcr.io/cloud-solutions-images/grpc_app /grpc_client \
      --host localhost:50051 --tlsCert certs/CA_crt.pem \
      -skipHealthCheck --servername grpc.domain.com
```