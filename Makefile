.PHONY: all clean sh build

all: build

up:
	scripts/development/up.sh

down:
	scripts/development/down.sh

migrate:
	scripts/development/migrate.sh

add_migration:
	scripts/development/add_migration.sh

build:
	scripts/development/build.sh

clean:
	scripts/development/cleanup.sh

sh:
	scripts/development/sh.sh

psql:
	scripts/development/psql.sh
