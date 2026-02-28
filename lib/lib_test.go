package lib

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestIp(t *testing.T) {
	_, err := GetLocalIP()
	if err != nil {
		t.Errorf("GetLocalIP isn't working: %s", err)
	}
}
func TestServerWithoutCache(t *testing.T) {
	go StartServer(".", ":80", false, false, "", "")
	resp, err := http.Get("http://127.0.0.1/")
	if err != nil {
		t.Errorf("Couldn't get for the test %s", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.Header.Get("Cache-Control") == "" {
		t.Error("Server not setting cache control")
	}
	if resp.StatusCode != 200 {
		t.Error("Server not returning 200")
	}
	if body == nil {
		t.Error("Body is nil, server not working")
	}
}

func TestInitialPrint(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printStartMessage("/", ":80", false)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if !strings.Contains(string(out), "go-live") {
		t.Error("Did not find go-live")
	}

	if !strings.Contains(string(out), "Serving: /") {
		t.Error("Did not find serving message")
	}

	if !strings.Contains(string(out), "http://localhost:80/") {
		t.Error("Did not find local network print")
	}
}

func TestInitialPrintHTTPS(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printStartMessage("/", ":9000", true)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if !strings.Contains(string(out), "go-live (HTTPS)") {
		t.Error("Did not find HTTPS indicator")
	}

	if !strings.Contains(string(out), "https://localhost:9000/") {
		t.Error("Did not find HTTPS local URL")
	}
}

func TestPrintServerInfo(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printServerInformation("/", ":80", false)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if !strings.Contains(string(out), "Net: http://") {
		t.Error("Did not net address")
	}

	if !strings.Contains(string(out), "Requests") {
		t.Error("Did not find requests")
	}
}

func TestPrintServerInfoHTTPS(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printServerInformation("/", ":9000", true)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if !strings.Contains(string(out), "Net: https://") {
		t.Error("Did not find HTTPS net address")
	}
}

func TestGenerateSelfSignedCert(t *testing.T) {
	cert, err := GenerateSelfSignedCert()
	if err != nil {
		t.Fatalf("GenerateSelfSignedCert failed: %s", err)
	}
	if len(cert.Certificate) == 0 {
		t.Error("Certificate is empty")
	}
	if cert.PrivateKey == nil {
		t.Error("PrivateKey is nil")
	}
}
