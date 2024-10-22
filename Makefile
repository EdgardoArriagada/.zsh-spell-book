define build
	go build -o ./bin/$(1) ./cmd/$(1)
endef

define dev
	ls ./{lib,cmd}/**/*.go | entr -c $(call build,$(1))
endef

.dev-target:
	(cd go-work && $(call dev,$(TARGET)))

dev-urlopen:
	$(MAKE) TARGET=zsb_charm_tmux_urlopen .dev-target

dev-zsh-open:
	$(MAKE) TARGET=zsb_open .dev-target

dev-zsb-clipcopy:
	$(MAKE) TARGET=zsb_clipcopy .dev-target
