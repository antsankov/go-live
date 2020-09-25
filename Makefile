build:
	go build -o ./bin/go-live && chmod +X ./bin/*

clean:
	rm -rf ./bin/*

run:
	go run main.go
