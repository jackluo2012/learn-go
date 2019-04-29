package main

import (
	"gopcp.v2/chapter7/rpc"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

/**


telnet localhost 1234
Trying ::1...
Connected to localhost.
Escape character is '^]'.
{"method":"DemoService.Div","params":[{"A":3,"B":0}],"id":1234}
{"id":1234,"result":null,"error":"Zero "}
{"method":"DemoService.Div","params":[{"A":3,"B":2}],"id":1}
{"id":1,"result":1.5,"error":null}

 */

func main() {
	rpc.Register(rpcdemo.DemoService{})
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept err %v", err)
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
}
