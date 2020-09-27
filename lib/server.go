package lib

import (
	"log"
	"net/http"
	"time"
)

var requests uint64 = 0
var epoch = time.Unix(0, 0).Format(time.RFC1123)
var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}
var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func incrementRequest() {
	requests++
}

// StartServer starts up the file server
func StartServer(dir string, port string, cache bool) {
	fs := http.FileServer(http.Dir(dir))
	if cache {
		http.Handle("/", useCache(fs))
	} else {
		http.Handle("/", noCache(fs))
	}
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func noCache(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}
		incrementRequest()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func useCache(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		incrementRequest()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
