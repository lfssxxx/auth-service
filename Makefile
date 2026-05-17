include .env
export

export PROJECT_ROOT=$(shell pwd)

auth-service-run:
	@go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/sso/main.go


env-up:
	@docker compose up -d auth-service-postgres

env-down:
	@docker compose down auth-service-postgres
	
env-port-forward:
	@docker compose up -d auth-service-port-forwarder

env-port-close:
	@docker compose down auth-service-port-forwarder