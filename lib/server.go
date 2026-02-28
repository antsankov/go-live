package lib

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

var requests uint64
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
	atomic.AddUint64(&requests, 1)
}

// StartServer starts up the file server
func StartServer(dir string, port string, cache bool, tlsEnabled bool, certFile string, keyFile string) error {
	go Printer(dir, port, tlsEnabled)
	fs := http.FileServer(http.Dir(dir))
	if cache {
		http.Handle("/", useCache(fs))
	} else {
		http.Handle("/", noCache(fs))
	}

	if tlsEnabled {
		if certFile != "" && keyFile != "" {
			server := &http.Server{
				Addr:         port,
				TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
			}
			return server.ListenAndServeTLS(certFile, keyFile)
		}
		// Self-signed certificate
		cert, err := GenerateSelfSignedCert()
		if err != nil {
			return err
		}
		ln, err := net.Listen("tcp", port)
		if err != nil {
			return err
		}
		server := &http.Server{
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
			ErrorLog:     tlsHandshakeFilter(),
		}
		return server.ServeTLS(ln, "", "")
	}

	return http.ListenAndServe(port, nil)
}

// tlsHandshakeFilter returns a logger that suppresses expected TLS handshake
// errors from clients that don't trust the self-signed certificate.
func tlsHandshakeFilter() *log.Logger {
	return log.New(tlsFilterWriter{}, "", 0)
}

type tlsFilterWriter struct{}

func (w tlsFilterWriter) Write(p []byte) (n int, err error) {
	msg := string(p)
	if strings.Contains(msg, "TLS handshake error") {
		return len(p), nil
	}
	return os.Stderr.Write(p)
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
