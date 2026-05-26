// Package api hosts the gRPC adapter for the quote bounded context.
//
// Files in this package are guarded by `//go:build grpcgen`. The default `go build`
// excludes them so the repository compiles without `proto/gen/` materialised.
//
// Workflow:
//
//	task proto:gen           # produces proto/gen/exchangeos/v1/*.pb.go
//	go build -tags grpcgen ./...
package api
