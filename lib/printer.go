package lib

import (
	"fmt"
	"os"
	"time"
)

func printStartMessage(path string, port string) {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
	fmt.Println("go-live\n--")
	fmt.Printf("Serving: %s\n", path)
	fmt.Printf("Local: http://localhost%s/\n", port)
}

// Printer prints out the information associated with the server on a loop.
func Printer(dir string, port string) {
	// Need to give time if there is a server error.
	time.Sleep(5 * time.Millisecond)
	start := time.Now()

	path, err := os.Getwd()
	// if there is an error or we are using a special path, use dir arg.
	if err != nil || dir != "./" {
		path = dir
	}

	printStartMessage(path, port)
	var lastIP string
	for {
		// Move to the fifth row, change if more print statements are added.
		fmt.Print("\033[5;1H")
		localIP, err := GetLocalIP()
		if err == nil {
			fmt.Printf("Net: http://%s%s/\n", localIP, port)
			if lastIP == "" {
				lastIP = localIP
			}
		} else {
			// This is case when connection disconnects
			if lastIP != localIP {
				printStartMessage(path, port)
				fmt.Print("\033[6;1H")
				lastIP = localIP
			} else {
				lastIP = ""
				fmt.Println()
			}
		}
		fmt.Println("\nRequests:", requests)
		fmt.Println(time.Since(start).Round(time.Second))
		time.Sleep(250 * time.Millisecond)
	}
}
