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
dev-airbnb-calculator:
	$(MAKE) TARGET=airbnb_calculator .dev-target

dev-get-repo-name:
	$(MAKE) TARGET=get_repo_name .dev-target

dev-repeatstr:
	$(MAKE) TARGET=repeatstr .dev-target

dev-zsb-charm-tmux-renametab:
	$(MAKE) TARGET=zsb_charm_tmux_renametab .dev-target

dev-zsb-charm-tmux-urlopen:
	$(MAKE) TARGET=zsb_charm_tmux_urlopen .dev-target

dev-zsb-clipcopy:
	$(MAKE) TARGET=zsb_clipcopy .dev-target

dev-zsb-open:
	$(MAKE) TARGET=zsb_open .dev-target

build:
	$(MAKE) TARGET=airbnb_calculator .build-target
	$(MAKE) TARGET=get_repo_name .build-target
	$(MAKE) TARGET=repeatstr .build-target
	$(MAKE) TARGET=zsb_charm_tmux_renametab .build-target
	$(MAKE) TARGET=zsb_charm_tmux_urlopen .build-target
	$(MAKE) TARGET=zsb_clipcopy .build-target
	$(MAKE) TARGET=zsb_open .build-target
