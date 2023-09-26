![go-live logo](./logo.png)

# go-live
[![Go Report](https://goreportcard.com/badge/github.com/antsankov/go-live)](https://goreportcard.com/report/github.com/antsankov/go-live)
[![GoDoc](https://godoc.org/github.com/antsankov/go-live?status.svg)](https://pkg.go.dev/github.com/antsankov/go-live@v1.2.0?tab=overview)
[![Snap Package](https://snapcraft.io/go-live/badge.svg)](https://snapcraft.io/go-live)


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
**More Info**: https://antsankov.gitbook.io/go-live/

## Install

### MacOS Intel (with Brew)
`brew tap antsankov/go-live && brew install go-live` 

### MacOS Apple Silicon/M1/M2 (with Brew)
* For ARM (Mac M1 / M2) - make sure your Brew is istalled to `opt/homebrew`. Brew does not do this by default, easiest way to do this is to install homebrew via the .pkg from the [`homebrew` github releases page](https://github.com/Homebrew/brew/releases). 
  
`brew tap antsankov/go-live && arch -arm64 brew install go-live`

### MacOS x64 (without Brew)

`curl -LJO https://github.com/antsankov/go-live/releases/download/v1.2.0/go-live-mac-x64.zip && unzip go-live-mac-x64.zip && mv go-live /usr/local/bin/go-live && chmod +x /usr/local/bin/go-live && go-live`

### MacOS Apple Silicon/M1/M2 (without Brew)

`curl -LJO https://github.com/antsankov/go-live/releases/download/v1.2.0/go-live-mac-arm64.zip && unzip go-live-mac-arm64.zip && mv go-live /usr/local/bin/go-live && chmod +x /usr/local/bin/go-live && go-live`

### Linux (using Snapcraft)
`snap install go-live`

### Linux x32 (Ubuntu/RHEL/etc.):
`wget https://github.com/antsankov/go-live/releases/download/v1.2.0/go-live-linux-x32 -O /usr/bin/go-live && chmod +x /usr/bin/go-live`

### Linux x64 (Ubuntu/RHEL/etc.):
`wget https://github.com/antsankov/go-live/releases/download/v1.2.0/go-live-linux-x64 -O /usr/bin/go-live && chmod +x /usr/bin/go-live`

### Linux ARM32 (Ubuntu/RHEL/etc.):
`wget https://github.com/antsankov/go-live/releases/download/v1.2.0/go-live-linux-arm32 -O /usr/bin/go-live && chmod +x /usr/bin/go-live`

### Linux ARM64 (Ubuntu/RHEL/etc.):
`wget https://github.com/antsankov/go-live/releases/download/v1.2.0/go-live-linux-arm64 -O /usr/bin/go-live && chmod +x /usr/bin/go-live`

### Docker
`docker pull antsankov/go-live`

To run (will serve current directory on port 9000):

`docker run --rm -v "${PWD}":/workdir -p 9000:9000 antsankov/go-live go-live`

### Windows

[Download Here and Execute](https://github.com/antsankov/go-live/releases/tag/v1.2.0)

- Chocolatey coming soon! (Help wanted)
- Make sure when running that all necessary ports are open and user has permissions (Help wanted)
- QT based front-end? (Help wanted)

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

## To Release
- For snapcraft it builds automatically when you push it
- For Mac and Homebrew, see https://github.com/mitchellh/gon
  - `gon gon.json`
  - Make sure to have XTools installed, and opened already.
  - You need to have a valid developer certficate - check `security find-identity -p codesigning`. If it is not valid, see https://developer.apple.com/forums/thread/86161 -- you need to check the info of the developer cert to see if the "Organizational Unit" certificate is installed.
  - For gon to work, you need to use the hacked version https://github.com/mitchellh/gon/issues/64#issuecomment-1336311570 to release on Apple Silicon
  - The "ac-password" in gon is an App specfic password for your Apple ID.
- For docker (remember for version and for latest): `sudo docker build -t antsankov/go-live:v1.2.0 .` and `sudo docker push antsankov/go-live:v1.2.0`
## Flags
```
  -h  Print help message for go-live 
  --help

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
- [ ] Android Support
- [x] Docker Support
- [ ] Benchmarking and performance tests. Large files, and concurrent connections.
- [x] Gif and Screenshots of it in use. 
- [x] Tutorial Use as Github Pages Alternative
- [x] Copy Paste from Terminal fix.
- [x] Finish Gitbook documentation. 
- [ ] HTTPS support.
- [x] Publish as a Go package.
- [x] Setup Unit tests.
- [x] Requests Counter
- [x] Ability to download as a binary.
- [x] Browser Opening
- [x] Finish Go Deps
- [x] Run as shell utility.
- [x] Figure out rotating print message.
- [x] Get local server going.
