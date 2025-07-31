module github.com/xgmsx/rsf/inventory

go 1.24.2

replace github.com/xgmsx/rsf/shared => ../shared

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.1
	github.com/xgmsx/rsf/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250728155136-f173205681a0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250721164621-a45f3dfb1074 // indirect
)
