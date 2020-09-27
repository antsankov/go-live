package lib

import (
	"fmt"
	"time"
)

func Printer(dir string, port string, url string) {
	start := time.Now()
	local, err1 := GetLocalIP()
	external, err2 := GetExternalIP()
	for {
		fmt.Println("\033[2J")
		fmt.Println("go-live\n--")
		fmt.Println("Serving: " + dir)
		fmt.Println("Port: " + port)
		fmt.Println("URL: " + url)
		if err1 == nil {
			fmt.Printf("Local: http://%s%s\n", local, port)
		}
		if err2 == nil {
			fmt.Printf("Internet: http://%s%s\n", external, port)
		}
		fmt.Println(time.Since(start).Round(time.Second))
		time.Sleep(100 * time.Millisecond)
	}
}
