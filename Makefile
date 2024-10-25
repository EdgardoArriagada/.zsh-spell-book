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

dev-zsb-charm-tmux-urlopen:
	$(MAKE) TARGET=zsb_charm_tmux_urlopen .dev-target

dev-zsb-clipcopy:
	$(MAKE) TARGET=zsb_clipcopy .dev-target

dev-zsb-open:
	$(MAKE) TARGET=zsb_open .dev-target

build:
	$(MAKE) TARGET=zsb_charm_tmux_urlopen .build-target
	$(MAKE) TARGET=zsb_clipcopy .build-target
	$(MAKE) TARGET=zsb_open .build-target
