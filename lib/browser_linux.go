package lib

import "os/exec"

func openBrowser(url string) {
	exec.Command("xdg-open", url).Start()
}
