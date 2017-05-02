package main

import (
	"fmt"
	"net/http"
	"github.com/TorkilSo/gowebsocket/websocket"
	"bufio"
	"os"
	"runtime"
	"strings"
	"net"
	//"strconv"
)

var p = fmt.Println

func serverInterface(s *websocket.Websocketserver)  {
	for {
		p("1. Print clients")
		p("2. Print number of connections")
		p("3. Print number of goroutines")

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		str := strings.Replace(text, "\n", "", -1)
		clients := s.GetClients()

		if str == "1" {

			for i := range clients {
				//Prints out remote address and local address of a client
				p(i+1,": Remote address: " + clients[i].RemoteAddr().String() + ", local address: " + clients[i].LocalAddr().String())
			}
		} else if str == "2"{
			//Number of clients is the length of s.GetClients()
			p("Number of connections: ", len(clients))
		} else if str == "3"{
			p("Number of Goroutines: ", runtime.NumGoroutine())
		} else {
			p("Please enter a number from the list above!")
		}
		p("\n")
	}
}



func main() {

	//Creates a websocketserver "object"
	server := websocket.Create("localhost", "3001")

	server.OnRecieve = func (msg []byte, client net.Conn)  {
		p("New message: ", server.ToString(msg))
		server.SendToAll(msg)
	}

	server.OnOpen = func (client net.Conn)  {
		p("New connection!\nRemote address: ", client.RemoteAddr().String(), "Local address: ", client.LocalAddr().String())
	}

	server.OnClose = func (client net.Conn)  {
		p("Client closed!\nRemote address: ", client.RemoteAddr().String(), "Local address: ", client.LocalAddr().String())
	}

	server.OnError = func (err string)  {
		p(err)
	}

	//Start the server
	server.Start()

	go serverInterface(&server)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":3000", nil)
}
