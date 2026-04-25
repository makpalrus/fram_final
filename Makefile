DB_URL=postgres://postgres:Kz123456@localhost:5432/fram?sslmode=disable
MIGRATE_PATH=migrations

migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1

migrate-create:
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

migrate-version:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" version

migrate-force:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" force $(v)