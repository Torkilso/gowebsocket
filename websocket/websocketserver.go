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
}

func (s *Websocketserver) GetClients() []net.Conn {
	return s.clients
}

func (s *Websocketserver) PingClient(c int)  {
	/*p("Pinging client...")
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
	}*/
}


func Create(host string, port string) Websocketserver {
	return Websocketserver{clients: make([]net.Conn, 0), CONN_HOST: host, CONN_PORT: port}
}

func handleMsg(msg []byte, s *Websocketserver) {
  decoded := decode(msg)
  enc := encode(decoded)
  writeToAll(enc, s)
}

func writeToAll(msg []byte, s *Websocketserver) {
  for i := range s.clients {
    s.clients[i].Write(msg)
  }
}

func closeConn(client net.Conn, s *Websocketserver) {
  for i := range s.clients {
    if s.clients[i] == client {
      s.clients = s.clients[:i + copy(s.clients[i:], s.clients[i+1:])]
      client.Close()
      break
    }
  }
}

func (s *Websocketserver) CloseClient(client int) {
	/*s.clients[client].Close()
	s.clients = s.clients[:client + copy(s.clients[client:], s.clients[client+1:])]
	*/
}

func handler(client *net.Conn, s *Websocketserver) {
	verified := handshake(*client)
  if(verified){

    s.clients = append(s.clients, *client)

    for {
      msg := make([]byte, 4096)
      (*client).Read(msg)
      c := fmt.Sprintf("%08b", byte(msg[0]))
      if c[4:len(c)] == "1000" {
        closeConn(*client, s)
        break
      }else if c[4:len(c)] == "1001"{
	      response := make([]byte,2)
	      response[0] = byte(138)
	      // p(fmt.Sprintf("%08b",byte(response[0])))
	      (*client).Write(response)
      }else{
	      go handleMsg(msg, s)
      }
    }
  }
}

func listen(s *Websocketserver){
	listener, err := net.Listen(CONN_TYPE, s.CONN_HOST+":"+s.CONN_PORT)

	if err != nil {
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
			p("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new thread (goroutine)
		go handler(&conn, s)
	}
}

func (s *Websocketserver)Start() {
	go listen(s)
}
