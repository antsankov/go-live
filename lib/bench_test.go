//go:build bench

package lib

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

var silentLogger = log.New(io.Discard, "", 0)

// setupTestDir creates a temp directory with files of various sizes for benchmarking.
func setupTestDir(b *testing.B) string {
	b.Helper()
	dir := b.TempDir()

	// Small file (HTML page ~1KB)
	small := []byte("<html><body><h1>Hello</h1><p>" + string(make([]byte, 1000)) + "</p></body></html>")
	os.WriteFile(filepath.Join(dir, "small.html"), small, 0644)

	// Medium file (~100KB, typical JS bundle)
	medium := make([]byte, 100*1024)
	for i := range medium {
		medium[i] = 'a' + byte(i%26)
	}
	os.WriteFile(filepath.Join(dir, "medium.js"), medium, 0644)

	// Large file (~1MB, image-sized)
	large := make([]byte, 1024*1024)
	for i := range large {
		large[i] = byte(i % 256)
	}
	os.WriteFile(filepath.Join(dir, "large.bin"), large, 0644)

	return dir
}

// newKeepAliveClient returns an HTTP client configured for connection reuse.
func newKeepAliveClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
			DisableKeepAlives:   false,
		},
	}
}

func benchmarkHandler(b *testing.B, path string, handler http.Handler) {
	ts := httptest.NewServer(handler)
	defer ts.Close()

	client := newKeepAliveClient()

	b.ResetTimer()
	for b.Loop() {
		resp, err := client.Get(ts.URL + path)
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

// --- No-cache middleware benchmarks ---

func BenchmarkNoCache_SmallFile(b *testing.B) {
	dir := setupTestDir(b)
	fs := http.FileServer(http.Dir(dir))
	benchmarkHandler(b, "/small.html", noCache(fs))
}

func BenchmarkNoCache_MediumFile(b *testing.B) {
	dir := setupTestDir(b)
	fs := http.FileServer(http.Dir(dir))
	benchmarkHandler(b, "/medium.js", noCache(fs))
}

func BenchmarkNoCache_LargeFile(b *testing.B) {
	dir := setupTestDir(b)
	fs := http.FileServer(http.Dir(dir))
	benchmarkHandler(b, "/large.bin", noCache(fs))
}

// --- Cache middleware benchmarks ---

func BenchmarkCache_SmallFile(b *testing.B) {
	dir := setupTestDir(b)
	fs := http.FileServer(http.Dir(dir))
	benchmarkHandler(b, "/small.html", useCache(fs))
}

func BenchmarkCache_MediumFile(b *testing.B) {
	dir := setupTestDir(b)
	fs := http.FileServer(http.Dir(dir))
	benchmarkHandler(b, "/medium.js", useCache(fs))
}

func BenchmarkCache_LargeFile(b *testing.B) {
	dir := setupTestDir(b)
	fs := http.FileServer(http.Dir(dir))
	benchmarkHandler(b, "/large.bin", useCache(fs))
}

// --- Raw file server (no middleware) for baseline ---

func BenchmarkRaw_SmallFile(b *testing.B) {
	dir := setupTestDir(b)
	fs := http.FileServer(http.Dir(dir))
	benchmarkHandler(b, "/small.html", fs)
}

func BenchmarkRaw_MediumFile(b *testing.B) {
	dir := setupTestDir(b)
	fs := http.FileServer(http.Dir(dir))
	benchmarkHandler(b, "/medium.js", fs)
}

func BenchmarkRaw_LargeFile(b *testing.B) {
	dir := setupTestDir(b)
	fs := http.FileServer(http.Dir(dir))
	benchmarkHandler(b, "/large.bin", fs)
}

// --- HTTP vs HTTPS/1.1 vs HTTP/2 comparison ---

func BenchmarkProtocol(b *testing.B) {
	dir := setupTestDir(b)

	cert, err := GenerateSelfSignedCert()
	if err != nil {
		b.Fatal(err)
	}

	files := []struct {
		name string
		path string
	}{
		{"1KB", "/small.html"},
		{"100KB", "/medium.js"},
		{"1MB", "/large.bin"},
	}

	for _, f := range files {
		b.Run(f.name, func(b *testing.B) {
			// HTTP (plaintext)
			b.Run("HTTP", func(b *testing.B) {
				fs := http.FileServer(http.Dir(dir))
				ts := httptest.NewServer(noCache(fs))
				defer ts.Close()
				client := newKeepAliveClient()
				b.ResetTimer()
				for b.Loop() {
					resp, err := client.Get(ts.URL + f.path)
					if err != nil {
						b.Fatal(err)
					}
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
				}
			})

			// HTTPS with HTTP/1.1 only (HTTP/2 disabled)
			b.Run("HTTPS1.1", func(b *testing.B) {
				fs := http.FileServer(http.Dir(dir))
				ts := httptest.NewUnstartedServer(noCache(fs))
				ts.Config.ErrorLog = silentLogger
				ts.TLS = &tls.Config{
					Certificates: []tls.Certificate{cert},
					NextProtos:   []string{"http/1.1"},
				}
				ts.StartTLS()
				defer ts.Close()
				client := ts.Client()
				client.Transport.(*http.Transport).TLSClientConfig.NextProtos = []string{"http/1.1"}
				b.ResetTimer()
				for b.Loop() {
					resp, err := client.Get(ts.URL + f.path)
					if err != nil {
						b.Fatal(err)
					}
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
				}
			})

			// HTTPS with HTTP/2 (default)
			b.Run("HTTP2", func(b *testing.B) {
				fs := http.FileServer(http.Dir(dir))
				ts := httptest.NewUnstartedServer(noCache(fs))
				ts.Config.ErrorLog = silentLogger
				ts.TLS = &tls.Config{
					Certificates: []tls.Certificate{cert},
				}
				ts.StartTLS()
				defer ts.Close()
				client := ts.Client()
				b.ResetTimer()
				for b.Loop() {
					resp, err := client.Get(ts.URL + f.path)
					if err != nil {
						b.Fatal(err)
					}
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
				}
			})
		})
	}
}

// --- Self-signed cert generation benchmark ---

func BenchmarkGenerateSelfSignedCert(b *testing.B) {
	for b.Loop() {
		_, err := GenerateSelfSignedCert()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// --- Concurrent connection scaling: HTTP vs HTTPS/1.1 vs HTTP/2 ---

func BenchmarkProtocolConcurrency(b *testing.B) {
	dir := setupTestDir(b)

	cert, err := GenerateSelfSignedCert()
	if err != nil {
		b.Fatal(err)
	}

	for _, concurrency := range []int{1, 10, 50} {
		b.Run(fmt.Sprintf("clients-%d", concurrency), func(b *testing.B) {

			b.Run("HTTP", func(b *testing.B) {
				fs := http.FileServer(http.Dir(dir))
				ts := httptest.NewServer(noCache(fs))
				defer ts.Close()
				client := newKeepAliveClient()
				b.SetParallelism(concurrency)
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						resp, err := client.Get(ts.URL + "/small.html")
						if err != nil {
							b.Fatal(err)
						}
						io.Copy(io.Discard, resp.Body)
						resp.Body.Close()
					}
				})
			})

			b.Run("HTTPS1.1", func(b *testing.B) {
				fs := http.FileServer(http.Dir(dir))
				ts := httptest.NewUnstartedServer(noCache(fs))
				ts.Config.ErrorLog = silentLogger
				ts.TLS = &tls.Config{
					Certificates: []tls.Certificate{cert},
					NextProtos:   []string{"http/1.1"},
				}
				ts.StartTLS()
				defer ts.Close()
				client := ts.Client()
				client.Transport.(*http.Transport).TLSClientConfig.NextProtos = []string{"http/1.1"}
				client.Transport.(*http.Transport).MaxIdleConnsPerHost = concurrency
				b.SetParallelism(concurrency)
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						resp, err := client.Get(ts.URL + "/small.html")
						if err != nil {
							b.Fatal(err)
						}
						io.Copy(io.Discard, resp.Body)
						resp.Body.Close()
					}
				})
			})

			b.Run("HTTP2", func(b *testing.B) {
				fs := http.FileServer(http.Dir(dir))
				ts := httptest.NewUnstartedServer(noCache(fs))
				ts.Config.ErrorLog = silentLogger
				ts.TLS = &tls.Config{
					Certificates: []tls.Certificate{cert},
				}
				ts.StartTLS()
				defer ts.Close()
				client := ts.Client()
				b.SetParallelism(concurrency)
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						resp, err := client.Get(ts.URL + "/small.html")
						if err != nil {
							b.Fatal(err)
						}
						io.Copy(io.Discard, resp.Body)
						resp.Body.Close()
					}
				})
			})
		})
	}
}

// --- Throughput benchmark (reports bytes/sec) ---

func BenchmarkThroughput(b *testing.B) {
	dir := setupTestDir(b)
	fs := http.FileServer(http.Dir(dir))

	ts := httptest.NewServer(noCache(fs))
	defer ts.Close()

	client := newKeepAliveClient()

	// Get file size for throughput calculation
	resp, err := client.Get(ts.URL + "/large.bin")
	if err != nil {
		b.Fatal(err)
	}
	size := resp.ContentLength
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	b.SetBytes(size)
	b.ResetTimer()
	for b.Loop() {
		resp, err := client.Get(ts.URL + "/large.bin")
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

// --- Atomic counter benchmark ---

func BenchmarkIncrementRequest(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			incrementRequest()
		}
	})
}

// --- GetLocalIP benchmark ---

func BenchmarkGetLocalIP(b *testing.B) {
	_, err := GetLocalIP()
	if err != nil {
		b.Skip("no network available")
	}
	b.ResetTimer()
	for b.Loop() {
		GetLocalIP()
	}
}

// --- Directory listing benchmark ---

func BenchmarkDirectoryListing(b *testing.B) {
	dir := setupTestDir(b)
	fs := http.FileServer(http.Dir(dir))

	ts := httptest.NewServer(noCache(fs))
	defer ts.Close()

	client := newKeepAliveClient()

	b.ResetTimer()
	for b.Loop() {
		resp, err := client.Get(ts.URL + "/")
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

// --- Connection overhead (new TCP conn per request) ---

func BenchmarkNewConnPerRequest(b *testing.B) {
	dir := setupTestDir(b)
	fs := http.FileServer(http.Dir(dir))

	ts := httptest.NewServer(noCache(fs))
	defer ts.Close()

	b.ResetTimer()
	for b.Loop() {
		transport := &http.Transport{DisableKeepAlives: true}
		client := &http.Client{Transport: transport}
		resp, err := client.Get(ts.URL + "/small.html")
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		transport.CloseIdleConnections()
	}
}
