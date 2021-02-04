package main

import (
	"os"
	"testing"
)

func TestIsSudo(t *testing.T) {
	sudo := isSudo()
	if os.Geteuid() == 0 && sudo != true {
		t.Errorf("Sudo checker is returning true")
	}
}
