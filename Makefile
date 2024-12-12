define build_dev
	go build -o ./bin/$(1) ./cmd/$(1)
endef

define build
	go build -ldflags "-s -w" -o ./bin/$(1) ./cmd/$(1)
endef

define dev
	ls ./{lib,cmd}/**/*.go | entr -c $(call build_dev,$(1))
endef

.dev-target:
	@if [[ -n "$(TARGET)" ]]; then \
		(cd go-work && $(call dev,$(TARGET))) \
	fi

.build-target:
	@if [[ -n "$(TARGET)" ]]; then \
		(cd go-work && $(call build,$(TARGET))) \
	fi

dev:
	$(MAKE) TARGET=$$(ls ./go-work/cmd | fzf) .dev-target

build:
	$(MAKE) TARGET=$$(ls ./go-work/cmd | fzf) .build-target

build-all:
	for target in $$(ls ./go-work/cmd); do \
		$(MAKE) TARGET=$$target .build-target; \
	done
