package main

import (
	"fmt"
	"net"
)

func main() {
	//1、监听端口
	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		fmt.Printf("listen fail, err: %v\n", err)
		return
	}
	//2.建立套接字连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept fail, err: %v\n", err)
			continue
		}
		//3. 创建处理协程
		go func(conn net.Conn) {
			defer conn.Close() //这里不填写，服务端出现CLOSE_WAIT状态，客户端会出现FINISH_WAIT_2状态(一段时间后，会由于服务端的探测机制，超过存活时间的会被释放掉)
			for {
				var buf [128]byte
				n, err := conn.Read(buf[:])
				if err != nil {	// 客户端关闭时，会受到EOF
					fmt.Printf("read from connect failed, err: %v\n", err)
					break
				}
				str := string(buf[:n])
				fmt.Printf("receive from client, data: %v\n", str)
			}
		}(conn)
	}
}
