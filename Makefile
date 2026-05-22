include .env
export

export PROJECT_ROOT=$(shell pwd)

auth-service-run:
	@go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/sso/main.go

auth-service-tests:
	@go test ${PROJECT_ROOT}/tests

env-up:
	@docker compose up -d auth-service-postgres

env-down:
	@docker compose down auth-service-postgres
	
env-port-forward:
	@docker compose up -d auth-service-port-forwarder

env-port-close:
	@docker compose down auth-service-port-forwarder

	


migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутсвует необходимый параметр seq. Пример make migrate-create seq=value"; \
		exit 1; \
	fi; \
	    docker compose run --rm  auth-service-migrate \
		create -ext sql -dir /migrations -seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутсвует необходимый параметр action. Пример make migrate-action action=up"; \
		exit 1; \
	fi; \
	 	docker compose run --rm auth-service-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@auth-service-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)" \

migrate-force:
	docker compose run --rm auth-service-migrate \
	-path /migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@auth-service-postgres:5432/${POSTGRES_DB}?sslmode=disable \
	force 1