package main

import (
	"log"
	"net/http"
	"webssh/lib"
)

func main() {

	http.Handle("/socket.io/", lib.InitServer())
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:3003...")
	log.Fatal(http.ListenAndServe(":3003", nil))
}
