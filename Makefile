include .env
export 

# Local start
start:
	@echo  "=== Start Miniotaur ==="
	docker-compose up -d minio
	S3_BUCKET_ENDPOINT=http://localhost:9000 \
	go run main.go

# Docker Start
startd:
	@echo  "=== Start Miniotaur ==="
	docker-compose up -d
	@make logs

down:
	docker compose down

logs:
	docker compose logs -f