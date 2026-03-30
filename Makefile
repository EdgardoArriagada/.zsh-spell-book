
# Allow `make go-build <name>` or `make go-dev <name>` to skip fzf selection
ifneq ($(filter $(firstword $(MAKECMDGOALS)),go-build go-dev rust-build rust-dev),)
  TARGET := $(word 2,$(MAKECMDGOALS))
  $(if $(TARGET),$(eval $(TARGET):;@:))
endif

define go-build-dev-fn
	go build -o ./bin/$(1) ./cmd/$(1)
endef

define go-build-fn
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o ./bin/$(1) ./cmd/$(1) && echo "$(1) ✅"
endef

define go-dev-fn
	ls ./{lib,cmd}/**/*.go | entr -c $(call go-build-dev-fn,$(1))
endef

.go-run:
	@if [[ -z "$(TARGET)" ]]; then \
		TARGET=$$(ls ./src-go/cmd | fzf); \
	fi; \
	if [[ -n "$$TARGET" ]]; then \
		(cd src-go && $(call $(ACTION),$$TARGET)); \
	fi

go-build:
	$(MAKE) ACTION='go-build-fn' TARGET='$(TARGET)' .go-run

go-dev:
	$(MAKE) ACTION='go-dev-fn' TARGET='$(TARGET)' .go-run

zsh-dev:
	@find . -name "*.zsh" -not -path "./src-go/*" | entr -c ./src-rust/bin/zsb_bundle

go-build-all:
	for target in $$(ls ./src-go/cmd); do \
		$(MAKE) TARGET=$$target ACTION='go-build-fn' .go-run; \
	done

go-test-all:
	cd src-go && go list -f '{{.Dir}}' -m | xargs -I{} sh -c 'cd "{}" && go test ./...'

# Rust targets

define rust-build-dev-fn
	cargo build -p $(1) && cp ./target/debug/$(1) ./bin/$(1)
endef

define rust-build-fn
	cargo build --release -p $(1) && cp ./target/release/$(1) ./bin/$(1) && echo "$(1) ✅"
endef

define rust-dev-fn
	find ./cmd -name "*.rs" | entr -c $(call rust-build-dev-fn,$(1))
endef

.rust-run:
	@if [[ -z "$(TARGET)" ]]; then \
		TARGET=$$(ls ./src-rust/cmd | fzf); \
	fi; \
	if [[ -n "$$TARGET" ]]; then \
		(cd src-rust && $(call $(ACTION),$$TARGET)); \
	fi

rust-build:
	$(MAKE) ACTION='rust-build-fn' TARGET='$(TARGET)' .rust-run

rust-dev:
	$(MAKE) ACTION='rust-dev-fn' TARGET='$(TARGET)' .rust-run

rust-build-all:
	for target in $$(ls ./src-rust/cmd); do \
		$(MAKE) TARGET=$$target ACTION='rust-build-fn' .rust-run; \
	done

rust-test-all:
	cd src-rust && cargo test --workspace
