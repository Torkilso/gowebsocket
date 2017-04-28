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
  //"runtime"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3001"
	CONN_TYPE = "tcp"
	magic_server_key = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
  msg_length = 128
)

var p = fmt.Println

// Handles incoming requests.

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
/*
Første byte inneholder typ beskrivelse
Andre byte inneholder lengden på dataen fra(/til) klienten
either two or eight bytes if the length does not fit in the second byte (the second byte is then a code saying how many bytes are used for the length)
the actual (raw) data
 */
func decode (rawBytes []byte) string {
	var idxMask int
	if rawBytes[1] == 126 {
		idxMask = 4
	} else if rawBytes[1] == 127 {
		idxMask = 10
	} else {
		idxMask = 2
	}

	masks := rawBytes[idxMask:idxMask + 4]

  first := 6

  for r := range rawBytes{
    if rawBytes[r] == 0 {
      first = r
      break
    }
  }

  if first < 6 {
    first = 6
  }

  data := rawBytes[idxMask + 4:first]
	decoded := make([]byte, len(data))

	for i, b := range data {
		decoded[i] = b ^ masks[i % 4]
	}
	return string(decoded)
}

//Ikke testet 27.04.17
func encode (message string) (result []byte) {
	rawBytes := []byte(message)
	var idxData int

	length := byte(len(rawBytes))
	if len(rawBytes) <= 125 { //one byte to store data length
		result = make([]byte, len(rawBytes)+2)
		result[1] = length
		idxData = 2
	} else if len(rawBytes) >= 126 && len(rawBytes) <= 65535 { //two bytes to store data length
		result = make([]byte, len(rawBytes)+4)
		result[1] = 126 //extra storage needed
		result[2] = ( length >> 8 ) & 255
		result[3] = ( length      ) & 255
		idxData = 4
	} else {
		result = make([]byte, len(rawBytes)+10)
		result[1] = 127
		result[2] = ( length >> 56 ) & 255
		result[3] = ( length >> 48 ) & 255
		result[4] = ( length >> 40 ) & 255
		result[5] = ( length >> 32 ) & 255
		result[6] = ( length >> 24 ) & 255
		result[7] = ( length >> 16 ) & 255
		result[8] = ( length >> 8 ) & 255
		result[9] = ( length       ) & 255
		idxData = 10
	}

	result[0] = 129 //only text is supported

	// put raw data at the correct index
	for i, b := range rawBytes {
		result[idxData+i] = b
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

func handleIncomingMsg(msg []byte, client net.Conn) {
  c := fmt.Sprintf("%08b", byte(msg[0]))
  if c[4:len(c)] == "1000" {
    closeConn(client)
  }

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
      msg := make([]byte, 32)
      client.Read(msg)
	    c := fmt.Sprintf("%08b", byte(msg[0]))
      select {
      case c[4:len(c)] == "1000":
	      closeConn(client)
	      break
      case c[4:len(c)] == "1001":
	      //PONG
      default:
	      decoded := decode(msg)
	      enc := encode(decoded)

	      writeToAll(enc)
      }
    }
  }
}

var clients = make([]net.Conn, 0)

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
  fmt.Sprintf("0%8b", byte(1))
	go startWss()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":3000", nil)
}
