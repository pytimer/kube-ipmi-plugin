
IMG ?= pytimer/kube-ipmi-plugin:v0.1

.PHONY: build
build:
	go build -o bin/kube-ipmi-plugin main.go

docker-build:
	docker build -t ${IMG} .