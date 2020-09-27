# go-live


A light-weight Go server that hosts the file directory it's in (HTML or otherwise) over HTTP. Can be used for local development or production static site serving as a drop-in personal Github-pages. By default, it hosts the directory it is executed in. 

Based on JavaScript's famous `live-server` utility. Supports Linux, Windows, and Mac, as well as ARM32/64.

Help wanted! Check out the TODO list if interested.

## Uses
* Local development of an HTML project. 
* Host a production static site on a server as a Github-Pages alternative. We recommend using nginx for HTTPS, and Cloudflare caching.
* A very lightweight network file-hosting server on a LAN or the Internet. Run it in the folder you want to share!
* Use in an embedded system to either share files or host a static website. 

## Flags
`-p / --port` : The port you want to run on. Default: `9000`

`-d / --dir` : The directory you want to host from. Default: `./` (current directory)

`-q / --quiet` : Set quiet mode to on to avoid opening browser on startup. Default `false`

`-c / --use-browser-cache` : Allow browser to cache pages. Bad for development, but good if you're hosting prod and can't use a load balancer or if HTML pages never change. Default `false`

`--help` : Help menu

Note: `index.html` is displayed automatically at the root of a directory.

**Example**: Serve a site on Port 80

`./go-live --dir ~/example.com/ --port 80`

##  Install
Requires Go to be installed.
```
git clone https://github.com/antsankov/go-live.git && cd go-live
make build
```

If you want to run directly do:
```
./bin/go-live --port 8888 --dir ../../example-site/
```

or you can move it to your path and then use it like a shell utility:

```
chmod +X ./bin/*
mv ./bin/go-live /usr/local/bin
go-live ...
```

## TODO
- [ ] Copy Paste from Terminal fix.
- [ ] Publish Gitbook documentation. 
- [ ] HTTPS support.
- [ ] Publish as a Go package.
- [ ] Local development mode refresh page. 
- [ ] Setup Unit tests.
- [x] Ability to download as a binary.
- [x] Browser Opening
- [x] Finish Go Deps
- [x] Run as shell utility.
- [x] Figure out rotating print message.
- [x] Get local server going.
