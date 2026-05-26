# ExchangeOS — Makefile (delegacao para Task — taskfile.dev é o primary runner).
# Existe para compatibilidade Unix tradicional e onboarding mais rapido.
# PowerShell mirror em scripts/exchangeos.ps1.

.PHONY: help install tidy proto build test lint fmt vet sec db compose clean dash \
        proto-gen proto-lint proto-breaking \
        build-api build-worker build-migrator build-cls-cycle build-eod build-mq-bridge build-cred-rotator \
        test-unit test-integration test-e2e test-cover \
        sec-secrets sec-trivy sec-govulncheck sec-cosign \
        db-up db-migrate db-seed db-reset \
        docker-build compose-up compose-down compose-logs \
        otel-up dash-update \
        hooks-pre-commit hooks-pre-push hooks-pre-merge

TASK := task

help:
	@$(TASK) --list-all

install: ; @$(TASK) install
tidy: ; @$(TASK) tidy

proto: proto-gen
proto-gen: ; @$(TASK) proto:gen
proto-lint: ; @$(TASK) proto:lint
proto-breaking: ; @$(TASK) proto:breaking

build: ; @$(TASK) build
build-api: ; @$(TASK) build:api
build-worker: ; @$(TASK) build:worker
build-migrator: ; @$(TASK) build:migrator
build-cls-cycle: ; @$(TASK) build:cls-cycle
build-eod: ; @$(TASK) build:eod
build-mq-bridge: ; @$(TASK) build:mq-bridge
build-cred-rotator: ; @$(TASK) build:cred-rotator

test: ; @$(TASK) test
test-unit: ; @$(TASK) test:unit
test-integration: ; @$(TASK) test:integration
test-e2e: ; @$(TASK) test:e2e
test-cover: ; @$(TASK) test:cover

lint: ; @$(TASK) lint
fmt: ; @$(TASK) fmt
vet: ; @$(TASK) vet

sec: sec-secrets sec-trivy sec-govulncheck
sec-secrets: ; @$(TASK) sec:secrets
sec-trivy: ; @$(TASK) sec:trivy
sec-govulncheck: ; @$(TASK) sec:govulncheck
sec-cosign: ; @$(TASK) sec:cosign

db-up: ; @$(TASK) db:up
db-migrate: ; @$(TASK) db:migrate
db-seed: ; @$(TASK) db:seed
db-reset: ; @$(TASK) db:reset

docker-build: ; @$(TASK) docker:build
compose-up: ; @$(TASK) compose:up
compose-down: ; @$(TASK) compose:down
compose-logs: ; @$(TASK) compose:logs

otel-up: ; @$(TASK) otel:up

xsd-download: ; @$(TASK) xsd:download
xsd-verify: ; @$(TASK) xsd:verify

dash: dash-update
dash-update: ; @$(TASK) dash-update

hooks-pre-commit: ; @$(TASK) hooks:pre-commit
hooks-pre-push: ; @$(TASK) hooks:pre-push
hooks-pre-merge: ; @$(TASK) hooks:pre-merge

clean: ; @$(TASK) clean
