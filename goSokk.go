package main

import (
	"net"
	"net/http"
	"bufio"
	"strings"
)

//socket structure
/*
map of connections map [*ws_conn] bool?
broadcast chan [] byte
register chan * ws_conn
unregister chan * ws_conn
 */

var (

	//required features of an http request the initial handshake?? or is this the 2nd
	hsHead = map[string] bool {
		"Host": true,
		"Upgrade":true,
		"Connection":true,
		"Sec-Websocket-Key":true,
		"Sec-Websocket-Origin":true,
		"Sec-Websocket-Version":true,
		"Sec-Websocket-Protocol":true,
		"Sec-Websocket-Accept":true,
	}
)


/*
The main socket struct needs?
http switching protocols / upgrade
match bytes from key [MAGIC STRING]
handle new connection
handle data from connection
handle data from connections
handle connection closed


 */
type go_sokk struct {



}

type ws_conn struct {
	conn *net.Conn

}

// define interface for specified methods that must be implemented for a websocket connection
type ws_c_i interface {
	Read(b [] byte) (n int, err error)
	Write(b [] byte) (n int, err error)
	Close() error
}

//does this need to return anything?
func (w *go_sokk) ws_handshake(reader *bufio.Reader, request * http.Request) (code int){ // statuscode

	//check for HTTP GET method
	if request.Method != "GET"{
		return http.StatusMethodNotAllowed // only support for GET calls
	}

	//Check if UPGRADE header exists and if its not requesting a websocket upgrade
	if request.Header.Get("Upgrade") != "websocket" {
		return http.StatusBadRequest
	}

	//the upgrade header does not exist, theres no point to carry on
	key := request.Header.Get("Sec-Websocket-Key")

	//if the key is empty, something went wrong somewhere
	if key == "" || key == " " || len(key) > 8 {   //TODO check number of bytes to be length of key? 16/32/69
		return http.StatusBadRequest
	}
	//TODO implement check or case/switch for Sec-WebSocket-Version ?? [13] IS THIS REQUIRED??!


	//TODO lag base-64 SHA1 res
	return http.StatusSwitchingProtocols // changing to websocket protocol
}