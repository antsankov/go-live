package lib

import (
	"fmt"
	"os"
	"time"
)

func scheme(tlsEnabled bool) string {
	if tlsEnabled {
		return "https"
	}
	return "http"
}

func printStartMessage(path string, port string, tlsEnabled bool) {
	// Clear the screen.
	fmt.Print("\033[2J")
	// Move the cursor to the upper-left corner of the screen.
	fmt.Print("\033[H")
	if tlsEnabled {
		fmt.Printf("go-live (HTTPS)\n--\n")
	} else {
		fmt.Printf("go-live\n--\n")
	}
	fmt.Printf("Serving: %s\n", path)
	fmt.Printf("Local: %s://localhost%s/\n", scheme(tlsEnabled), port)
}

func printServerInformation(path string, port string, tlsEnabled bool) {
	// Move to the fifth row, 1st column change if more print statements are added.
	fmt.Print("\033[5;1H")
	localIP, err := GetLocalIP()
	if err == nil && localIP != "" {
		fmt.Printf("Net: %s://%s%s/\033[K\n", scheme(tlsEnabled), localIP, port)
	} else {
		// If there is no network connection, erase the line.
		fmt.Print("\033[K")
		fmt.Println()
	}
	fmt.Println("\nRequests:", requests)
}

// Printer prints out the information associated with the server on a loop.
func Printer(dir string, port string, tlsEnabled bool) {
	// Need to give time if there is a server error.
	time.Sleep(5 * time.Millisecond)
	start := time.Now()

	path, err := os.Getwd()
	// if there is an error or we are using a special path, use dir arg.
	if err != nil || dir != "./" {
		path = dir
	}

	printStartMessage(path, port, tlsEnabled)
	for {
		time.Sleep(500 * time.Millisecond)
		printServerInformation(path, port, tlsEnabled)
		// Move to the timeSince row, and clear it.
		fmt.Print("\033[8;1H")
		fmt.Print("\033[K")
		fmt.Printf("%s\n", time.Since(start).Round(time.Second))
	}
}
