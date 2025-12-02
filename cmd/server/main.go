package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/modelcontrolinterface/mcim/internal/rpc"
)

func main() {
	r := rpc.NewRouter()

	fmt.Println("Starting server on :1234")
	log.Fatal(http.ListenAndServe(":1234", r))
}
