package main

import (
	"fmt"
	"net/http"
)

// 定义函数类型，方便后面引用
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine 实现了 http.Handler 接口
type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	fmt.Printf("[Route Debug] Registering %s\n", key) // 打印个日志看看
	engine.router[key] = handler
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// 核心：请求分发器
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, r)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL.Path)
	}
}

func main() {
	r := New()

	// 你看！现在的用法是不是很有“框架感”了？
	// 我们不再用 http.HandleFunc，而是用 r.GET
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	fmt.Println("Gee Server running on :9999")
	http.ListenAndServe(":9999", r)
}
