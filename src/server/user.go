package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	user_addr := conn.RemoteAddr().String()

	new_user := &User{
		Name:   user_addr,
		Addr:   user_addr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go new_user.ListenMessager()

	return new_user
}

func (this *User) Online() {
	// 用户上线，将用户加入到OnlineMap表中
	this.server.map_lock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.map_lock.Unlock()

	// 广播用户上线消息
	this.server.BroadCast(this, "已上线")
}

func (this *User) Offline() {
	// 用户下线，将用户从OnlineMap表中删除
	this.server.map_lock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.map_lock.Unlock()

	// 广播用户下线消息
	this.server.BroadCast(this, "拜拜嘞，下线咯")

}

// 给当前User对应客户端发消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前用户都有哪些
		this.server.map_lock.Lock()
		for _, user := range this.server.OnlineMap {
			online_msg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			this.SendMsg(online_msg)
		}
		this.server.map_lock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 消息格式：`rename|张三`
		msg_array := strings.Split(msg, "|")
		new_name := msg_array[1]

		// 判断name是否存在
		_, ok := this.server.OnlineMap[new_name]
		if ok {
			this.SendMsg("当前用户名被使用！\n")
		} else {
			this.server.map_lock.Lock() //------> LOCK
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[new_name] = this
			this.server.map_lock.Unlock() // ------> Unlock

			this.Name = new_name
			this.SendMsg("您已成功更新用户名：" + this.Name + "\n")
		}
	} else if len(msg) == 4 && msg == "exit" {
		this.conn.Close()
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 消息格式： `to|Bob|message`

		// 1、获取对方用户名
		msg_array := strings.Split(msg, "|")
		remote_user_name := msg_array[1]
		if remote_user_name == "" {
			this.SendMsg("消息格式不正确，请使用`to|Alice|How you doing?`\n")
			return
		}

		// 2、根据用户名得到对方User对象
		remote_user, ok := this.server.OnlineMap[remote_user_name]
		if !ok {
			this.SendMsg("该用户名不存在\n")
			return
		}

		// 3、根据消息内容，通过对方的User对象将消息内容发送过去
		content := msg_array[2]
		if content == "" {
			this.SendMsg("消息格式不正确，请使用`to|Alice|How you doing?`\n")
			return
		}
		remote_user.SendMsg(this.Name + "对您说：" + content)

	} else {
		this.server.BroadCast(this, msg)
	}
}

// 监听当前User channel的方法，一旦有消息就发送给对端客户
func (this *User) ListenMessager() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
