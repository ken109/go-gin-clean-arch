include .env

connection_string := $(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?multiStatements=True

# 
# migration
# 
.PHONY: goose-create-sql
goose-create-sql:
	@docker compose exec app goose -dir driver/migrations create $(FILE) sql

.PHONY: goose-status
goose-status:
	@docker compose exec app goose -dir driver/migrations mysql "$(connection_string)" status

.PHONY: goose-up
goose-up:
	@docker compose exec app goose -allow-missing -dir driver/migrations/ mysql "$(connection_string)" up
