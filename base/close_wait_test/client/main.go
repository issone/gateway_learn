package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main()  {
	doSend()
	fmt.Print("doSend over")
}


//  netstat -Aaln |grep 9090 可以用这个命令查看连接状态

func doSend() {
	//1、连接服务器
	conn, err := net.Dial("tcp", "localhost:9090")
	defer conn.Close()	//这里不填写会无法释放连接，使其处于连接状态（ESTABLISHED），而不是CLOSE状态
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}
	//2、读取命令行输入
	inputReader := bufio.NewReader(os.Stdin)
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
		// 5、回复服务器信息
		_, err = conn.Write([]byte(trimmedInput))
		if err != nil {
			fmt.Printf("write failed , err : %v\n", err)
			break
		}
	}

	// 测试客户端不执行conn.Close时，取消下面的注释，否则看不到效果，因为程序执行完，进程关闭，进程里的资源也
	//fmt.Println("client close")
	//time.Sleep(100* time.Second)
}
