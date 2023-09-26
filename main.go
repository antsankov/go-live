package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/antsankov/go-live/lib"
	"github.com/pkg/browser"
)

// VERSION of Package
const VERSION = "1.2.0"


func main() {
	var _quiet bool
	flag.BoolVar(&_quiet, "q", false, "Quiet stops go-live from opening the browser when started.")
	flag.BoolVar(&_quiet, "quiet", false, "")
	var _cache bool
	flag.BoolVar(&_cache, "c", false, "Allow browser caching of pages. Can lead to stale results, off by default.")
	flag.BoolVar(&_cache, "cache", false, "")
	var _port string
	flag.StringVar(&_port, "p", "9000", "Set port to serve on.")
	flag.StringVar(&_port, "port", "9000", "")
	var _version bool
	flag.BoolVar(&_version, "v", false, "Print the version of go-live.")
	flag.BoolVar(&_version, "version", false, "")
	var _dir string
	flag.StringVar(&_dir, "d", "./", "Select the directory you want to serve. Serves all subpaths that user has read permissions for.")
	flag.StringVar(&_dir, "dir", "./", "")
	var _serve bool
	flag.BoolVar(&_serve, "s", false, "Start in server mode on port 80 and in quiet.")
	flag.BoolVar(&_serve, "serve", false, "")

	flag.Parse()

	if _version || (len(os.Args) >= 2 && os.Args[1] == "version") {
		fmt.Printf("v%s (%s/%s)\n", VERSION, runtime.GOOS, runtime.GOARCH)
		return
	}

	if _dir != "./" {
		// Check if last char is a slash, if not add it.
		if _dir[len(_dir)-1] != '/' {
			_dir = _dir + "/"
		}
	}

	// Check if port begins with ":", if not add it.
	if _port[0] != ':' {
		_port = ":" + _port
	}

	var err error
	if _serve {
		err = lib.StartServer(_dir, ":80", true)
	} else {
		// If user is sudo we don't launch the browser.
		if !_quiet && !isSudo() {
			browser.OpenURL(fmt.Sprintf("http://localhost%s", _port))
		}
		err = lib.StartServer(_dir, _port, _cache)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// IsSudo checks if user is sudo
func isSudo() bool {
	return (os.Geteuid() == 0)
}
