package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// NewUser 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   "name_" + userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	// 单独起一个协程，来监听用户的channel
	go user.ListenMessage()
	return user
}

// ListenMessage 方法用于监听C 通道并发送用户消息
// 这是一个持续运行的方法，会不断从通道中获取消息并发送给用户连接
func (this *User) ListenMessage() {
	// 无限循环，持续监听消息通道
	for {
		// 从通道中获取消息
		// this.C 是一个通道，用于接收需要发送给用户的消息
		msg := <-this.C

		// 将消息加上换行符后发送给用户的连接
		// this.conn 表示与用户的网络连接
		// []byte(msg + "\n") 将消息转换为字节数组并添加换行符
		_, err := this.conn.Write([]byte(msg + "\n"))
		// 如果发送过程中出现错误，则直接返回结束方法
		// 这通常意味着连接已经断开或出现其他网络问题
		if err != nil {
			return
		}
	}
}

// 用户上线
func (user *User) online() {
	// 用户上线，加入onlineMap
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()
	//广播当前消息，给其他人
	user.server.BroadCast(user, "用户已上线")
}

// 用户下线
func (user *User) offline() {
	// 用户下线，删除onlineMap
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()
	//广播当前消息，给其他人
	user.server.BroadCast(user, "用户已下线")
}

// 用户处理消息
func (this *User) doMessage(msg string) {
	if msg == "who" {
		this.server.mapLock.Lock()
		for name, _ := range this.server.OnlineMap {
			onlineMsg := "用户：" + name + " 在线" + "\n"
			this.sendMessage(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 消息格式"rename|张三"
		newName := strings.Split(msg, "|")[1]

		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.sendMessage("当前用户名被使用")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.sendMessage("您已经成功修改用户名为：" + newName + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 消息格式"to|张三|msg"
		peopleName := strings.Split(msg, "|")[1]
		if peopleName == "" {
			this.sendMessage("您输入的格式有误，请使用 \"to|张三|你好啊\"，格式 \n")
			return
		}
		people, ok := this.server.OnlineMap[peopleName]
		if !ok {
			this.sendMessage("该用户名不存在 \n")
			return
		}
		p_msg := strings.Split(msg, "|")[2]
		if p_msg == "" {
			this.sendMessage("您输入的格式有误，请使用 \"to|张三|你好啊\"，格式 \n")
			return
		}
		people.sendMessage(this.Name + "对您说：" + p_msg + "\n")
	} else {
		this.server.BroadCast(this, msg)
	}
}

// 给当前用户对应的客户端发消息
func (this *User) sendMessage(msg string) {
	_, err := this.conn.Write([]byte(msg))
	if err != nil {
		return
	}
}
