package rpc_support

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// 启动一个 rpc 的服务
func ServeRpc(host string, service interface{}) error {
	//注册
	rpc.Register(service)
	//监听端口
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}
	log.Printf("Listening on %s", host)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept err %v", err)
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
	return nil
}
// 调用一个客户端
func NewClient(host string) (*rpc.Client, error) {
	conn, err := net.Dial("tcp", host)
	return jsonrpc.NewClient(conn), err

}
