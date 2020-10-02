<img src="https://user-images.githubusercontent.com/2533512/94706954-16a92800-0300-11eb-97a1-3524d22d7c6d.png" width="75" height="75">

# go-live
[![Go Report](https://goreportcard.com/badge/github.com/antsankov/go-live)](https://goreportcard.com/report/github.com/antsankov/go-live) [![GoDoc](https://godoc.org/github.com/antsankov/go-live?status.svg)](https://pkg.go.dev/github.com/antsankov/go-live/)

**Docs**: https://antsankov.gitbook.io/go-live/

A fast, portable Go command line utility that serves a file directory over HTTP. Can be used for local web development, production static-site serving, or as a network file host. By default, `go-live` serves the directory it is executed in.

Based on JavaScript's famous `live-server` utility. Supports Linux, Windows, and Mac, as well as ARM. See TODO list if interested in helping.

*To use*: Run `go-live` in your terminal while in directory you want to serve.

## Example

![go-live-demo](https://user-images.githubusercontent.com/2533512/94636832-5554c900-0293-11eb-8aea-585f8d007fab.gif)

## Use-Cases
* Local development of an HTML/JS project (can serve any frontend code).
* Host a production static site yourself as a GitHub Pages alternative.
* A lightweight network file host that can serve over a LAN or the whole Internet.
* Run on an embedded system or Kubernetes cluster to share files and host a static website on a network (full binary is less than 5MB). 

## Install

### [Download latest Binary](https://github.com/antsankov/go-live/releases)

### MacOS

`brew tap antsankov/go-live`

`brew install go-live` 

(without Brew installed)

`curl -LJO https://github.com/antsankov/go-live/releases/download/v1.0.0/go-live-mac.zip && unzip -d /usr/local/bin/go-live && chmod +x /usr/local/bin/go-live && go-live`

### Linux x32:
`wget https://github.com/antsankov/go-live/releases/download/v1.0.0/go-live-linux-x32 -O /usr/bin/go-live && chmod +x /usr/bin/go-live`

### Linux x64:
`wget https://github.com/antsankov/go-live/releases/download/v1.0.0/go-live-linux-x64 -O /usr/bin/go-live && chmod +x /usr/bin/go-live`

- Deb packages and snap coming soon (Help wanted)
- Verify checksum of binary in checksum.txt
- Need ARM? Check the releases page.

### Windows

[Download Here and Execute](https://github.com/antsankov/go-live/releases/tag/v1.0.0)

- Chocolatey coming soon! (Help wanted)
- Make sure when running that all necessary ports are open (Help wanted)
- Verify checksum of binary in checksum.txt.

### Go Get (must have Go installed)
`GO111MODULE=on go get github.com/antsankov/go-live`

### Install From Source (must have Go installed)
```
git clone https://github.com/antsankov/go-live.git && cd go-live
make build && ./bin/go-live
```
### Cross Compile for multiple systems
```
git clone https://github.com/antsankov/go-live.git && cd go-live
make cross-compile && ls release/
```


## Flags
```
  -c	Allow browser caching of pages. Can lead to stale results, off by default.
  --cache

  -d string
    	Select the directory you want to serve. Serves all subpaths that user has read permissions for. (default "./")
  --dir string
    	 (default "./")
  -p string
    	Set port to serve on. (default "9000")
  --port string
    	 (default "9000")
  -q	Quiet stops go-live from opening the browser when started.
  --quiet

  -s	Start in server mode on port 80 and in quiet.
  --serve

  -v	Print the version of go-live.
  --version
```

Note: `index.html` is displayed automatically at the root of a directory.

**Example**: Serve a static site over Port 80

`sudo go-live --dir ~/example.com/ --serve`

## TODO (Help Wanted)
- [ ] Docker Support
- [ ] Benchmarking and performance tests. Large files, and concurrent connections.
- [x] Gif and Screenshots of it in use. 
- [ ] Tutorial Use as Github Pages Alternative
- [X] Copy Paste from Terminal fix.
- [x] Finish Gitbook documentation. 
- [ ] HTTPS support.
- [x] Publish as a Go package.
- [ ] Setup Unit tests.
- [x] Requests Counter
- [x] Ability to download as a binary.
- [x] Browser Opening
- [x] Finish Go Deps
- [x] Run as shell utility.
- [x] Figure out rotating print message.
- [x] Get local server going.