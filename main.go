package main

import (
	"fmt"
	"net/http"
	"./websocket"
	"bufio"
	"os"
	"runtime"
	"strings"
	"strconv"
)

var p = fmt.Println

func serverInterface(s *websocket.Websocketserver)  {
	for {
		p("1. Print clients")
		p("2. Ping client")
		p("3. Close Client")
		p("4. Print number of goroutines")

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		str := strings.Replace(text, "\n", "", -1)
		clients := s.GetClients()

		if str == "1" {

			for i := range clients {
				p(i+1,": Remote address: " + clients[i].RemoteAddr().String() + ", local address: " + clients[i].LocalAddr().String())
			}
		} else if str == "2" {
			p("Choose client:")
			for i := range clients {
				p(i+1,": Remote address: " + clients[i].RemoteAddr().String() + ", local address: " + clients[i].LocalAddr().String())
			}
			text, _ := reader.ReadString('\n')
			str := strings.Replace(text, "\n", "", -1)
			stri, err := strconv.ParseInt(str, 16, 16)

			if err != nil {
				p(err)
			} else {
				if stri > 0 && int(stri) < len(clients)+1 {
					s.PingClient(int(stri)-1)
				} else {
					p("Not valid!")
				}
			}

		} else if str == "3" {
			p("Choose client:")
			for i := range clients {
				p(i+1,": Remote address: " + clients[i].RemoteAddr().String() + ", local address: " + clients[i].LocalAddr().String())
			}
			text, _ := reader.ReadString('\n')
			str := strings.Replace(text, "\n", "", -1)
			stri, err := strconv.ParseInt(str, 16, 16)

			if err != nil {
				p(err)
			} else {
				if stri > 0 && int(stri) < len(clients)+1 {
					s.CloseClient(int(stri)-1)
				} else {
					p("Not valid!")
				}
			}

		} else if str == "4"{
			p(runtime.NumGoroutine())
		} else {
			p("Please enter a number from the list above!")
		}
		p("\n")
	}
}


func main() {


	server := websocket.Create("localhost", "3001")
	server.Start()

	go serverInterface(&server)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":3000", nil)
}

