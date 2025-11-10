.PHONY: build run seed swagger test clean

build:
	go build -o bin/todo-server ./cmd/server

run: build
	./bin/todo-server -p 8080

seed:
	go run scripts/seed.go

swagger:
	swag init -g cmd/server/main.go -o docs

test:
	go test ./...

clean:
	rm -rf bin/