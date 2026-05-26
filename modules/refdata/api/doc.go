// Package api hosts the gRPC adapter for the refdata bounded context.
//
// Files in this package are guarded by `//go:build grpcgen`. The default `go build`
// excludes them so the repository compiles without `proto/gen/` materialised. To
// enable:
//
//	task proto:gen           # produces proto/gen/exchangeos/v1/*.pb.go
//	go build -tags grpcgen ./...
//
// Once enabled, the adapter translates pb.GetCurrencyRequest → application.GetCurrency
// (and so on) and registers with the gRPC server via cmd/api/grpc_register_proto.go.
package api
