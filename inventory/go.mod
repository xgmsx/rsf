module github.com/xgmsx/rsf/inventory

go 1.24.2

replace github.com/xgmsx/rsf/shared => ../shared

require (
	github.com/brianvoe/gofakeit/v7 v7.4.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.1
	github.com/stretchr/testify v1.10.0
	github.com/xgmsx/rsf/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250728155136-f173205681a0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250721164621-a45f3dfb1074 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
