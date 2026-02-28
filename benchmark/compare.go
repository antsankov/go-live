// Benchmark comparing go-live vs caddy vs miniserve vs live-server.
// Run: go run benchmark/compare.go
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

const (
	duration = 3 * time.Second
	warmup   = 500 * time.Millisecond
)

type server struct {
	name  string
	port  string
	start func(dir, port string) *exec.Cmd
}

func main() {
	dir := setupTestDir()
	defer os.RemoveAll(dir)

	goLiveBin := buildGoLive()
	caddyBin := findBin("caddy", filepath.Join(os.Getenv("HOME"), "go", "bin", "caddy"))
	miniserveBin := findBin("miniserve", "")
	liveServerBin := findBin("live-server", "")

	servers := []server{
		{
			name: "go-live",
			port: "9100",
			start: func(dir, port string) *exec.Cmd {
				return exec.Command(goLiveBin, "-q", "-p", port, "-d", dir)
			},
		},
	}

	if caddyBin != "" {
		servers = append(servers, server{
			name: "caddy",
			port: "9101",
			start: func(dir, port string) *exec.Cmd {
				return exec.Command(caddyBin, "file-server", "--listen", ":"+port, "--root", dir)
			},
		})
	}

	if miniserveBin != "" {
		servers = append(servers, server{
			name: "miniserve",
			port: "9102",
			start: func(dir, port string) *exec.Cmd {
				return exec.Command(miniserveBin, "-p", port, dir)
			},
		})
	}

	if liveServerBin != "" {
		servers = append(servers, server{
			name: "live-server",
			port: "9103",
			start: func(dir, port string) *exec.Cmd {
				return exec.Command(liveServerBin, dir, "--port="+port, "--no-browser", "--quiet")
			},
		})
	}

	files := []struct {
		name string
		path string
	}{
		{"1KB", "/small.html"},
		{"50KB", "/medium.bin"},
		{"1MB", "/large.bin"},
	}

	concurrencies := []int{1, 10, 50}

	fmt.Println("File server benchmark")
	fmt.Println("=====================")
	fmt.Printf("Duration per test: %s\n", duration)
	fmt.Printf("Servers: ")
	for i, s := range servers {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(s.name)
	}
	fmt.Println()
	fmt.Println()

	// Print header
	fmt.Printf("%-8s  %-7s", "File", "Clients")
	for _, s := range servers {
		fmt.Printf("  %15s", s.name)
	}
	fmt.Println()

	fmt.Printf("%-8s  %-7s", "----", "-------")
	for range servers {
		fmt.Printf("  %15s", "-------")
	}
	fmt.Println()

	for _, f := range files {
		for _, c := range concurrencies {
			results := make([]int64, len(servers))

			for i, s := range servers {
				cmd := s.start(dir, s.port)
				cmd.Stdout = io.Discard
				cmd.Stderr = io.Discard
				cmd.Start()

				url := fmt.Sprintf("http://127.0.0.1:%s%s", s.port, f.path)
				waitForServer(url)
				results[i] = benchmark(url, c)

				cmd.Process.Kill()
				cmd.Wait()
				time.Sleep(200 * time.Millisecond)
			}

			fmt.Printf("%-8s  %-7d", f.name, c)
			for _, r := range results {
				fmt.Printf("  %12d/s", r)
			}
			fmt.Println()
		}
	}
}

func setupTestDir() string {
	dir, _ := os.MkdirTemp("", "bench-compare-*")

	small := make([]byte, 1024)
	for i := range small {
		small[i] = 'a' + byte(i%26)
	}
	os.WriteFile(filepath.Join(dir, "small.html"), small, 0644)

	// 50KB file
	medium := make([]byte, 50*1024)
	for i := range medium {
		medium[i] = byte(i % 256)
	}
	os.WriteFile(filepath.Join(dir, "medium.bin"), medium, 0644)

	// 1MB file
	large := make([]byte, 1024*1024)
	for i := range large {
		large[i] = byte(i % 256)
	}
	os.WriteFile(filepath.Join(dir, "large.bin"), large, 0644)

	return dir
}

func buildGoLive() string {
	fmt.Print("Building go-live... ")
	out := filepath.Join(os.TempDir(), "go-live-bench")
	cmd := exec.Command("go", "build", "-ldflags", "-s -w", "-trimpath", "-o", out, ".")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build go-live: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("done")
	return out
}

func findBin(name string, fallback string) string {
	path, err := exec.LookPath(name)
	if err == nil {
		return path
	}
	if fallback != "" {
		if _, err := os.Stat(fallback); err == nil {
			return fallback
		}
	}
	fmt.Fprintf(os.Stderr, "warning: %s not found, skipping\n", name)
	return ""
}

func waitForServer(url string) {
	client := &http.Client{Timeout: 200 * time.Millisecond}
	for i := 0; i < 50; i++ {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Fprintf(os.Stderr, "warning: server at %s did not start in time\n", url)
}

// benchmark returns requests per second.
func benchmark(url string, concurrency int) int64 {
	transport := &http.Transport{
		MaxIdleConnsPerHost: concurrency,
		DisableKeepAlives:   false,
	}
	client := &http.Client{Transport: transport}

	// Warmup
	for i := 0; i < concurrency*2; i++ {
		resp, err := client.Get(url)
		if err == nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
	time.Sleep(warmup)

	var count atomic.Int64
	var wg sync.WaitGroup
	stop := make(chan struct{})

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					resp, err := client.Get(url)
					if err != nil {
						continue
					}
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
					count.Add(1)
				}
			}
		}()
	}

	time.Sleep(duration)
	close(stop)
	wg.Wait()

	return count.Load() / int64(duration.Seconds())
}
