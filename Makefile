.PHONY: lint
lint:
	golangci-lint run --disable=typecheck

dev:
	air
