module github.com/xgmsx/rsf/payment

go 1.24.2

replace github.com/xgmsx/rsf/shared => ../shared

require (
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.1
	google.golang.org/grpc v1.74.2
)

require (
	go.opentelemetry.io/otel v1.37.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250728155136-f173205681a0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250721164621-a45f3dfb1074 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
