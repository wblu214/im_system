package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServeIp    string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServeIp:    serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		return nil
	}

	client.conn = conn
	return client
}

func (client *Client) menu() bool {
	var sflag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	_, err := fmt.Scanln(&sflag)
	if err != nil {
		return false
	}
	if sflag >= 0 && sflag <= 3 {
		client.flag = sflag
		return true
	} else {
		fmt.Println("请输入合法的数字")
		return false
	}
}

func (client *Client) updateName() bool {
	fmt.Println("请输入用户名:")
	_, err := fmt.Scanln(&client.Name)
	if err != nil {
		return false
	}

	sendMsg := "rename|" + client.Name + "\n"
	_, w_err := client.conn.Write([]byte(sendMsg))
	if w_err != nil {
		fmt.Println("conn.Write err: ", w_err)
		return false
	}
	return true

}
func (client *Client) DealResponse() {
	//永久阻塞监听网络，conn有数据就copy到stdout标准输出上
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) PublicChat() {
	var chatStr string
	fmt.Println(">>>请输入聊天内容...")
	_, err := fmt.Scanln(&chatStr)
	if err != nil {
		return
	}
	for chatStr != "exit" {
		//发送服务器

		//消息不为空则发送
		if len(chatStr) != 0 {
			sendMsg := chatStr + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err: ", err)
				break
			}
		}
		chatStr = ""
		fmt.Println(">>>请输入聊天内容...")
		_, err := fmt.Scanln(&chatStr)
		if err != nil {
			return
		}
	}
}

func (client *Client) PrivateChat() {
	client.queryUsers()
	var peopleName string
	var chatStr string

	fmt.Println(">>>请输入聊天用户名，exit退出!!")
	fmt.Scanln(&peopleName)

	for peopleName != "exit" {
		fmt.Println(">>>请输入聊天内容，exit退出!!")
		fmt.Scanln(&chatStr)
		//消息不为空则发送
		for chatStr != "exit" {
			if len(chatStr) != 0 {
				sendMsg := "to|" + peopleName + "|" + chatStr + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write err: ", err)
					break
				}
			}
			chatStr = ""
			fmt.Println(">>>请输入聊天内容，exit退出!!")
			_, err := fmt.Scanln(&chatStr)
			if err != nil {
				return
			}
		}

		peopleName = ""
		client.queryUsers()
		fmt.Println(">>>请输入聊天用户名，exit退出!!")
		fmt.Scanln(&peopleName)
	}

}
func (client *Client) queryUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return
	}
}
func (client *Client) Run() {
	for {
		for client.menu() != true {
		}
		//处理不同模式
		switch client.flag {
		case 0:
			//退出
			fmt.Println("正在退出...")
			return
		case 1:
			//公聊模式
			fmt.Println("公聊模式...")
			client.PublicChat()
			break
		case 2:
			//私聊模式
			fmt.Println("私聊模式...")
			client.PrivateChat()
			break
		case 3:
			//更新用户名
			fmt.Println("更新用户名...")
			client.updateName()
			break
		}
	}
}

var serverIp string
var serverPort int

// go run client.go -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "默认为本地127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "默认为8888")

}

func main() {
	// 命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>连接服务器失败！！！")
		return
	}
	//单独启动一个协程去处理server的回执消息
	go client.DealResponse()

	fmt.Println(">>>连接服务器成功！！！")
	client.Run()

	// 用户选择退出后，关闭连接
	if client.conn != nil {
		client.conn.Close()
	}
	fmt.Println(">>>已退出聊天系统！！！")
}
