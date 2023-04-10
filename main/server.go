package main

import (
	"fmt"
	"net"
)

//初始化端口
func runServer(addr string)(err error){
	listener, err := net.Listen("tcp",addr)
	if err != nil {
		fmt.Println("listen failed,", err)
		return
	}

	for  {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept failed,", err)
			continue
		}

		go process(conn)
	}
}

// 创建Client,处理请求
func process(conn net.Conn){
	client := &Client{
		conn: conn,
	}
	defer conn.Close()
	err := client.Process()
	if err != nil{
		fmt.Println("client process failed,",err)
		return
	}
}
