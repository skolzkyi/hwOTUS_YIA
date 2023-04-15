BIN := "./bin/calendar"
SCHEDULER_BIN := "./bin/calendar_scheduler"
SENDER_BIN := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"
DSN="imapp:LightInDark@/OTUSFinalLab?parseTime=true"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

mysql-img:
	docker-compose -f ./deployments/docker-compose.yaml up -d
	
mysql-img-down:
	docker-compose -f ./deployments/docker-compose.yaml down

migrate-goose:
	goose --dir=migrations mysql $(DSN) up

generate:
	rm -rf internal/server/grpc/pb
	mkdir -p internal/server/grpc/pb
	protoc api/EventService.proto  --go_out=internal/server/grpc/pb --go-grpc_out=internal/server/grpc/pb 

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(SCHEDULER_BIN) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(SENDER_BIN) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(BIN) -config ./configs/config.env
	$(SCHEDULER_BIN) -config ./configs/config_sheduler.env
	$(SENDER_BIN) -config ./configs/config_sender.env

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... 

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

up:
	docker-compose -f build/docker-compose.yaml up

down:
	docker-compose -f build/docker-compose.yaml down

.PHONY: generate build run build-img mysql-img migrate-goose run-img version test lint
