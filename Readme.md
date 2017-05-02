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
	then import GoWebsocket to your file
    ```GO
    package main
    import(
    	...
        "github.com/TorkilSo/gowebsocket/websocket"
        ...
    )
    ```
* Alternatively you can download GoWebsocket and add it to your GOPATH manually
 
 *Make sure your GOPATH looks something like this* 
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
 ```GO
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

*Example backend GO*

*Parameters for OnRecieve, OnOpen, OnClose and OnError needs to be as shown in this example*
```GO
    func main(){
        server := websocket.Create("localhost", "3001")
        server.Start()

        clients := server.GetClients()
	
	
	server.OnReceive = func(msg[]byte, client net.Conn){
		//do stuff with message and/or client
		//for example send messeage to all clients
		server.SendToAll(msg) 
		//or to specific client
		server.send(msg,client)
	}
	
	server.OnOpen = func(client net.Conn){
		//do stuff on connection 
		//for example closing connection
		client.close()
	}
	
	server.OnClose = func(client net.Conn){
		//do stuff on closing connection
		server.Send_string("Goodbye")
	}
	
	server.OnError = func(err string){
		//something something error handling
	}
    }
```
*Example frontend JAVASCRIPT*
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
### methods 
<dl>
<dt><strong>Create(host,port)</strong></dt>
<dd>Creates an instance of GoWebsocket on the specified host and port</dd>
<dt><strong>Start()</strong></dt>
<dd>Websocket starts listening </dd>
<dt><strong>SendString(string,net.Conn)</strong></dt>
<dd>Sends given string to specified client </dd>
<dt><strong>SendStringAsJSON(string,net.Conn)</strong></dt>
<dd>Sends given string in the format {"msg":msg} to specified client</dd>
<dt><strong>SendStringToAll(string)</strong></dt>
<dd>Sends given string to all clients on websocket</dd>
<dt><strong>SendStringToAllAsJSON(string)</strong></dt>
<dd>Sends given string in the format {"msg":msg} to all clients on websocket</dd>
<dt><strong>GetLatestMsg()</strong></dt>
<dd>returns the latest message sent over websocket</dd>
<dt><strong>GetLatestCLient()</strong></dt>
<dd>returns the latest client connected to the websocket</dd>
<dt><strong>ToString([]byte)</strong></dt>
<dd>converts byte-array to string 
</dd>
</dl>

