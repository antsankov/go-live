package lib

import (
	"log"
	"net/http"
)

func StartServer(dir string, port string) {
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
