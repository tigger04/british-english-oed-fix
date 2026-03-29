# ABOUTME: Build, test, install, and release targets for the sanitize CLI tool.
# Standard entry points so any user can run make build/test/install without
# needing to know Go toolchain details.

BINARY := sanitize
BUILD_DIR := .
GO := go

.PHONY: build test test-one-off install clean release

build:
	$(GO) build -o $(BUILD_DIR)/$(BINARY) ./cmd/sanitize/

test: build
	$(GO) test ./pkg/spelling/ -v
	$(GO) test ./tests/regression/ -v

test-one-off:
ifdef ISSUE
	$(GO) test ./tests/one_off/ -v -run "$(ISSUE)"
else
	$(GO) test ./tests/one_off/ -v
endif

install: build
	cp $(BUILD_DIR)/$(BINARY) ~/bin/$(BINARY)

clean:
	rm -f $(BUILD_DIR)/$(BINARY)

sync:
	git add --all
	git commit -m "chore: sync"
	git pull
	git push
