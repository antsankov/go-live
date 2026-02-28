package lib

import "os/exec"

func openBrowser(url string) {
	exec.Command("open", url).Start()
}
