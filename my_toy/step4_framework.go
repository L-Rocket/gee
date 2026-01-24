package main

/*
================================================================================
第四步：构建完整的Web框架
================================================================================

现在我们有了：
1. ✅ 路由功能（Router）
2. ✅ 上下文封装（Context）

接下来：将它们组合成一个完整的框架

Web框架的核心组成部分：
┌─────────────────────────────────────┐
│         Engine (引擎)                │  ← 框架的入口
│  - 管理路由器                         │
│  - 提供 GET/POST 等注册方法            │
│  - 实现 http.Handler 接口             │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│         Router (路由器)              │  ← 路由管理
│  - 存储路由映射表                     │
│  - 匹配请求到处理函数                 │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│         Context (上下文)             │  ← 请求/响应封装
│  - 便捷的参数获取                     │
│  - 便捷的响应方法                     │
└─────────────────────────────────────┘

这就是你的 Gee 框架的架构！
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ============================================================================
// Context - 上下文（和step3一样，这里再写一遍方便理解）
// ============================================================================

type Ctx struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	Path       string
	Method     string
	StatusCode int
}

func newCtx(w http.ResponseWriter, r *http.Request) *Ctx {
	return &Ctx{
		Writer:  w,
		Request: r,
		Path:    r.URL.Path,
		Method:  r.Method,
	}
}

func (c *Ctx) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Ctx) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Ctx) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Ctx) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Ctx) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain; charset=utf-8")
	c.Status(code)
	fmt.Fprintf(c.Writer, format, values...)
}

func (c *Ctx) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Ctx) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html; charset=utf-8")
	c.Status(code)
	fmt.Fprint(c.Writer, html)
}

// ============================================================================
// HandlerFunc - 处理函数类型
// ============================================================================

type Handler func(*Ctx)

// ============================================================================
// Router - 路由器
// ============================================================================

type FrameworkRouter struct {
	handlers map[string]Handler
}

func newFrameworkRouter() *FrameworkRouter {
	return &FrameworkRouter{
		handlers: make(map[string]Handler),
	}
}

func (r *FrameworkRouter) addRoute(method, path string, handler Handler) {
	key := method + "-" + path
	r.handlers[key] = handler
	log.Printf("Route %4s - %s", method, path)
}

func (r *FrameworkRouter) handle(c *Ctx) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

// ============================================================================
// Engine - 框架引擎（这是对外的主要接口）
// ============================================================================

type Engine struct {
	router *FrameworkRouter
}

// New 创建一个新的框架实例
func New() *Engine {
	return &Engine{
		router: newFrameworkRouter(),
	}
}

// GET 注册GET路由
func (e *Engine) GET(path string, handler Handler) {
	e.router.addRoute("GET", path, handler)
}

// POST 注册POST路由
func (e *Engine) POST(path string, handler Handler) {
	e.router.addRoute("POST", path, handler)
}

// Run 启动HTTP服务器
func (e *Engine) Run(addr string) error {
	log.Printf("服务器启动在 http://localhost%s", addr)
	return http.ListenAndServe(addr, e)
}

// ServeHTTP 实现 http.Handler 接口
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 为每个请求创建一个Context
	c := newCtx(w, req)
	// 交给路由器处理
	e.router.handle(c)
}

// ============================================================================
// 使用框架 - 看看多么优雅！
// ============================================================================

func main() {
	// 1. 创建框架实例
	r := New()

	// 2. 注册路由
	r.GET("/", func(c *Ctx) {
		c.HTML(http.StatusOK, `
			<h1>🎉 欢迎使用我的Web框架！</h1>
			<h2>这个框架包含：</h2>
			<ul>
				<li>✅ Engine - 框架引擎</li>
				<li>✅ Router - 路由管理</li>
				<li>✅ Context - 上下文封装</li>
			</ul>
			<h2>测试路由：</h2>
			<ul>
				<li><a href="/hello?name=张三">GET /hello</a></li>
				<li><a href="/api/user">GET /api/user</a></li>
				<li>POST /login (需要curl测试)</li>
			</ul>
		`)
	})

	r.GET("/hello", func(c *Ctx) {
		name := c.Query("name")
		if name == "" {
			name = "World"
		}
		c.String(http.StatusOK, "Hello, %s!\n", name)
	})

	r.GET("/api/user", func(c *Ctx) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"username": "张三",
			"age":      25,
			"email":    "zhangsan@example.com",
		})
	})

	r.POST("/login", func(c *Ctx) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		c.JSON(http.StatusOK, map[string]interface{}{
			"username": username,
			"password": password,
			"status":   "success",
			"message":  "登录成功",
		})
	})

	// 3. 启动服务器
	log.Println("测试命令:")
	log.Println("  curl http://localhost:8006/")
	log.Println("  curl http://localhost:8006/hello?name=李四")
	log.Println("  curl http://localhost:8006/api/user")
	log.Println("  curl -X POST http://localhost:8006/login -d 'username=admin&password=123456'")

	r.Run(":8006")
}

/*
================================================================================
🎓 知识总结：Web框架的本质
================================================================================

1. Web框架是什么？
   - 在 net/http 标准库之上的封装
   - 提供更便捷的API，简化Web开发
   - 核心功能：路由、上下文、中间件

2. Web框架由哪些部分构成？

   ┌──────────────────────────────────────────────┐
   │  Engine (引擎)                                │
   │  - 框架的入口和管理者                          │
   │  - 实现 http.Handler 接口                     │
   │  - 提供 GET/POST/Run 等方法                   │
   └──────────────────────────────────────────────┘
                    ↓ 管理
   ┌──────────────────────────────────────────────┐
   │  Router (路由器)                              │
   │  - 维护 路径→处理函数 的映射表                  │
   │  - 负责匹配请求到对应的处理函数                 │
   └──────────────────────────────────────────────┘
                    ↓ 创建
   ┌──────────────────────────────────────────────┐
   │  Context (上下文)                             │
   │  - 封装 Request 和 Response                   │
   │  - 提供便捷的参数获取和响应方法                 │
   └──────────────────────────────────────────────┘
                    ↓ 传递给
   ┌──────────────────────────────────────────────┐
   │  HandlerFunc (处理函数)                       │
   │  - 业务逻辑                                   │
   │  - func(c *Context) { ... }                  │
   └──────────────────────────────────────────────┘

3. HTTP处理流程：

   客户端请求
        ↓
   http.ListenAndServe (Go标准库)
        ↓
   Engine.ServeHTTP (框架入口)
        ↓
   创建 Context (封装请求和响应)
        ↓
   Router.handle (查找路由)
        ↓
   HandlerFunc (执行业务逻辑)
        ↓
   响应返回给客户端

4. 和你的项目对比：

   你的 day2-http-context 项目：
   - ✅ gee.Engine  → 就是这里的 Engine
   - ✅ gee.router  → 就是这里的 Router
   - ✅ gee.Context → 就是这里的 Context

   完全一样的架构！

5. 下一步可以学习：
   - 🚀 动态路由（/user/:id）
   - 🚀 中间件（日志、认证、恢复）
   - 🚀 路由分组
   - 🚀 模板渲染
   - 🚀 静态文件服务

现在你已经理解了Web框架的本质！
你的 Gee 框架就是按照这个思路构建的。

从 step1 → step2 → step3 → step4，你看到了框架是如何一步步演进的：
1. 从最基础的 HTTP 处理开始
2. 添加路由功能
3. 引入Context简化代码
4. 封装成完整的框架

这就是所有Web框架的核心原理！
================================================================================
*/
