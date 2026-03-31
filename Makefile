GO ?= go
NPM ?= npm
NPX ?= npx

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags="-s -w -X main.version=$(VERSION)"

.PHONY: build build-ui dev-backend dev-frontend test docker-build clean

build-ui:
	cd ui && $(NPM) install && $(NPX) vite build

build: build-ui
	$(GO) build $(LDFLAGS) -o bin/nats-dashboard ./cmd/nats-dashboard

dev-backend:
	$(GO) run ./cmd/nats-dashboard -config config.yaml

dev-frontend:
	cd ui && $(NPX) vite

test:
	$(GO) test -count=1 -timeout 120s github.com/machmqtt/nats-dashboard/...

docker-build:
	docker build -t nats-dashboard .

clean:
	rm -rf bin/ internal/api/dist/assets/ ui/node_modules/
