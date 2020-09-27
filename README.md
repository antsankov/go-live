# go-live
![Github Releases (by Release)](https://img.shields.io/github/downloads/antsankov/go-live/total.svg) ![Go Report](https://goreportcard.com/badge/github.com/antsankov/go-live) [![GoDoc](https://godoc.org/github.com/antsankov/go-live?status.svg)](https://godoc.org/github.com/antsankov/go-live)

**Docs**: https://antsankov.gitbook.io/go-live/

A simple, portable Go server that hosts a file directory over HTTP. Can be used for local web development or production static site serving. Can also be used as a network file sharing server. By default, `go-live` serves the directory it is executed in.

Based on JavaScript's famous `live-server` utility. Supports Linux, Windows, and Mac, as well as ARM.

Help wanted! Check out the TODO list if interested.

## When to Use
* Local development and serving of an HTML project (can run any compiled web-project code), prevents caching.
* Host a production static site yourself as a GitHub Pages alternative.
* A lightweight network file-hosting server. Can be used on a LAN or the Internet. Simply run `go-live` in the folder you want to share.
* Use in an embedded IoT system to either share files on a network or host a static website.

## Install

### [Download latest Binary](https://github.com/antsankov/go-live/releases)

### MacOS

`wget https://github.com/antsankov/go-live/releases/download/latest/go-live-mac-x64 && chmod 755 go-live-mac-x64 && mv go-live-mac-x64 /usr/local/bin`

- Homebrew coming soon!
- Verify checksum of binary in checksum.txt

### Ubuntu

`wget https://github.com/antsankov/go-live/releases/download/v0.0.2/go-live-linux-x32 && chmod 755 go-live-linux-x32 && mv go-live-linux-x32 /usr/local/bin/go-live`

- Snap coming soon! (Help wanted)
- Verify checksum of binary in checksum.txt

### Windows

[Download Here and Execute](https://github.com/antsankov/go-live/releases/download/v0.0.2/go-live-windows-x64.exe)

- Chocolatey coming soon! (Help wanted)
- Verify checksum of binary in checksum.txt.

### Go Get (must have Go installed)
`GO111MODULE=on go get github.com/antsankov/go-live`

### Install From Source (must have Go installed)
`make build && ./bin/go-live`

### Cross Compile for multiple systems
`make cross-compile && ls release/`

(Need 32-bit or ARM? Check the releases page.)

## Flags
`-p / --port` : The port you want to run on. Default: `9000`

`-d / --dir` : The directory you want to host from. Default: `./` (current directory)

`-q / --quiet` : Set quiet mode to on to avoid opening browser on startup. Default `false`

`-c / --use-browser-cache` : Allow browser to cache pages. Bad for development, but good if you're hosting prod and can't use a load balancer or if HTML pages never change. Default `false`

`-v / --version / version` : Print the Version of go-live.

`--help` : Help menu

Note: `index.html` is displayed automatically at the root of a directory.

**Example**: Serve a static site over Port 80

`go-live --dir ~/example.com/ --port 80 --quiet`

## TODO
- [ ] Tutorial Use as Github Pages Alternative
- [ ] Copy Paste from Terminal fix.
- [x] Finish Gitbook documentation. 
- [ ] HTTPS support.
- [x] Publish as a Go package.
- [ ] Local development refresh on file change. 
- [ ] Setup Unit tests.
- [x] Requests Counter
- [x] Ability to download as a binary.
- [x] Browser Opening
- [x] Finish Go Deps
- [x] Run as shell utility.
- [x] Figure out rotating print message.
- [x] Get local server going.

## Release Steps
1. checksum.txt
1. version.go
1. Download links on README.