package service

import (
	"fmt"
	"log"
	"net/http"
)

func StartServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Response from server %d", port)
	})

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
