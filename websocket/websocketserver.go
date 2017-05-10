package websocket

import (
	"fmt"
	"net"
	"os"
)

const (
	CONN_TYPE = "tcp"
)

var p = fmt.Println

type Websocketserver struct {
	clients []net.Conn
	CONN_HOST string
	CONN_PORT string
	latestMsg []byte
	latestClient net.Conn
	OnRecieve fn_b
	OnError fn_e
	OnOpen fn_c
	OnClose fn_c
}

/*
 * Types for functions in the Websocketserver struct
 */
type fn_b func([]byte, net.Conn)
type fn_c func(net.Conn)
type fn_e func(string)

/*
 * Returns a new Websocketserver struct
 */
func Create(host string, port string) Websocketserver {
	return Websocketserver{clients: make([]net.Conn, 0), CONN_HOST: host, CONN_PORT: port, OnRecieve: func(a []byte, b net.Conn){}, OnOpen: func(c net.Conn){}, OnClose: func(c net.Conn){}, OnError: func(err string){}}
}

/*
 * Returns all the clients for a Websocketserver struct
 */
func (s *Websocketserver) GetClients() []net.Conn {
	return s.clients
}

/*
 * Send message to client client, input in byteslice
 */
func (s *Websocketserver) Send(msg []byte, client net.Conn){
	decoded := decode(msg)
	enc := encode(decoded)
	client.Write(enc)
}

/*
 * Send message to one client, input in string
 */
func (s *Websocketserver) SendString(msg string, client net.Conn){
	enc := encode(msg)
	client.Write(enc)
}

/*
 * Send message to one client in the format {"msg":msg}, input in string
 */
func (s *Websocketserver) SendStringAsJSON(msg string, client net.Conn){
	msg_f := "{\"msg\":\"" + msg + "\"}"
	enc := encode(msg_f)
	client.Write(enc)
}

/*
 * Send message to all connected clients, input in string
 */
func (s *Websocketserver) SendStringToAll(msg string){
	enc := encode(msg)
	for i := range s.clients {
		s.clients[i].Write(enc)
	}
}

/*
 * Send message to all connected clients in the format {"msg":msg}, input in string
 */
func (s *Websocketserver) SendStringToAllAsJSON(msg string){
	msg_f := "{\"msg\":\"" + msg + "\"}"
	enc := encode(msg_f)
	for i := range s.clients {
		s.clients[i].Write(enc)
	}
}

/*
 * Send message to all connected clients, input in byteslice
 */
func (s *Websocketserver) SendToAll(msg []byte) {
	decoded := decode(msg)
	enc := encode(decoded)
  for i := range s.clients {
    s.clients[i].Write(enc)
  }
}

/*
 * Returns latest message recieved from all clients
 */
func (s *Websocketserver) GetLatestMsg() []byte {
	return s.latestMsg
}

/*
 * Returns latest client connected
 */
func (s *Websocketserver) GetLatestClient() net.Conn {
	return s.latestClient
}

/*
 * Decodes message in byteslice to string
 */
func (s *Websocketserver) ToString(msg []byte) string {
	return decode(msg)
}

/*
 * Starts the listener in a new goroutine
 */
func (s *Websocketserver)Start() {
	go listen(s)
}

/*
 * Closes a connection.
 * Removes the connection from the Websocketserver clients slice.
 */
func closeConn(client net.Conn, s *Websocketserver) {
  for i := range s.clients {
    if s.clients[i] == client {

			//Remove the client from the Websocketserver slice
      s.clients = s.clients[:i + copy(s.clients[i:], s.clients[i+1:])]

			//Closes the connection
      client.Close()
			go s.OnClose(client)
      break
    }
  }
}

/*
 * Removes client from clientslice in Websocketserver struct when error is read from the handler
 */
func closeConnErr(client net.Conn, s *Websocketserver) {
  for i := range s.clients {
    if s.clients[i] == client {

			//Remove the client from the Websocketserver slice
      s.clients = s.clients[:i + copy(s.clients[i:], s.clients[i+1:])]

			//Runs OnClose() specified by user of library
			go s.OnClose(client)
      break
    }
  }
}

/*
 * Handles incoming messages from a client.
 * Each client has a handler running in a seperate goroutine
 */
func handler(client *net.Conn, s *Websocketserver) {
	verified := handshake(*client)
	if(verified){

		//Add client to slice in Websocketserver
    s.clients = append(s.clients, *client)

		go s.OnOpen(*client)
		s.latestClient = *client

		//Listen for incoming messages
    for {
			//4KB buffer
      msg := make([]byte, 4096)

			//Waits for incoming message
      length, err_r := (*client).Read(msg)
			//Check for errors
			if err_r != nil {
				go s.OnError(err_r.Error())
				closeConnErr(*client,s)
				break
			}

			//Extract the first byte to find opcode
      c := fmt.Sprintf("%08b", byte(msg[0]))

			//If client sent close
      if c[4:len(c)] == "1000" {
        closeConn(*client, s)
        break

      }else if c[4:len(c)] == "1001"{
				//Ping

	      response := make([]byte,2)
	      response[0] = byte(138)

				//Sends back pong
	      (*client).Write(response)
      }else{
				//Handle message in a new goroutine
	      //go handleMsg(msg, s)
				s.latestMsg = msg
				go s.OnRecieve(msg[:length], *client);
      }
    }
  }
}

/*
 * Listens for new connections/clients.
 */
func listen(s *Websocketserver){

	listener, err := net.Listen(CONN_TYPE, s.CONN_HOST+":"+s.CONN_PORT)

	if err != nil {
		s.OnError(err.Error())
		p("Error listening:", err.Error())
		os.Exit(1)
	}

	//Executed when the application closes.
	defer listener.Close()

	p("Listening on " + s.CONN_HOST + ":" + s.CONN_PORT)

	for {

		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			go s.OnError(err.Error())
			//p("Error accepting: ", err.Error())
			//os.Exit(1)
		} else {
			// Handle connections in a new thread (goroutine)
			go handler(&conn, s)
		}
	}
}
