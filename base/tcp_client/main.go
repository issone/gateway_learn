package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	doSend()
	fmt.Print("doSend over")
	//doSend()
	//fmt.Print("doSend over")
	//select {}
}

func doSend() {
	//1、连接服务器
	conn, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}
	defer conn.Close()
	//2、读取命令行输入
	inputReader := bufio.NewReader(os.Stdin)
	//var buf [128]byte
	for {
		// 3、一直读取直到读到\n
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Printf("read from console failed, err: %v\n", err)
			break
		}
		// 4、读取Q时停止
		trimmedInput := strings.TrimSpace(input)
		if trimmedInput == "Q" {
			break
		}
		// 5、给服务器发信息
		_, err = conn.Write([]byte(trimmedInput))
		if err != nil {
			fmt.Printf("write failed , err : %v\n", err)
			break
		}

	}
}
