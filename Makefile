build:
	go mod tidy && go mod vendor
	go build -o ./bin/go-live && chmod +X ./bin/*

lint: 
	gofmt -l -s -w .

cross-compile:
	env GOOS=linux GOARCH=arm go build -o ./release/go-live-linux-arm32 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=linux GOARCH=arm64 go build -o ./release/go-live-linux-arm64 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=linux GOARCH=386 go build -o ./release/go-live-linux-x32 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=linux GOARCH=amd64 go build -o ./release/go-live-linux-x64 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=darwin GOARCH=amd64 go build -o ./release/go-live-mac-x64 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=windows GOARCH=amd64 go build -o ./release/go-live-windows-x64.exe -ldflags "-s -w" -trimpath -mod=readonly
	
clean:
	rm -rf ./bin/*
	rm -rf ./release/*

run:
	./bin/go-live
