.PHONY: build test
build:
	cargo build --locked

test:
	cargo test --locked
	cargo build --release --locked
	python3 scripts/check-performance.py
	python3 scripts/smoke-test.py
