package service

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
)

func StartServer(port int) {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "Response from server %d", port)
	}

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server at %s\n", addr)
	if err := fasthttp.ListenAndServe(addr, requestHandler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
