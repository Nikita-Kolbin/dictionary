lint:
	golangci-lint run ./...

compose_up:
	docker compose --env-file ./env/.env up -d

compose_down:
	docker compose down

compose_drop:
	docker compose down -v

compose_rebuild:
	docker compose down -v
	docker compose --env-file ./env/.env up -d --build --force-recreate

# Накатить миграции
compose_migrate:
	docker compose --env-file ./env/.env -p dictionary up -d migrate --build --force-recreate

# Перебилдить бэк
compose_backend:
	docker compose --env-file ./env/.env -p dictionary up -d backend --build --force-recreate

# migrate create -ext sql -dir ./migrations -seq <name>
