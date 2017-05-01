# GoWebsocket

## About
`GoWebsocket` is a websocket implemented in Go

## Prerequisites
* Install `Go`
https://golang.org/dl/
https://golang.org/doc/install

## Build
* Download websocket to project directory
  ```
    project
    │   README.md
    │   main.go
    │
    └───folder1
    │   │   site.html
    │   │   style.css
    │   │   ...
    │
    └───folder2
    │   │   go_server.go
    │   │   ...
    │
    └───websocket
    │   │   encode_test.go
    │   │   encoding.go
    │   │   handshake.go
    │   │	websocketserver.go
    │   │   ...
    │
	```

* Import websocket-folder
	```Go
		package main
    	import(
    	...
    	"./websocket"
        ...
    	)
    ```


## Using GoWebsocket
Create a server by calling the method *create(host, port)*
and start it with the method *start()*

You can access all clients connected to your socket by calling *GetClients()*

Each client has a LocalAddr and a RemoteAddr which you can access
```GO
    func main(){
        server := websocket.Create("localhost", "3001")
        server.Start()

        go example_function(&server)
        clients := server.GetClients()

        http.Handle("/", http.FileServer(http.Dir("./static")))
        http.ListenAndServe(":3000", nil)
    }
```
```javascript
    this.ws = new WebSocket("ws://localhost:3001");


    this.ws.onmessage =  (event) => {
    //do stuff here
    }

    this.ws.onerror = (event) => {

    }

    this.ws.onopen = (event) => {

    }

    this.ws.onclose = (event) => {

    }
```





