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
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServeIp:    serverIp,
		ServerPort: serverPort,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		return nil
	}

	client.conn = conn
	return client
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
	time.Sleep(time.Second * 100000)
}
