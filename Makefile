.PHONY: build test check clean

build:
	mkdir -p dist
	go build -trimpath -o dist/spotify-tui ./cmd/spotify-tui
	cargo build --manifest-path playback/Cargo.toml --locked --release
	install -m 755 playback/target/release/spotify-tui-playback dist/.spotify-tui-playback.new
	mv -f dist/.spotify-tui-playback.new dist/spotify-tui-playback

test:
	go test -race ./...
	cargo test --manifest-path playback/Cargo.toml --locked

check:
	test -z "$$(gofmt -l cmd internal)"
	go vet ./...
	cargo fmt --manifest-path playback/Cargo.toml --check
	cargo clippy --manifest-path playback/Cargo.toml --locked --all-targets -- -D warnings

clean:
	rm -rf dist
