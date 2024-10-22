define build
	go build -o ./bin/$(1) ./cmd/$(1)
endef

define dev
	ls ./{lib,cmd}/**/*.go | entr -c $(call build,$(1))
endef

.build-target:
	(cd go-work && $(call build,$(TARGET)))

.dev-target:
	(cd go-work && $(call dev,$(TARGET)))

build-urlopen:
	$(MAKE) TARGET=zsb_charm_tmux_urlopen .build-target

dev-urlopen:
	$(MAKE) TARGET=zsb_charm_tmux_urlopen .dev-target

