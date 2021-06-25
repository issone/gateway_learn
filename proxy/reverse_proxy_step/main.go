package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

var addr = "127.0.0.1:2002"

func main() {
	//127.0.0.1:2002/xxx
	//127.0.0.1:2003/base/xxx
	rs := "http://127.0.0.1:2003/base"
	url1, err1 := url.Parse(rs)
	if err1 != nil {
		log.Println(err1)
	}
	proxy := NewSingleHostReverseProxy(url1)
	log.Println("Starting httpserver at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}

func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	// 新建一个proxy
	// 代理路径rs是http://127.0.0.1:2003/base
	// 请求的路径如果是 http://127.0.0.1:2002/dir
	// 则实际路径为 http://127.0.0.1:2003/base/dir

	//对于路径 http://127.0.0.1:2002/dir?name=123而言，RayQuery: name=123，Scheme: http，Host: 127.0.0.1:2002

	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		// 自定义重写URL规则

		// 默认情况 127.0.0.1:2002/abc ==> 127.0.0.1:2003/base/abc
		// 增加path为dir开头时，替换dir
		//127.0.0.1:2002/dir/abc ==> 127.0.0.1:2003/base/abc

		// http://127.0.0.1:2002/dir123/abc ==> http://127.0.0.1:2003/base/123/abc

		re, _ := regexp.Compile("^/dir(.*)");
		req.URL.Path = re.ReplaceAllString(req.URL.Path, "$1")

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		//target.Path : /base
		//req.URL.Path : /dir
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}

	modifyFunc := func(res *http.Response) error {
		// 修改返回内容
		if res.StatusCode != 200 {
			//return errors.New("error statusCode")

			// oldPayload为原始内容
			oldPayload, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}

			// 在原始内容上增加前缀
			newPayLoad := []byte("hello " + string(oldPayload))

			// 将byte切片转换为ReadCloser对象
			res.Body = ioutil.NopCloser(bytes.NewBuffer(newPayLoad))
			// 修改内容后，要修改ContentLength
			res.ContentLength = int64(len(newPayLoad))
			res.Header.Set("Content-Length", fmt.Sprint(len(newPayLoad)))
		}
		return nil
	}

	errorHandler := func(res http.ResponseWriter, req *http.Request, err error) {
		res.Write([]byte(err.Error()))

	}

	return &httputil.ReverseProxy{Director: director, ModifyResponse: modifyFunc, ErrorHandler: errorHandler}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
