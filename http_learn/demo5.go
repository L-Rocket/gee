package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ================= Context 开始 =================
type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
}

func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

func (c *Context) Header(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 封装好的 JSON 方法
func (c *Context) JSON(code int, obj interface{}) {
	c.Header("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 封装好的 Query 方法
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// ================= Context 结束 =================

// 修改 HandlerFunc：只接收 *Context
type HandlerFunc func(*Context)

type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// ServeHTTP 中初始化 Context
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r) // 这里创建 Context
	key := c.Method + "-" + c.Path
	if handler, ok := engine.router[key]; ok {
		handler(c) // 传递 Context
	} else {
		c.Header("Content-Type", "text/plain")
		c.Status(404)
		c.Writer.Write([]byte("404 NOT FOUND: " + c.Path))
	}
}

func main() {
	r := New()

	// 体验 1：更简单的 Query 获取
	// 访问：/hello?name=Duke
	r.GET("/hello", func(c *Context) {
		name := c.Query("name") // 直接拿参数，不用解析 r.URL
		if name == "" {
			name = "Guest"
		}
		c.Writer.Write([]byte("Hello " + name))
	})

	// 体验 2：极简的 JSON 返回
	// 访问：/api/student
	r.GET("/api/student", func(c *Context) {
		// H 是 map[string]interface{} 的别名，Gee 框架里常用，这里直接写 map
		c.JSON(200, map[string]interface{}{
			"name":   "Geektutu",
			"age":    20,
			"skills": []string{"Go", "Docker", "Context"},
		})
	})

	fmt.Println("Gee Server running on :9999")
	http.ListenAndServe(":9999", r)
}

// func main() {
// 	r := New()

// 	// 1. 我们尝试注册一个看起来像动态参数的路由
// 	// 我们心中期望：:name 可以匹配任何字符串
// 	r.GET("/hello/:name", func(c *Context) {
// 		c.JSON(200, map[string]string{
// 			"name": "这里应该显示名字",
// 		})
// 	})

// 	// 2. 只有精准匹配才能成功
// 	r.GET("/hello/world", func(c *Context) {
// 		c.Writer.Write([]byte("精确匹配 world 成功"))
// 	})

// 	http.ListenAndServe(":9999", r)
// }
