MAIN_IMAGE := app:latest
MAIN_PACKAGE_PATH := ./cmd/app/
BINARY_NAME := bin/app

CONFIG_PATH ?= ./etc/config.yml
MIGRATE_DSN ?= postgres://postgres@localhost:5432/?sslmode=disable

.PHONY: build
build:
	CGO_ENABLED=0 go build -o ${BINARY_NAME} ${MAIN_PACKAGE_PATH}

.PHONY: run
run: build
	${BINARY_NAME} -config ${CONFIG_PATH}

.PHONY: docker.image
docker.image:
	docker build -t ${MAIN_IMAGE} .

.PHONY: migrate.up
migrate.up:
	migrate -source=file://migrations/ -database "${MIGRATE_DSN}" up

.PHONY: migrate.down
migrate.down:
	migrate -source=file://migrations/ -database "${MIGRATE_DSN}" down

