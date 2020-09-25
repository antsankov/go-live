package lib

import (
	"fmt"
	"time"
)

func Printer(dir string, port string) {
	start := time.Now()
	for {
		fmt.Println("\033[2J")
		fmt.Println("go-live\n--")
		fmt.Println("Serving: " + dir)
		fmt.Println("Port: " + port)
		fmt.Println(time.Since(start).Round(time.Second))
		time.Sleep(100 * time.Millisecond)
	}
}
