package main

import (
	"bytes"
	"fmt"
	"gateway_learn/base/unpack/unpack"
)

func main() {
	//类比接收缓冲区 net.Conn
	bytesBuffer := bytes.NewBuffer([]byte{})

	//发送
	if err := unpack.Encode(bytesBuffer, "hello world 0!!!"); err != nil {
		panic(err)
	}
	if err := unpack.Encode(bytesBuffer, "hello world 1!!!"); err != nil {
		panic(err)
	}

	//读取
	for {
		if bt, err := unpack.Decode(bytesBuffer); err == nil {
			fmt.Println(string(bt))
			continue
		}
		break
	}
}
