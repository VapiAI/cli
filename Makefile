.PHONY: build test
build:
	cargo build --locked

test:
	python3 scripts/test-release.py
	python3 scripts/test-packaging.py
	cargo test --locked
	cargo build --release --locked
	python3 scripts/check-performance.py
	python3 scripts/smoke-test.py
