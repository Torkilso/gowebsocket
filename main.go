package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"bufio"
	//"net/textproto"
	//"regexp"
	//"strings"
)

func main() {
	go startWss()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":3000", nil)
}

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3001"
	CONN_TYPE = "tcp"
)

var p = fmt.Println

func startWss() {
	// Listen for incoming connections.
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		p("Error listening:", err.Error())
		os.Exit(1)
	}

	//Executed when the application closes.
	defer listener.Close()

	p("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			p("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new thread (goroutine)
		go handler(conn)
	}
}

// Handles incoming requests.
func handler(client net.Conn) {
	handshake(client)
}

func handshake(client net.Conn) {
	status, key := parseKey(client)
  if status != 101 {
    //reject
    reject(client)
  } else {
    //Complete handshake

    p(key)

  }
}

func parseKey(client net.Conn)(code int, k string) {
  bufReader := bufio.NewReader(client)
  request, err := http.ReadRequest(bufReader)

  if err != nil {
      p(err)
  }

  if request.Header.Get("Upgrade") != "websocket" {
		return http.StatusBadRequest, ""
	} else {
    key := request.Header.Get("Sec-Websocket-Key")
    return http.StatusSwitchingProtocols, key
  }
}

func reject(client net.Conn) {
	reject := "HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\nIncorrect request"
	client.Write([]byte(reject))
	client.Close();
}
