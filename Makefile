.PHONY: build dev build-web build-go clean

VERSION ?= v0.1.0

build: build-web build-go

build-web:
	cd web && npm install && npm run build

build-go:
	go build -o easytalk ./cmd/easytalk

dev:
	go run ./cmd/easytalk

clean:
	rm -f easytalk