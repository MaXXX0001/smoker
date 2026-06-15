.PHONY: proto test build up down tidy

# Перегенерувати gRPC-стаби з proto/.
proto:
	buf generate

# Усі тести (домен + парсери, без мережі).
test:
	go test ./...

# Зібрати всі бінарі локально.
build:
	go build ./...

tidy:
	go mod tidy

# Підняти/зупинити весь стек.
up:
	docker compose up --build

down:
	docker compose down
