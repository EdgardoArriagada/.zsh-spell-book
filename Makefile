define build
	go build -o ./bin/$(1) ./$(1)
endef

build-urlopen:
	(cd go-work && $(call build,"urlopen"))


dev-urlopen:
	(cd go-work && ls **/*.go | entr -c $(call build,"urlopen"))

