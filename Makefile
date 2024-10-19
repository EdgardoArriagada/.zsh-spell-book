define build
	(cd go-work && go build -o ./bin/$(1) ./$(1))
endef

build-urlopen:
	@$(call build,"urlopen")


