include .env
export

export PROJECT_ROOT=$(shell pwd)

auth-service-run:
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/sso/main.go
env-up:
	@docker compose up -d auth-service-postgres auth-service-port-forwarder
env-down:
	@docker compose down auth-service-postgres auth-service-port-forwarder
	
