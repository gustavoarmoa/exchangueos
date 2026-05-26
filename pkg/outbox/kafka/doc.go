// Package kafka — outbox.Publisher backed by franz-go (github.com/twmb/franz-go/pkg/kgo).
//
// Guarded by `//go:build kafka` — the franz-go dependency is heavy enough that
// we want the default build to stay slim. Enable via:
//
//	go build -tags kafka ./...
//
// Production-tuned defaults:
//
//	RequiredAcks:        acks=all (LeaderAndISR commit)
//	Compression:         zstd
//	Idempotent producer: enabled (prevents duplicates on broker retry)
//	Max in-flight:       1 (with idempotent producer for strict ordering per partition)
//	Retries:             effectively infinite (delivery is guaranteed by outbox)
package kafka
