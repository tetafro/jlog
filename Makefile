.PHONY: build
build:
	cargo build --release

.PHONY: install
install:
	install target/release/jlog $(HOME)/.cargo/bin

.PHONY: publish
publish:
	cargo publish
