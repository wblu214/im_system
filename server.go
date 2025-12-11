package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	IP   string
	Port int
	//在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	//消息广播的channel
	Message chan string
}

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) Handler(conn net.Conn) {
	//fmt.Println("socket连接建立成功！！！")

	user := NewUser(conn, this)
	user.online()

	isLive := make(chan bool)
	// 接收用户发送的消息
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered from panic: %v\n", r)
			}
			user.offline() // 确保用户下线
		}()

		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			// 先处理错误
			if err != nil {
				if err != io.EOF {
					fmt.Println("io read err", err)
				}
				return
			}
			if n == 0 {
				return
			}
			// 安全地处理消息，避免n-1可能导致的越界
			var msg string
			if n > 0 {
				msg = string(buf[:n])
				// 如果需要去除最后一个字符
				if len(msg) > 0 {
					msg = msg[:len(msg)-1]
				}
			}
			user.doMessage(msg)
			isLive <- true
		}
	}()

	// 持续监听用户活动
	for {
		select {
		case <-isLive:
			//不用做任何操作，继续循环
		case <-time.After(time.Second * 100):
			user.sendMessage("由于长时间未操作，你被强制踢出!!!")
			//关闭通道
			close(user.C)
			// 释放链接
			err := conn.Close()
			if err != nil {
				return
			}
			return
		}
	}
}

// BroadCast 方法用于向服务器广播消息
func (this *Server) BroadCast(user *User, msg string) {
	// 将用户地址、用户名和消息内容用"|"符号连接成一个字符串
	sendMsg := user.Addr + "|" + user.Name + "|" + msg
	// 将格式化后的消息发送到服务器的消息通道
	this.Message <- sendMsg
}

// ListenMessage 监听Message广播消息的go routine，一旦有消息就发送给全部在线User
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for name, user := range this.OnlineMap {
			user.C <- msg
			log.Println("当前广播用户为", name)
		}
		this.mapLock.Unlock()
	}
}

// Start 启动服务器的接口
func (this *Server) Start() {
	//监听socket
	listenr, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.IP, this.Port))

	if err != nil {
		fmt.Println("listen this port error:", err)
		return
	}

	// 记得close listen 套接字
	defer func(listenr net.Listener) {
		err := listenr.Close()
		if err != nil {
			fmt.Println("close socket error:", err)
			return
		}
	}(listenr)

	go this.ListenMessage()
	for {
		//accept
		conn, err := listenr.Accept()
		if err != nil {
			fmt.Println("accept socket error:", err)
			continue
		}
		// do handler
		go this.Handler(conn)
	}

	//close listen socket

}
