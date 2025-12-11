package main

import "fmt"

func main() {
	srv := NewServer("127.0.0.1", 8888)
	fmt.Println("启动成功！，正在监听中...")
	srv.Start()
}
