package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	option     int
}

func NewClient(servIp string, servPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   servIp,
		ServerPort: servPort,
		option:     999,
	}

	// 连接Server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", servIp, servPort))
	if err != nil {
		fmt.Println("net.Dial ERROR:", err)
		return nil
	}

	client.conn = conn

	// 返回对象
	return client
}

// goroutine处理server的响应包
func (this *Client) DealResponse() {
	io.Copy(os.Stdout, this.conn)
}

func (this *Client) menu() bool {
	var option int

	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("0. 退出")

	fmt.Scanf("%d", &option)
	if option >= 0 && option <= 3 {
		this.option = option
		return true
	} else {
		fmt.Println(">>>>请输入合法的数字范围<<<<")
		return false
	}
}

func (this *Client) SelectUser() {
	send_msg := "who\n"
	_, err := this.conn.Write([]byte(send_msg))
	if err != nil {
		fmt.Println("ERROR SelectUser():", err)
		return
	}
}

func (this *Client) PrivateChat() {
	this.SelectUser()

	var remote_name string
	fmt.Println(">>>>>>请输入你要对话的用户: (exit退出)")
	fmt.Scanln(&remote_name)

	for remote_name != "exit" {
		var chat_msg string
		fmt.Println(">>>>>>请输入你的聊天内容: (exit退出)")
		fmt.Scanln(&chat_msg)

		for chat_msg != "exit" {
			// 消息不为空，则发送给服务器
			if len(chat_msg) != 0 {
				send_msg := "to|" + remote_name + "|" + chat_msg + "\n\n"
				_, err := this.conn.Write([]byte(send_msg))
				if err != nil {
					fmt.Println("ERROR sending msg:", err)
					break
				}
			}

			chat_msg = ""
			fmt.Println(">>>>输入聊天内容，exit退出.")
			fmt.Scanln(&chat_msg)
		}

		this.SelectUser()
		fmt.Println(">>>>>>请输入你要对话的用户: (exit退出)")
		fmt.Scanln(&remote_name)
	}
}

func (this *Client) PublicChat() {
	// 提示用户输入消息
	var chat_msg string
	fmt.Println(">>>>输入聊天内容，exit退出.")
	fmt.Scanln(&chat_msg) //不知道为什么你不能输入一个带空格的字符串

	for chat_msg != "exit" {
		// 消息不为空，则发送给服务器
		if len(chat_msg) != 0 {
			send_msg := chat_msg + "\n"
			_, err := this.conn.Write([]byte(send_msg))
			if err != nil {
				fmt.Println("ERROR sending msg:", err)
				break
			}
		}

		chat_msg = ""
		fmt.Println(">>>>输入聊天内容，exit退出.")
		fmt.Scanln(&chat_msg)
	}

}

func (this *Client) UpdateName() bool {
	fmt.Print(">>>>请输入用户名: ")
	fmt.Scanln(&this.Name)

	send_msg := "rename|" + this.Name + "\n"
	_, err := this.conn.Write([]byte(send_msg))
	if err != nil {
		fmt.Println("ERROR in sending msg:", err)
		return false
	}

	return true
}

func (this *Client) Run() {
	for {
		for this.menu() != true {
		}

		// 根据option处理不同的业务
		switch this.option {
		case 1:
			// 公聊
			this.PublicChat()
			break
		case 2:
			// 私聊
			this.PrivateChat()
			break
		case 3:
			// 更新用户名
			if this.UpdateName() != true {
				fmt.Println("更新名字失败！")
			}
			break
		case 0:
			this.conn.Write([]byte("exit\n"))
			return
		default:
			// 退出
			fmt.Println("一个奇怪的输入")
			return
		}
	}
}

var ServerIp string
var ServerPort int

func init() {
	// ./client -ip 127.0.0.1 -p 8888
	// 注册flag的一些操作
	flag.StringVar(&ServerIp, "ip", "127.0.0.1", "设置服务器IP地址")
	flag.IntVar(&ServerPort, "p", 8888, "设置服务器端口")
}

func main() {
	// 命令行解析
	flag.Parse()

	client := NewClient(ServerIp, ServerPort)
	if client == nil {
		fmt.Println("xxxxxxxxxx 连接服务器失败...")
		return
	}

	// 单独开一个goroutine去处理server的响应包
	go client.DealResponse()

	fmt.Println(">>>>>>>>>> 连接服务器成功 ...")

	// 启动客户端业务
	client.Run()
}
