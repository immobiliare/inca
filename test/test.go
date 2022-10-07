package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", httpRequestHandler)
	if err := http.ListenAndServeTLS(":65535", os.Args[1], os.Args[2], nil); err != nil {
		log.Fatalln(err.Error())
	}
}

func httpRequestHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, World!\n")
}
