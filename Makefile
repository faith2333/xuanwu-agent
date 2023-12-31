VERSION=$(shell git describe --tags --always)
IMAGE_REPO="swr.cn-east-3.myhuaweicloud.com/turingsyn"

.PHONY: build
# build binary file for tscli, output path is ./bin/
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: build-linux
# build binary file for tscli, output path is ./bin/
build-linux:
	mkdir -p bin/ && CGO_ENABLED=0  GOOS=linux  GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: build-image
build-image: build-linux
	docker build -t $(IMAGE_REPO)/xuanwu-agent:$(VERSION) .
	docker push $(IMAGE_REPO)/xuanwu-agent:$(VERSION)

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help