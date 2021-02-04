package lib

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestIp(t *testing.T) {
	_, err := GetLocalIP()
	if err != nil {
		t.Errorf("GetLocalIP isn't working: %s", err)
	}
}
func TestServerWithoutCache(t *testing.T) {
	go StartServer(".", ":80", false)
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
