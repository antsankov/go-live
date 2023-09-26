package lib

import (
	"fmt"
	"os"
	"time"
)

// VERSION of Package
const VERSION = "1.2.0"

func printStartMessage(path string, port string) {
	// Clear the screen.
	fmt.Print("\033[2J")
	// Move the cursor to the upper-left corner of the screen.
	fmt.Print("\033[H")
	fmt.Printf("go-live\n--\n")
	fmt.Printf("Serving: %s\n", path)
	fmt.Printf("Local: http://localhost%s/\n", port)
}

func printServerInformation(path string, port string) {
	// Move to the fifth row, 1st column change if more print statements are added.
	fmt.Print("\033[5;1H")
	localIP, err := GetLocalIP()
	if err == nil && localIP != "" {
		fmt.Printf("Net: http://%s%s/\n", localIP, port)
	} else {
		// If there is no network connection, erase the line.
		fmt.Print("\033[K")
		fmt.Println()
	}
	fmt.Println("\nRequests:", requests)
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
	for {
		time.Sleep(100 * time.Millisecond)
		printServerInformation(path, port)
		// Move to the timeSince row, and clear it.
		fmt.Print("\033[8;1H")
		fmt.Print("\033[K")
		fmt.Printf("%s\n", time.Since(start).Round(time.Second))
	}
}
