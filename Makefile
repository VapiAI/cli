.PHONY: build test generate
build:
	cargo build --locked

test:
	cargo test --locked
	python3 scripts/smoke-test.py
	cd npm-wrapper && npm test
	python3 scripts/check-version.py

generate:
	python3 scripts/generate.py
