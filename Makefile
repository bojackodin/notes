MAIN_IMAGE := app:latest
MAIN_PACKAGE_PATH := ./cmd/app/
BINARY_NAME := bin/app

CONFIG_PATH ?= ./etc/config.yml

.PHONY: build
build:
	CGO_ENABLED=0 go build -o ${BINARY_NAME} ${MAIN_PACKAGE_PATH}

.PHONY: run
run: build
	${BINARY_NAME} -config ${CONFIG_PATH}

.PHONY: docker.image
docker.image:
	docker build -t ${MAIN_IMAGE} .

