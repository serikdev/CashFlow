run:
	go run cmd/api/main.go

migrate:
	bash scripts/migrate.sh up

migrate-down:
	bash scripts/migrate.sh down

migrate-status:
	bash scripts/migrate.sh status

migrate-redo:
	bash scripts/migrate.sh redo

migrate-create:
	@read -p "Migration name: " name; bash scripts/migrate.sh create $$name



