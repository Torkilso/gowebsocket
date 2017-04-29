package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"bufio"
	"bytes"
	"encoding/base64"
	"crypto/sha1"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3001"
	CONN_TYPE = "tcp"
	magic_server_key = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

var p = fmt.Println
var clients = make([]net.Conn, 0)

func hand(str string)(key string){
	h:=sha1.New()
	h.Write([]byte(str))
	key = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}

func parseKey(client net.Conn) (code int, k string) {
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

//Funnet på nett

func decode (inputBytes []byte) string {
	mask := 2
	if inputBytes[1]-128 == 126 {
		mask = 4
	} else if inputBytes[1]-128 == 127 {
		mask = 10
	}

	masks := inputBytes[mask:mask + 4]

  dataEnd := mask + 4

  for r := mask + 4; r <= len(inputBytes); r++ {
    if inputBytes[r] == 0 {
      check := true
      for t := 1; t <= 10; t++ {
        if inputBytes[r+t] != 0 {
          check = false
          break
        }
      }
      if check {
        dataEnd = r
        break
      }
    }
  }

  data := inputBytes[mask + 4:dataEnd]
	decoded := make([]byte, len(data))

	for i, b := range data {
		decoded[i] = b ^ masks[i % 4]
	}
	return string(decoded)
}

func encode (message string) (result []byte) {

	input := []byte(message)
	var dataIndex int

	length := byte(len(input))

	if len(input) <= 125 { //one byte to store data length
		result = make([]byte, len(input)+2)
		result[1] = length
		dataIndex = 2
	} else if len(input) >= 126 && len(input) <= 65535 { //two bytes to store data length
		result = make([]byte, len(input)+4)

		result[1] = 126 //extra storage needed
		result[2] = byte(len(input) >> 8)
		result[3] = length

		dataIndex = 4
	} else {
		result = make([]byte, len(input)+10)
		result[1] = 127
		result[2] = byte(len(input) >> 56)
		result[3] = byte(len(input) >> 48)
		result[4] = byte(len(input) >> 40)
		result[5] = byte(len(input) >> 32)
		result[6] = byte(len(input) >> 24)
		result[7] = byte(len(input) >> 16)
		result[8] = byte(len(input) >> 8)
		result[9] = length
		dataIndex = 10
	}

	result[0] = 129

	// put data at the correct index
	for i, b := range input {
		result[dataIndex+i] = b
	}
	return result
}

func handshake(client net.Conn) bool {
	status, key := parseKey(client)
	if status != 101 {
		//reject
		reject(client)
    return false
	} else {
		//Complete handshake
		var t = hand(key + magic_server_key)
		var buff bytes.Buffer
		buff.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
		buff.WriteString("Connection: Upgrade\r\n")
		buff.WriteString("Upgrade: websocket\r\n")
		buff.WriteString("Sec-WebSocket-Accept:")
		buff.WriteString(t + "\r\n\r\n")
		client.Write(buff.Bytes())
    return true
	}
}

func handleMsg(msg []byte) {
  decoded := decode(msg)
  enc := encode(decoded)
  writeToAll(enc)
}

func writeToAll(msg []byte) {
  for i := range clients {
    clients[i].Write(msg)
  }
}

func closeConn(client net.Conn) {
  for i := range clients {
    if clients[i] == client {
      clients = clients[:i + copy(clients[i:], clients[i+1:])]
      client.Close()
      break
    }
  }
  client.Close()
}

func handler(client net.Conn) {
	verified := handshake(client)
  if(verified){

    clients = append(clients, client)

    for {
      msg := make([]byte, 4096)
      client.Read(msg)

      c := fmt.Sprintf("%08b", byte(msg[0]))
      if c[4:len(c)] == "1000" {
        closeConn(client)
        break
      }else if c[4:len(c)] == "1001"{
	      response := make([]byte,2)
	      response[0] = byte(138)
	      // p(fmt.Sprintf("%08b",byte(response[0])))
	      client.Write(response)
      }else{
	      go handleMsg(msg)
      }
    }
  }
}

func startWss() {
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

func main() {
	go startWss()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":3000", nil)



}
