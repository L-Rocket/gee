package main

import (
	"fmt"
	"net/http"
)

// 这是具体的业务逻辑（厨师）
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!") // 往 w 里写，就是发回给浏览器
}

func noHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "No, World!") // 往 w 里写，就是发回给浏览器
}

func main() {
	// 1. 注册路由（写菜单）
	// 此时只是把 "/hello" -> helloHandler 的关系存入 DefaultServeMux
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/No", noHandler)

	// 2. 启动服务（开门营业）
	// 监听 8080 端口。
	// 第二个参数是 nil，说明使用上面注册好的 DefaultServeMux
	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
