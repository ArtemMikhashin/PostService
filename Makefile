include .env

db.up:
	docker compose up -d postgres

db.down:
	docker compose down

db.exec:
	docker exec -it postservice-postgres-1 psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

migrate.up:
	docker run --rm \
		--network postservice_posts_net \
		-v $(PWD)/migrations:/migrations \
		migrate/migrate:v4.15.2 \
		-path /migrations \
		-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:5432/$(POSTGRES_DB)?sslmode=disable" \
		up

.PHONY: db.up db.down db.exec migrate.up