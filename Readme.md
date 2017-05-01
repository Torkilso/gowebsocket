# GoWebsocket

## About
`GoWebsocket` is a websocket implemented in Go


## Prerequisites
* Install `Go`
https://golang.org/dl/
https://golang.org/doc/install

## Build
* Download GoWebsocket to your GOPATH
	```sh
    	go get github.com/TorkilSo/gowebsocket
    ```
	to use GoWebsocket you only need to import it to your file
    ```GO
    package main
    import(
    	...
        "github.com/TorkilSo/gowebsocket/websocket"
        ...
    )
    ```
* Alternatively you can download GoWebsocket and add it to your GOPATH manually
  Make sure your *GOPATH* looks something like this 
  ```
    GOPATH
    └───src
    	└───websocket
    │   │   	encoding.go
    │   │   	handshake.go
    │   │		websocketserver.go
    │
    │
	```

 then import websocket in your file
	```Go
		package main
    	import(
    	...
    	"websocket"
        ...
    	)
    ```


## Using GoWebsocket
Create a server by calling the method *Create(host, port)*
and start it with the method *Start()*

You can access all clients connected to your socket by calling *GetClients()*

Each client has a LocalAddr and a RemoteAddr which you can access

*See main.go for further examples*

*GO*
```GO
    func main(){
        server := websocket.Create("localhost", "3001")
        server.Start()

        clients := server.GetClients()
    }
```
*JAVASCRIPT*
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





