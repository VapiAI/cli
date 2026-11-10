.PHONY: build test generate preview-publish
build:
	cargo build --locked

test:
	cargo test --locked
	python3 scripts/smoke-test.py

generate:
	python3 scripts/generate.py

# Preview generates publishing files locally; it never releases a package.
preview-publish:
	fern generate --group publish --version $$(python3 -c 'import tomllib; print(tomllib.load(open("Cargo.toml", "rb"))["package"]["version"])') --preview --output generated-publish
