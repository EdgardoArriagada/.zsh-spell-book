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
	(cd go-work && $(call dev,$(TARGET)))

.build-target:
	(cd go-work && $(call build,$(TARGET)))

# Setup
dev:
	$(MAKE) TARGET=$$(ls ./go-work/cmd | fzf) .dev-target

build:
	for target in $$(ls ./go-work/cmd); do \
		$(MAKE) TARGET=$$target .build-target; \
	done
