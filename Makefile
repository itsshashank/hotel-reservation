build:
	@go build -o bin/server cmd/main.go

run: build
	@./bin/server

seed:
	@go run scripts/seed.go

test:
	@go test -v -race -count=1 ./...

db:
	@docker rm -f mongodb
	@docker run --name mongodb -p 27017:27017 -d mongo:latest
	@sleep 1
	@export MONGODB_URI=$(docker inspect -f '{{ .NetworkSettings.IPAddress }}' mongodb)