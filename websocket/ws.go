package websocket

import (
    "fmt"
    "net"
    "os"
    "bufio"
    "net/textproto"
)

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
        // Handle connections in a new goroutine.
        go handler(conn)
    }
}

// Handles incoming requests.
func handler(conn net.Conn) {
  shakeIt(conn)

  p("Got connection from:",conn)

	bufReader := bufio.NewReader(conn)
  tp := textproto.NewReader(bufReader)

  var headers []string

  for {
      line, _ := tp.ReadLine()
      headers = append(headers, line)

      if line == "" {
          break
      }
  }

  p(headers)
  conn.Write([]byte("Hei"))
  conn.Close()
}

func shakeIt(client net.Conn) {
  var shake = "HTTP/1.1 101 Web Socket Protocol Handshake\r\n"
	client.Write([]byte(shake))
}
