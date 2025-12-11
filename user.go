package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// NewUser 创建一个用户的API
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: "name_" + userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
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
