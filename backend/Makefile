MIGRATE = migrate
MIGRATIONS_DIR = db/migrations
DB_URL = "mysql://root:fadel123@tcp(localhost:3307)/new_ticket?charset=utf8mb4&parseTime=True&loc=Local"



## Create new migration
create-migration:
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) $(name)

## Run migrations
migrate-up:
	$(MIGRATE) -database $(DB_URL) -path $(MIGRATIONS_DIR) up

## Rollback 1 step
migrate-down:
	$(MIGRATE) -database $(DB_URL) -path $(MIGRATIONS_DIR) down 1

## Drop all migrations
migrate-drop:
	$(MIGRATE) -database $(DB_URL) -path $(MIGRATIONS_DIR) drop -f

### Start docker development
mysql:
	docker container start mysql-container1

app:
	go run cmd/web/main.go