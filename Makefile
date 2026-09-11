.PHONY: build install test clean api-check

SPEC_URL ?= https://www.mobileops.at/openapi.yaml
SPEC_FILE ?= /tmp/mobileops-openapi.yaml

build:
	go build -o bin/mobileops ./cmd/mobileops

install:
	go install ./cmd/mobileops

test:
	go test ./...

api-check:
	@if [ ! -f "$(SPEC_FILE)" ] || [ "$(SPEC_FILE)" = "/tmp/mobileops-openapi.yaml" ]; then \
		echo "Downloading $(SPEC_URL) -> $(SPEC_FILE)"; \
		curl -fsSL "$(SPEC_URL)" -o "$(SPEC_FILE)"; \
	fi
	MOBILEOPS_OPENAPI_SPEC="$(SPEC_FILE)" go test ./internal/commands -run 'TestAPI' -v

clean:
	rm -rf bin/
