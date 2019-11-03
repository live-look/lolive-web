.PHONY: all clean sh build

all: build

build:
	scripts/development/build.sh

clean:
	scripts/development/cleanup.sh

sh:
	scripts/developmen/sh.sh
