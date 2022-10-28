APP_BIN = build/app

build: clean $(APP_BIN)

$(APP_BIN):
	go build -o $(APP_BIN) main.go
	./build/app

clean:
	rm -rf build || true
