run:
	go run cmd/main.go cmd/api.go

build:
	go build -o bin/app.exe cmd/main.go cmd/api.go

test:
	go test ./...

swag:
	swag init -g cmd/api.go -o docs

migrate-up:
	goose -dir internal/adapters/mysql/migrations mysql "root:root@tcp(localhost:3306)/go_starter_project?parseTime=true" up

migrate-down:
	goose -dir internal/adapters/mysql/migrations mysql "root:root@tcp(localhost:3306)/go_starter_project?parseTime=true" down
