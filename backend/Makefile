.PHONY: migrate migrate-down swagger

migrate:
	go run cmd/migrations/main.go -up

migrate-down:
	go run cmd/migrations/main.go -down

swagger:
	$(HOME)/go/bin/swag init -g cmd/main.go

run:
	go run cmd/main.go

dev: swagger run