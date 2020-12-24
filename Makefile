build:
	go mod tidy && go mod vendor
	go build -o ./bin/go-live && chmod +X ./bin/*
	echo "Executable Ready in ./bin/go-live"

docker:
	docker build -t antsankov/go-live:latest .

format: 
	gofmt -l -s -w .

# Compiled with Arch Go package guidelines and removed whitespace.
cross-compile:
	env GOOS=linux GOARCH=arm go build -o ./release/go-live-linux-arm32 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=linux GOARCH=arm64 go build -o ./release/go-live-linux-arm64 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=linux GOARCH=386 go build -o ./release/go-live-linux-x32 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=linux GOARCH=amd64 go build -o ./release/go-live-linux-x64 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=darwin GOARCH=amd64 go build -o ./release/go-live-mac-x64 -ldflags "-s -w" -trimpath -mod=readonly 
	env GOOS=darwin GOARCH=arm64 go build -o ./release/go-live-mac-arm64 -ldflags "-s -w" -trimpath -mod=readonly 
	env GOOS=windows GOARCH=386 go build -o ./release/go-live-windows-x32.exe -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=windows GOARCH=amd64 go build -o ./release/go-live-windows-x64.exe -ldflags "-s -w" -trimpath -mod=readonly

clean:
	rm -rf ./bin/*
	rm -rf ./release/*
	rm -rf ./vendor/*

run:
	./bin/go-live
