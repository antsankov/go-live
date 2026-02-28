package lib

import (
	"os/exec"
	"strings"
)

func openBrowser(url string) {
	r := strings.NewReplacer("&", "^&")
	exec.Command("cmd", "/c", "start", r.Replace(url)).Start()
}
