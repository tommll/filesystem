.PHONY: dep local-db clean run test

FLAGS := -ldflags "-X main.version=$(VERSION)"

setup:
	$(info $(M) go mod vendor)
	go mod vendor

dep: ; $(info Ensuring vendored dependencies are up-to-date...)
	go build $(FLAGS)

local-db:
	@echo "======================== Setup DB (Postgres) ========================"
	@docker-compose -f docker-compose.yml -p fs down
	@docker-compose -f docker-compose.yml -p fs up -d --build
	@echo "Waiting for database connection..."
	@while ! docker exec fs_postgres_1 pg_isready -h localhost -p 5432 > /dev/null; do \
		sleep 1; \
	done
	@echo "Migrate db..."
	@docker cp ./app/database/fs.sql fs_postgres_1:/fs.sql
	@docker exec fs_postgres_1 psql -U postgres -p 5432 -c "CREATE DATABASE fs ENCODING'UTF-8';" || true
	@docker exec fs_postgres_1 psql -U postgres -p 5432 fs -f /fs.sql
