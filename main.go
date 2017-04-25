package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"bufio"
	"net/textproto"
	"regexp"
	"strings"
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

	//for loop? hvis alle skal ha hver sin traad
	p()
	client.Close()
}

func handshake(client net.Conn) {
	key, found := parseKey(client)

	if !found {
		reject(client)
		return
	}



	p(key)
	return
}

func parseKey(client net.Conn) (string, bool) {
	bufReader := bufio.NewReader(client)
	tp := textproto.NewReader(bufReader)

	var headers []string
	var key string
	var keyFound bool

	for {
		line, _ := tp.ReadLine()
		headers = append(headers, line)

		if line == "" {
			break
		}
	}

	for i := 0; i < len(headers); i++ {
		s := strings.Split(headers[i], ": ")
		matchKey, errKey := regexp.MatchString(s[0], "Sec-WebSocket-Key")
		if errKey != nil {
			p(errKey)
		}
		if matchKey {
			if len(s) > 1 {
				keyFound = true
				key = s[1]
			}
		}
	}

	if keyFound {
		return key, true
	} else {
		return "", false
	}
}

func reject(client net.Conn) {
	reject := "HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\nIncorrect request"
	client.Write([]byte(reject))
	client.Close();
}
