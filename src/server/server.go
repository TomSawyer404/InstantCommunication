package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户表
	OnlineMap map[string]*User
	map_lock  sync.RWMutex

	// 用于广播的channel
	Message chan string
}

// Server的构造函数
func NewServer(ip string, port int) *Server {
	new_server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return new_server
}

// 监听Message广播消息的channel的goroutine，一旦有消息就发送给全部的在线user
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		//将msg发送给在线User
		this.map_lock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.map_lock.Unlock()
	}
}

func (this *Server) BroadCast(user *User, msg string) {
	send_msg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- send_msg
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("Connection established successfully.")

	new_user := NewUser(conn, this)
	new_user.Online()

	// 监听用户是否活跃的channel
	isLive := make(chan bool)

	// 接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			read_bytes, err := conn.Read(buf)
			if read_bytes == 0 {
				new_user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			// 提前用户消息（去除`\n`）
			msg := string(buf[:read_bytes-1])

			// 用户针对msg进行消息处理
			new_user.DoMessage(msg)

			// 用户的任意消息，代表当前用户是一个活跃的
			isLive <- true
		}
	}()

	// 当前goroutine阻塞
	for {
		select {
		case <-isLive:
			// 当前用户是活跃的，应该重置定时器
			// 一旦本case被触发，以下的case表达式都会执行，但不会进入冒号下的代码

		case <-time.After(time.Second * 100):
			// time.After本质是一个管道，超时后会往管道写一个数据
			new_user.SendMsg("你被踢了")
			close(new_user.C)
			conn.Close()
			return // 或rutime.Goexit()
		}
	}
}

// 启动Server的接口
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println()
		return
	}
	defer listener.Close()

	// 启动监听Messager的goroutine
	go this.ListenMessager()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener acceppt err:", err)
			continue
		}

		// do handler
		go this.Handler(conn)
	}
}
