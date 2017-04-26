package websocket

import (
	"net"

	"bytes"
	"encoding/base64"
	"crypto/sha1"
)

var (
	magic_server_key = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

func validate_key(key string) string{
	h := sha1.New()
	h.Write([]byte(key+magic_server_key))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func switch_pc(conn net.Conn, key string){ // TODO -> check if buffwriter is better?
	var buff bytes.Buffer
	buff.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
	buff.WriteString("Connection: Upgrade\r\n")
	buff.WriteString("Upgrade: websocket\r\n")
	buff.WriteString("Sec-WebSocket-Accept:")
	buff.WriteString(validate_key(key) + "\r\n\r\n")
	conn.Write(buff.Bytes())

}
/*
import hashlib, base64
def handshake(client):
    data = client.recv(1024)
    #Get headers
    headers = get_headers(data)
    magic_mike = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
    #Key = sha-1(websocket-key from client + magic_mike), then return base64 encoding
    key=base64.b64encode(hashlib.sha1(str.encode(headers['Sec-WebSocket-Key']+magic_mike)).digest()) #encode b64 the sha1 of the key
    print(key)
    #Response
    response = "HTTP/1.1 101 Switching Protocols\r\n"
    response+= "Upgrade: websocket\r\n"
    response+= "Connection: Upgrade\r\n"
    response+= "Sec-WebSocket-Accept:"+key+"\r\n\r\n"
    client.send(response)

def get_headers (data):
    headers = {}
    lines = data.splitlines()
    for l in lines:
        parts = l.split(": ", 1)
        if len(parts) == 2:
            headers[parts[0]] = parts[1]
    return headers

#example data
data = "GET /chat HTTP/1.1\r\n"
data+="Host: example.com:8000\r\n"
data+="Upgrade: websocket\r\n"
data+="Connection: Upgrade\r\n"
data+="Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\n"
data+="Sec-WebSocket-Version: 13\r\n"

 */