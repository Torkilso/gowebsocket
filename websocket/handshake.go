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
	buff.WriteString(validate_key(key) + "\r\n")
	conn.Write(buff.Bytes())

}