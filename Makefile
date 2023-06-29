include .env
export 

start:
	docker-compose up -d
	go run main.go