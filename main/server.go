package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 创建一个server接口
func NewServer(ip string, port int) *Server {
	Server := &Server{
		Ip:   ip,
		Port: port,
	}
	return Server
}
func (this *Server) Handler(conn net.Conn) {
	//当前连接的服务

}

// 启动服务器的接口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprint("%s;&d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("listen err", err)
		return
	}
	defer listener.Close()
	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept err", err)
			continue
		}
		//do handler
		go this.Handler(conn)
	}
	//close listen socket

}
