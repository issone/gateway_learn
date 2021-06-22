package main

import (
	"log"
	"net/http"
	"time"
)

// 关于timeout https://www.jianshu.com/p/4925bfcae05a

var (
	Addr = ":1210"
)

func main() {
	// 创建路由器
	mux := http.NewServeMux()
	// 设置路由规则
	mux.HandleFunc("/bye", sayBye)
	// 创建服务器
	server := &http.Server{
		Addr:         Addr,
		WriteTimeout: time.Second * 3, //请求响应超时时间（从收到请求，并将响应结果返回给客户端的总时间）
		Handler:      mux,
	}
	// 监听端口并提供服务
	log.Println("Starting httpserver at " + Addr)
	log.Fatal(server.ListenAndServe())
}

func sayBye(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	w.Write([]byte("bye bye ,this is httpServer"))
}
