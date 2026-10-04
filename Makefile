.PHONY: build dev build-web build-go test vet release-local clean

VERSION ?= v0.1.0

build: build-web build-go

build-web:
	cd web && npm ci && npm run build

build-go:
	go build -ldflags "-X easytalk/internal/buildinfo.Version=$(VERSION)" -o easytalk ./cmd/easytalk

dev:
	go run ./cmd/easytalk

test:
	go test ./...

vet:
	go vet ./...

release-local:
	./scripts/build-release.sh

clean:
	rm -f easytalk
	rm -rf release