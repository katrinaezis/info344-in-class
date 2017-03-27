package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// * passes by reference, no star passes by value
func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	// must add headers before writing body
	w.Header().Add("Content-Type", "text/plain")

	w.Write([]byte("hello " + name))
}

func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		log.Fatal("please set ADDR enviro var")
	}

	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("server is listening at %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
