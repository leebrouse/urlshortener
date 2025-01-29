IMAGE_NAME=postgres
CONTAINER_NAME=postgres
POSTGRES_USER=root
POSTGRES_PASSWORD=root

# Default target to start PostgreSQL container
run_postgres:
	docker run -d \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-p 5432:5432 \
		--name $(CONTAINER_NAME) $(IMAGE_NAME)

# Stop and remove the PostgreSQL container
stop_postgres:
	docker stop $(CONTAINER_NAME)

# Restart the PostgreSQL container
restart_postgres: docker restart $(CONTAINER_NAME)

# Check the status of the PostgreSQL container
status_postgres:
	docker ps -f name=$(CONTAINER_NAME)

# Remove the Docker image (if needed)
remove-image:
	docker rmi $(IMAGE_NAME)

# migrate
install_migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# sqlc
install_sqlc:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# redis
run_redis:
	docker run -d \
	 -p 6379:6379 \
	 --name redis_urls \
	 redis


databaseURL="postgres://root:root@localhost:5432/urldb?sslmode=disable"
#migrate db up
migrate_up:
	migrate -path="./database/migrate" -database $(databaseURL) -verbose up

#migrate db drop
migrate_drop:
	migrate -path="./database/migrate" -database $(databaseURL) -verbose drop -f