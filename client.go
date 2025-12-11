package main

import (
	"flag"
	"fmt"
	"net"
	"time"
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
func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		//处理不同模式
		switch client.flag {
		case 1:
			//公聊模式
			fmt.Println("公聊模式...")
			break
		case 2:
			//私聊模式
			fmt.Println("私聊模式...")
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
	fmt.Println(">>>连接服务器成功！！！")
	client.Run()
	time.Sleep(time.Second * 100000)
}
