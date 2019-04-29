package main

import (
	"fmt"
	"gopcp.v2/chapter7/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	rpcClient, err := jsonrpc.Dial("tcp", ":1234")
	if err != nil {

	}
	var result float64
	err = rpcClient.Call("DemoService.Div", rpcdemo.Args{A: 10, B: 2}, &result)
	if err != nil {

	}
	fmt.Print(result)
}
