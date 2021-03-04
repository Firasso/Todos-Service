.PHONY: build up down deps-dev test infra secrets

deps-dev:
	@if ! command -v CompileDaemon > /dev/null ; then \
		echo ">> Installing CompileDaemon"; \
		go get github.com/githubnemo/CompileDaemon; \
	fi

## Builds, (re)creates, starts, and de-attaches containers for services in docker-compose.yml
up:
	@docker-compose up --build

## Stops and removes docker containers ( for services not defined in docker-compose.yml )
down:
	@docker-compose down --remove-orphans

test:
	@docker-compose  -f docker-compose.test.yml up --build --abort-on-container-exit

build:
	@echo ">> building binaries"
	@go build -o ./build/Demo-app main.go

check-swagger:
	which swagger || GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check-swagger
	swagger generate spec -o ./docs/swagger.yaml --scan-models

infra:
	@docker-compose -f docker-compose.infra.yml run --rm aws

secrets:
	@kubectl create secret generic env-secrets --from-env-file=.env
