export ENV := test
export GAME_DB_HOST := localhost
export GAME_DB_PORT := 5432
export GAME_DB_USER := admin
export GAME_DB_PASSWORD := 12345
export GAME_DB_NAME := game-api
export GAME_DB_CONNECTION := postgresql://$(GAME_DB_USER):$(GAME_DB_PASSWORD)@$(GAME_DB_HOST):$(GAME_DB_PORT)/$(GAME_DB_NAME)?sslmode=disable

dev:
	go run ./cmd/api/

migrate-create:
	migrate create -ext=.sql -dir=./migrations $(name)

migrate-push:
	migrate -path=./migrations -database=$(GAME_DB_CONNECTION) up

migrate-down:
	migrate -path=./migrations -database=$(GAME_DB_CONNECTION) down 1
