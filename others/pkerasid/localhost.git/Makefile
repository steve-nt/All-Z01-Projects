CARGO ?= cargo
BIN   := localhost
CONFIG ?= config/default.toml
SIEGE_TARGET ?= http://127.0.0.1:8080/
SIEGE_TIME   ?= 60S

.PHONY: all build run release test lint fmt fmt-check clippy clean siege check

all: lint test build

build:
	$(CARGO) build

release:
	$(CARGO) build --release

run: build
	$(CARGO) run -- $(CONFIG)

test:
	$(CARGO) test --all-features

fmt:
	$(CARGO) fmt --all

fmt-check:
	$(CARGO) fmt --all -- --check

clippy:
	$(CARGO) clippy --all-targets --all-features -- -D warnings

lint: fmt-check clippy

check:
	$(CARGO) check --all-targets --all-features

siege:
	siege -b -t$(SIEGE_TIME) $(SIEGE_TARGET)

clean:
	$(CARGO) clean
