
# Allow `make go-build <name>` or `make go-dev <name>` to skip fzf selection
ifneq ($(filter $(firstword $(MAKECMDGOALS)),go-build go-dev),)
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
		TARGET=$$(ls ./go-work/cmd | fzf); \
	fi; \
	if [[ -n "$$TARGET" ]]; then \
		(cd go-work && $(call $(ACTION),$$TARGET)); \
	fi

go-build:
	$(MAKE) ACTION='go-build-fn' TARGET='$(TARGET)' .go-run

go-dev:
	$(MAKE) ACTION='go-dev-fn' TARGET='$(TARGET)' .go-run

zsh-dev:
	@find . -name "*.zsh" -not -path "./go-work/*" | entr -c ./go-work/bin/zsb_bundle

go-build-all:
	for target in $$(ls ./go-work/cmd); do \
		$(MAKE) TARGET=$$target ACTION='go-build-fn' .go-run; \
	done

go-test-all:
	cd go-work && go list -f '{{.Dir}}' -m | xargs -I{} sh -c 'cd "{}" && go test ./...'
