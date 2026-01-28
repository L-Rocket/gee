package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ==========================================
// 1. Trie 树节点定义 (Day 3 核心)
// ==========================================
type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang。只有在终点节点不为空
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为 true
}

// 辅助函数：为了 insert，找到匹配的子节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 辅助函数：为了 search，找到所有匹配的子节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 递归插入节点
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern // 只有终点才赋值 pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)

	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 递归搜索节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

// ==========================================
// 2. Context (Day 2 + Day 3 更新)
// ==========================================
type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	Params     map[string]string // 【Day 3 新增】存放解析出来的参数
	StatusCode int
}

func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		Params: make(map[string]string),
	}
}

// 【Day 3 新增】获取路由参数的方法
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) Header(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) JSON(code int, obj interface{}) {
	c.Header("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.Header("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// ==========================================
// 3. Router (Day 3 核心逻辑)
// ==========================================
type HandlerFunc func(*Context)

type router struct {
	roots    map[string]*node       // 存储每种请求方式的 Trie 树根节点
	handlers map[string]HandlerFunc // 存储 handler
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 解析路径: /hello/world -> ["hello", "world"]
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

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern

	// 如果该方法还没有树，创建一个根节点
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}

	// 插入 Trie 树
	r.roots[method].insert(pattern, parts, 0)
	// 注册 Handler
	r.handlers[key] = handler
}

// 核心查找逻辑
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	// 在 Trie 树中搜索
	n := root.search(searchParts, 0)

	if n != nil {
		// 如果找到了，解析参数
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

// ==========================================
// 4. Engine (入口)
// ==========================================
type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.router.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.router.addRoute("POST", pattern, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r)

	// Day 3 修改：使用 router.getRoute 获取节点和参数
	node, params := engine.router.getRoute(c.Method, c.Path)

	if node != nil {
		c.Params = params // 将解析出的参数放入 Context
		key := c.Method + "-" + node.pattern
		engine.router.handlers[key](c)
	} else {
		c.String(404, "404 NOT FOUND: %s\n", c.Path)
	}
}

// ==========================================
// Main 测试函数
// ==========================================
func main() {
	r := New()

	// 1. 静态路由
	r.GET("/", func(c *Context) {
		c.String(200, "Hello Geektutu")
	})

	// 2. 动态参数路由 (:name)
	// 访问 /hello/geektutu -> 匹配成功，name=geektutu
	r.GET("/hello/:name", func(c *Context) {
		// 使用 c.Param 获取 Trie 树解析出的参数
		c.String(200, "hello %s, you are at %s\n", c.Param("name"), c.Path)
	})

	// 3. 通配符路由 (*filepath)
	// 访问 /assets/css/style.css -> 匹配成功，filepath=css/style.css
	r.GET("/assets/*filepath", func(c *Context) {
		c.JSON(200, map[string]string{
			"filepath": c.Param("filepath"),
		})
	})

	fmt.Println("Gee Server running on :9999")
	http.ListenAndServe(":9999", r)
}
