package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	//消息广播的channel
	Message chan string
}

// 创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听Message广播消息channel的gorountine，一旦有消息就发送给全部user
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		//将msg发送给全部的user
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + "说:" + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//用户上线 将用户加入到onlineMap中
	user := NewUser(conn, this)
	user.Online()
	//广播当前用户上线消息
	this.BroadCast(user, "已上线")
	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read err", err)
				return
			}
			//提取用户的消息 去除\n
			msg := string(buf[:n])
			//将消息进行广播
			user.DoMessage(msg)
		}
	}()
	//当前handler阻塞
	select {}
}

// 启动服务器的接口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprint("%s;&d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("listen err", err)
		return
	}
	//close listen socket
	defer listener.Close()
	//启动监听Message的goroutine
	go this.ListenMessage()
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

}
