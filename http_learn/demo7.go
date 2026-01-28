package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ==========================================
// 1. Context & Engine (基础设施 - 已完成)
// ==========================================

// Context 封装 (Day 2 内容，已写好不用动)
type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	Params     map[string]string // 这里的 Params 需要你稍后在 Router 里填充
	StatusCode int
}

func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w, Req: req, Path: req.URL.Path, Method: req.Method,
		Params: make(map[string]string),
	}
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)
	json.NewEncoder(c.Writer).Encode(obj)
}

// Engine (框架入口，已写好不用动)
type HandlerFunc func(*Context)

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.router.addRoute("GET", pattern, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r)
	engine.router.handle(c) // 把请求交给 Router 处理
}

// ==========================================
// 2. Trie Tree & Router (你的战场 !!!)
// ==========================================

// 树节点定义
type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang (只有终点节点才存这个，中间节点为空)
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点
	isWild   bool    // 是否模糊匹配，part 含有 : 或 * 时为 true
}

// Router 定义
type router struct {
	roots    map[string]*node       // 每种 Method 对应一棵树 (GET, POST)
	handlers map[string]HandlerFunc // 存储 pattern -> handler 的映射
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 辅助工具：解析路径 (已提供)
// 输入: "/hello//duke" -> 输出: ["hello", "duke"]
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// -----------------------------------------------------------
// 任务 A: 实现 Trie 树的插入逻辑 (递归)
// -----------------------------------------------------------
func (n *node) insert(pattern string, parts []string, height int) {
	// TODO: 1. 检查是否到达层级终点 (len(parts) == height)
	//          如果是，设置 n.pattern 并返回

	// TODO: 2. 获取当前的 part (parts[height])

	// TODO: 3. 检查 n.children 中是否已经存在匹配该 part 的子节点
	//          (你需要自己实现一个辅助函数 matchChild 或者在这里直接写循环)

	// TODO: 4. 如果没找到，创建一个新节点并加入 children
	//          注意设置 isWild (如果 part 开头是 : 或 *)

	// TODO: 5. 递归调用子节点的 insert
}

// -----------------------------------------------------------
// 任务 B: 实现 Trie 树的搜索逻辑 (递归)
// -----------------------------------------------------------
func (n *node) search(parts []string, height int) *node {
	// TODO: 1. 检查终止条件
	//          (匹配完所有 parts) 或者 (当前节点是通配符 *)
	//          注意：如果匹配完了但 n.pattern 为空，说明这只是个中间路径，不是完整路由，要返回 nil

	// TODO: 2. 获取当前的 part

	// TODO: 3. 找出所有可能的子节点 (matchChildren)
	//          静态匹配：child.part == part
	//          动态匹配：child.isWild == true

	// TODO: 4. 遍历这些候选子节点，递归调用 search
	//          只要有一个返回了非 nil 结果，就立即返回该结果

	return nil
}

// -----------------------------------------------------------
// 任务 C: 实现 Router 的注册与查找
// -----------------------------------------------------------

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern

	// 确保 root 存在
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}

	// TODO: 调用 root.insert
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	// 1. 解析请求路径
	// searchParts := parsePattern(c.Path)

	// 2. 获取对应 Method 的 root 节点

	// TODO: 3. 调用 root.search 查找节点
	// node := ...

	// TODO: 4. 如果 node 不为空：
	//    a. 解析参数 (Param)
	//       你需要遍历 node.pattern 的 parts 和 searchParts
	//       如果是 :name，则 c.Params["name"] = ...
	//       如果是 *filepath，则 c.Params["filepath"] = ... (注意 * 可能包含多个 /)
	//    b. 从 r.handlers 中找到 handler 并执行 c.handlers[key](c)

	// 5. 如果 node 为空，返回 404
	// c.String(404, "404 NOT FOUND: %s\n", c.Path)
}

// ==========================================
// 3. 测试用例 (如果不报错且输出正确，你就成功了)
// ==========================================
func main() {
	r := New()

	// 1. 基础路由
	r.GET("/", func(c *Context) {
		c.String(200, "URL.Path = %q\n", c.Path)
	})

	// 2. 动态路由
	r.GET("/hello/:name", func(c *Context) {
		c.String(200, "hello %s, you are at %s\n", c.Param("name"), c.Path)
	})

	// 3. 通配符路由
	r.GET("/assets/*filepath", func(c *Context) {
		c.JSON(200, map[string]string{"filepath": c.Param("filepath")})
	})

	fmt.Println("Starting server at :9999")
	// 启动服务器
	http.ListenAndServe(":9999", r)
}
