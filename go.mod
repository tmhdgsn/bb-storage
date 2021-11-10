module github.com/buildbarn/bb-storage

go 1.16

replace github.com/gordonklaus/ineffassign => github.com/gordonklaus/ineffassign v0.0.0-20201223204552-cba2d2a1d5d9

require (
	dmitri.shuralyov.com/go/generated v0.0.0-20170818220700-b1254a446363 // indirect
	github.com/aws/aws-sdk-go-v2 v1.11.0
	github.com/aws/aws-sdk-go-v2/config v1.10.0
	github.com/aws/aws-sdk-go-v2/credentials v1.6.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.18.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.9.0
	github.com/bazelbuild/remote-apis v0.0.0-20211004185116-636121a32fa7
	github.com/go-redis/redis/extra/redisotel v0.3.0
	github.com/go-redis/redis/v8 v8.11.4
	github.com/golang/mock v1.6.0
	github.com/google/go-jsonnet v0.17.0
	github.com/google/uuid v1.3.0
	github.com/gordonklaus/ineffassign v0.0.0-20210914165742-4cc7213b9bc8
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/jmespath/go-jmespath v0.4.0
	github.com/klauspost/compress v1.13.6
	github.com/lazybeaver/xorshift v0.0.0-20170702203709-ce511d4823dd
	github.com/prometheus/client_golang v1.11.0
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.26.1
	go.opentelemetry.io/contrib/propagators/b3 v1.1.1
	go.opentelemetry.io/otel v1.1.0
	go.opentelemetry.io/otel/exporters/jaeger v1.1.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.1.0
	go.opentelemetry.io/otel/sdk v1.1.0
	go.opentelemetry.io/otel/trace v1.1.0
	go.opentelemetry.io/proto/otlp v0.10.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20211109184856-51b60fd695b3
	google.golang.org/genproto v0.0.0-20211104193956-4c6863e31247
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
	mvdan.cc/gofumpt v0.1.1
)
