package main

/*
================================================================================
第一步：理解 Go 的 net/http 标准库
================================================================================

什么是 HTTP 服务器？
- 监听某个端口（如9999）
- 接收客户端的HTTP请求（浏览器、curl等）
- 处理请求并返回响应

Go 的 net/http 提供了两个核心概念：
1. http.Handler 接口 - 定义了如何处理HTTP请求
2. http.ListenAndServe - 启动HTTP服务器
*/

import (
	"fmt"
	"log"
	"net/http"
)

// ============================================================================
// 方式1：使用函数作为处理器（最简单）
// ============================================================================

func handler1(w http.ResponseWriter, r *http.Request) {
	// w: ResponseWriter - 用来写响应内容
	// r: Request - 包含请求的所有信息（URL、方法、头部、参数等）

	fmt.Fprintf(w, "Hello! 你访问的路径是: %s\n", r.URL.Path)
	fmt.Fprintf(w, "请求方法: %s\n", r.Method)
}

// ============================================================================
// 方式2：实现 http.Handler 接口（推荐，这是框架的基础）
// ============================================================================

// http.Handler 接口定义：
// type Handler interface {
//     ServeHTTP(ResponseWriter, *Request)
// }
// 任何类型只要实现了 ServeHTTP 方法，就是一个 Handler

type MyHandler struct {
	message string
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", h.message)
	fmt.Fprintf(w, "你访问的路径是: %s\n", r.URL.Path)
}

// ============================================================================
// 方式3：使用 http.HandlerFunc 适配器（函数转Handler）
// ============================================================================

func handler3(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "这是通过 HandlerFunc 包装的处理器\n")
	fmt.Fprintf(w, "路径: %s, 方法: %s\n", r.URL.Path, r.Method)
}

func main() {
	// 选择你想测试的方式（取消注释对应的代码）

	// ========== 测试方式1：函数处理器 ==========
	// http.HandleFunc("/", handler1)
	// log.Println("方式1启动: http://localhost:8001")
	// log.Fatal(http.ListenAndServe(":8001", nil))

	// ========== 测试方式2：自定义Handler ==========
	// handler := &MyHandler{message: "欢迎使用自定义Handler!"}
	// log.Println("方式2启动: http://localhost:8002")
	// log.Fatal(http.ListenAndServe(":8002", handler))

	// ========== 测试方式3：HandlerFunc包装 ==========
	log.Println("方式3启动: http://localhost:8003")
	log.Fatal(http.ListenAndServe(":8003", http.HandlerFunc(handler3)))

	/*
		测试命令：
		curl http://localhost:8003/
		curl http://localhost:8003/hello
		curl http://localhost:8003/any/path

		注意：方式2中，无论访问什么路径，都会被同一个handler处理
		这就是为什么我们需要"路由"！
	*/
}
