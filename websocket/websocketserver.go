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
	OnRecieve fn
	OnError fn
	OnOpen fn
	OnClose fn
}

type fn func()

// PingClient() does not work
//
/*func (s *Websocketserver) PingClient(c int)  {
	p("Pinging client...")
	ping := make([]byte,128)
	ping[0] = byte(137)
	s.clients[c].Write(ping)

	pong := make([]byte, 128)
	s.clients[c].Read(pong)

	b := fmt.Sprintf("%08b", byte(pong[0]))

	if b[4:len(b)] == "1010" {
		p("Got pong")
	} else {
		p("No pong")
	}
}*/

/*
 * Returns a new Websocketserver struct
 */
func Create(host string, port string) Websocketserver {
	return Websocketserver{clients: make([]net.Conn, 0), CONN_HOST: host, CONN_PORT: port, OnRecieve: func(){}, OnOpen: func(){}, OnClose: func(){}, OnError: func(){}}
}

/*
 * Returns all the clients for a Websocketserver struct
 */
func (s *Websocketserver) GetClients() []net.Conn {
	return s.clients
}

func (s *Websocketserver) Send(msg []byte, client net.Conn){
	decoded := decode(msg)
	enc := encode(decoded)
	client.Write(enc)
}

/*
 * Writes incoming message to all connections to a Websocketserver struct
 */
func (s *Websocketserver) SendToAll(msg []byte) {
	decoded := decode(msg)
	enc := encode(decoded)
  for i := range s.clients {
    s.clients[i].Write(enc)
  }
}

func (s *Websocketserver) GetLatestMsg() []byte {
	return s.latestMsg
}
func (s *Websocketserver) GetLatestClient() net.Conn {
	return s.latestClient
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

		go s.OnOpen()

		s.latestClient = *client
		//Listen for incoming messages
    for {
			//4KB buffer
      msg := make([]byte, 4096)

			//Waits for incoming message
      (*client).Read(msg)

			//Extract the first byte to find opcode
      c := fmt.Sprintf("%08b", byte(msg[0]))

			//If client sent close
      if c[4:len(c)] == "1000" {
				go s.OnClose()
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
				s.latestClient = *client
				go s.OnRecieve();

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
		s.OnError()
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
			go s.OnError()
			p("Error accepting: ", err.Error())
			os.Exit(1)
		}

		// Handle connections in a new thread (goroutine)
		go handler(&conn, s)
	}
}
