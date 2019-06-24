FROM golang:1.11 AS build
ENV PROJECT grpc_health_proxy
WORKDIR /src/$PROJECT
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go install -a -tags netgo -ldflags=-w

FROM gcr.io/distroless/base
COPY --from=build /go/bin/grpc_health_proxy /bin/grpc_health_proxy
EXPOSE 8080
ENTRYPOINT [ "/bin/grpc_health_proxy" ]
