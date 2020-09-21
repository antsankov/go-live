/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"./cmd"
)

func printer(start time.Time, dir string, port string) {

	for {
		fmt.Println("\033[2J")
		fmt.Println("go-live\n--")
		fmt.Println("Serving: " + dir)
		fmt.Println("Port: " + port)
		fmt.Println(time.Since(start).Round(time.Second))
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	cmd.Execute()

	start := time.Now()
	dir := "./static"
	port := ":3000"

	go printer(start, dir, port)

	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func printInitial() {
	log.Println("Done")
}
