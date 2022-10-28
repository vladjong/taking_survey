APP_BIN = build/app

build: clean $(APP_BIN)

$(APP_BIN):
	go build -o $(APP_BIN) cmd/main/main.go
	./build/app

clean:
	rm -rf build || true

lint:
	golangci-lint run