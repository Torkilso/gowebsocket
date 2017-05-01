package websocket

import (
	"encoding/base64"
	"crypto/sha1"
  "bufio"
  "net"
  "net/http"
  "bytes"
)

const (
	magic_server_key = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

/*
 * Encodes the key with SHA1 and base 64.
 */
func encodeKey(str string)(key string){
	h:=sha1.New()
	h.Write([]byte(str))
	key = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}

/*
 * Extracts the key from an incoming request.
 */
func parseKey(client net.Conn) (code int, k string) {
	bufReader := bufio.NewReader(client)
	request, err := http.ReadRequest(bufReader)
	if err != nil {
		p(err)
	}

	//Send bad request if not websocket
	if request.Header.Get("Upgrade") != "websocket" {
		return http.StatusBadRequest, ""
	} else {
		//Extract key and return it
		key := request.Header.Get("Sec-Websocket-Key")
		return http.StatusSwitchingProtocols, key
	}
}

/*
 * Reject and close connection
 */
func reject(client net.Conn) {
	reject := "HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\nIncorrect request"
	client.Write([]byte(reject))
	client.Close();
}

/*
 * Creates handshake
 */
func handshake(client net.Conn) bool {
	//Parse the key from the connection
	status, key := parseKey(client)

	//Check the status
	if status != 101 {
		//reject
		reject(client)
    return false
	} else {
		//Complete handshake
		var t = encodeKey(key + magic_server_key)
		var buff bytes.Buffer
		buff.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
		buff.WriteString("Connection: Upgrade\r\n")
		buff.WriteString("Upgrade: websocket\r\n")
		buff.WriteString("Sec-WebSocket-Accept:")
		buff.WriteString(t + "\r\n\r\n")

		//Write handshake to client
		client.Write(buff.Bytes())
    return true
	}
}
