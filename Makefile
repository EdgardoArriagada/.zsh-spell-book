define build
	go build -o ./bin/$(1) ./$(1)
endef

build-urlopen:
	(cd go-work && $(call build,"zsb_charm_tmux_urlopen"))


dev-urlopen:
	(cd go-work && ls **/*.go | entr -c $(call build,"zsb_charm_tmux_urlopen"))

