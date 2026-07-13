VERSION ?= 0.0.0-dev
COMMIT ?= $(shell git rev-parse --short=12 HEAD)
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS = -X github.com/vvitem/item_all/internal/buildinfo.Version=$(VERSION) -X github.com/vvitem/item_all/internal/buildinfo.Commit=$(COMMIT) -X github.com/vvitem/item_all/internal/buildinfo.BuildTime=$(BUILD_TIME)

.PHONY: test frontend-test frontend-build dev build

test:
	go test ./...
	cd frontend && pnpm typecheck && pnpm lint && pnpm test:run

frontend-test:
	cd frontend && pnpm test:run

frontend-build:
	cd frontend && pnpm build

dev:
	wails dev

build:
	wails build -clean -trimpath -ldflags "$(LDFLAGS)"
