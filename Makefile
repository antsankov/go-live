VERSION ?= dev
LDFLAGS := -s -w -X main.VERSION=$(VERSION)

build:
	go mod tidy
	go build -o ./bin/go-live -ldflags "$(LDFLAGS)" -trimpath && chmod +x ./bin/*
	echo "Executable Ready in ./bin/go-live"

docker:
	docker build -t antsankov/go-live:latest .

format:
	gofmt -l -s -w .

cross-compile:
	rm -rf ./release && mkdir -p ./release
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./release/go-live-mac-x64 -ldflags "$(LDFLAGS)" -trimpath
	env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./release/go-live-mac-arm64 -ldflags "$(LDFLAGS)" -trimpath
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./release/go-live-linux-x64 -ldflags "$(LDFLAGS)" -trimpath
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./release/go-live-linux-arm64 -ldflags "$(LDFLAGS)" -trimpath
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./release/go-live-windows-x64.exe -ldflags "$(LDFLAGS)" -trimpath
	env CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o ./release/go-live-windows-arm64.exe -ldflags "$(LDFLAGS)" -trimpath

test:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

bench:
	go test ./lib/ -tags bench -bench=. -benchmem -run='^$$' -timeout=300s

bench-compare:
	go run benchmark/compare.go

clean:
	rm -rf ./bin/*
	rm -rf ./release/*

run:
	./bin/go-live
