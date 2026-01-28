package main

import (
	"fmt"
	"net/http"
)

// 1. 定义我们自己的核心结构体
// 目前它是空的，以后里面会放 路由表(Router)、中间件(Middleware) 等
type Engine struct{}

// 2. 只有实现了 ServeHTTP 方法，Engine 才能被视为一个 Handler
// 这是一个“万能入口”。所有的请求都会进到这里！
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	case "/hello":
		// 这里我们再次用到了 Header 和 Fprintf
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		// 统一处理 404
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL.Path)
	}
}

func main() {
	// 3. 实例化我们的 Engine
	engine := new(Engine)

	fmt.Println("Gee Server is running on :8080")

	// 4. 【关键一步】篡位！
	// 不再传 nil，而是传入我们的 engine
	// 也就是告诉 http 包：收到请求别自己瞎处理了，全部转交给 engine.ServeHTTP
	http.ListenAndServe(":8080", engine)
}
