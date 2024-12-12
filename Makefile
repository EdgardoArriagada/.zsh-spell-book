define build_dev
	go build -o ./bin/$(1) ./cmd/$(1)
endef

define build
	go build -ldflags "-s -w" -o ./bin/$(1) ./cmd/$(1)
endef

define dev
	ls ./{lib,cmd}/**/*.go | entr -c $(call build_dev,$(1))
endef

.run:
	@if [[ -n "$(TARGET)" ]]; then \
		(cd go-work && $(call $(ACTION),$(TARGET))) \
	fi

.choose-target:
	$(MAKE) TARGET=$$(ls ./go-work/cmd | fzf) .run

build:
	$(MAKE) ACTION='build' .choose-target

dev:
	$(MAKE) ACTION='dev' .choose-target

build-all:
	for target in $$(ls ./go-work/cmd); do \
		$(MAKE) TARGET=$$target ACTION='build' .run; \
	done
