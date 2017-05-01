package main

import (
	"net/http"
	"websocket"
)

func main() {

	go websocket.Start("localhost", "3001")

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":3000", nil)
}

