package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rs1 := &RealServer{Addr: "127.0.0.1:2003"}
	rs1.Run()
	rs2 := &RealServer{Addr: "127.0.0.1:2004"}
	rs2.Run()

	/*
		SIGINT信号为ctrl+c发送出来的信号，该信号只能结束前台进程，信号被当前进程树接收到，也就是说，不仅当前进程会收到信号，它的子进程也会收到；

		SIGTERM可以被阻塞，kill不加任何参数的时候发送的就是这个信号，只有当前进程收到信号，子进程不会收到。如果当前进程被kill了，那么它的子进程的父进程将会是init，也就是pid为1的进程；

		SIGKILL是不能被捕获的，程序收到这个信号后，一定会退出。这就是kill -9一定能保证将程序杀死的原因。
	*/

	//监听关闭信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

type RealServer struct {
	Addr string
}

func (r *RealServer) Run() {
	log.Println("Starting httpserver at " + r.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.HelloHandler)
	mux.HandleFunc("/base/error", r.ErrorHandler)
	mux.HandleFunc("/test_http_string/test_http_string/aaa", r.TimeoutHandler)
	server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}

func (r *RealServer) HelloHandler(w http.ResponseWriter, req *http.Request) {
	//127.0.0.1:8008/abc?sdsdsa=11
	//r.Addr=127.0.0.1:8008
	//req.URL.Path=/abc
	//fmt.Println(req.Host)	// 客户端地址 IP:PORT
	upath := fmt.Sprintf("http://%s%s\n", r.Addr, req.URL.Path)
	realIP := fmt.Sprintf("RemoteAddr=%s,X-Forwarded-For=%v,X-Real-Ip=%v\n", req.RemoteAddr, req.Header.Get("X-Forwarded-For"), req.Header.Get("X-Real-Ip"))
	header := fmt.Sprintf("headers =%v\n", req.Header)
	io.WriteString(w, upath)
	io.WriteString(w, realIP)
	io.WriteString(w, header)

}

func (r *RealServer) ErrorHandler(w http.ResponseWriter, req *http.Request) {
	upath := "error handler"

	w.WriteHeader(500)

	io.WriteString(w, upath)
}

func (r *RealServer) TimeoutHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(6 * time.Second)
	upath := "timeout handler"
	w.WriteHeader(200)
	io.WriteString(w, upath)
}
