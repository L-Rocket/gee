package main

import (
	"encoding/json" // 【新增】为了处理 JSON，必须引入这个包
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
	fmt.Printf("[Route Debug] Registering %s\n", key)
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

	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	// ==========================================
	// 【重点体验区域】原生写法实现 JSON 返回
	// ==========================================
	r.GET("/api/student", func(w http.ResponseWriter, req *http.Request) {
		// 1. 准备数据 (这是唯一的业务逻辑)
		data := map[string]interface{}{
			"name":   "Duke Student",
			"age":    22,
			"course": "Computer Architecture",
		}

		// 2. 痛苦点一：手动序列化
		// 必须处理 error，如果这里错了，后面都得崩
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 3. 痛苦点二：手动设置 Header
		// 必须在 WriteHeader 之前设置，否则无效！如果不设，客户端可能无法正确解析。
		w.Header().Set("Content-Type", "application/json")

		// 4. 痛苦点三：手动写状态码
		w.WriteHeader(http.StatusOK)

		// 5. 痛苦点四：写入字节流
		w.Write(jsonData)
	})
	// ==========================================

	fmt.Println("Gee Server running on :9999")
	http.ListenAndServe(":9999", r)
}
