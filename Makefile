MIGRATE=migrate
DB_URL=postgres://prabhat:prabhat@localhost:5432/GO?sslmode=disable
MIGRATIONS_DIR=src/core/database/migration

.PHONY: migrate-up migrate-down migrate-new

migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

migrate-new:
	@if [ -z "$(name)" ]; then \
		echo "❌ Please provide a name. Usage: make migrate-new name=create_users_table"; \
		exit 1; \
	fi
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)
