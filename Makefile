define build
	go build -o ./bin/$(1) ./cmd/$(1)
endef

build-urlopen:
	(cd go-work && $(call build,"zsb_charm_tmux_urlopen"))


dev-urlopen:
	(cd go-work && ls ./cmd/**/*.go | entr -c $(call build,"zsb_charm_tmux_urlopen"))

