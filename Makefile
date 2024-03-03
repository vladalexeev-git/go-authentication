include .env
export

compose-up:
	docker-compose up --build -d postgres && docker-compose logs -f
.PHONY: compose-up

compose-down:
	docker-compose down --remove-orphans
.PHONY: compose-down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(NAME)
.PHONY: migrate-create

run:
	go mod tidy && go mod download && \
	GIN_MODE=debug CGO_ENABLED=0 go run -tags migrate ./cmd/app
.PHONY: run