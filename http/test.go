package coreHttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wike2023/wike-go/lib/controller"
	zaplog "github.com/wike2023/wike-go/lib/log"
	"time"
)

type Router struct {
	log *zaplog.Logger
	*controller.Controller
}

func (r *Router) Get(ctx *gin.Context) {
	//	fmt.Println(r.log)
	r.ParsePage(ctx)
	r.log.Info("这个是依赖注入的日志")
	r.List(ctx, []int{1, 2, 3, 4, 5}, 5)
}
func (r *Router) Init() {
	r.Controller = controller.BaseController
}
func (r *Router) Path() string {
	return "/api/v2"
}
func (this *Router) Build(r *gin.RouterGroup) {
	r.GET("/hello", this.Get)
}
func RouterInit(log *zaplog.Logger) *Router {
	r := &Router{log: log}
	r.Init()
	return r
}

type Router2 struct {
	log  *zaplog.Logger
	svc  *Svc
	svc2 *Svc2
	svc3 *Svc3
	*controller.Controller
}

func (r *Router2) Init() {
	r.Controller = controller.BaseController
}
func (r *Router2) Get(ctx *gin.Context) {
	r.svc.String()
	r.log.Info("这个是依赖注入的日志")
	fmt.Println(r.svc2.String())
	fmt.Println(r.svc3.String())
	r.ParsePage(ctx)
	fmt.Println(r.Page)
	fmt.Println(r.PageSize)
	r.Success(ctx, "这个是成功的页", nil)

}

func RouterInit2(log *zaplog.Logger, svc *Svc, svc3 *Svc3, svc2 *Svc2) *Router2 {
	t := &Router2{log: log, svc: svc, svc2: svc2, svc3: svc3}
	t.Init()
	return t
}
func (r *Router2) Path() string {
	return "/api/v2"
}
func (this *Router2) Build(r *gin.RouterGroup) {
	//	this.log.Warnln("111111")
	//this.log.Error("333333")
	r.GET("/ping", this.Get)
	r.GET("/wike2", TestMiddleware, this.Get)
	r.Use(func(ctx *gin.Context) {
		fmt.Println("进入了中间件*****************")
	})
	r.GET("/wike", this.Get)
	r.Use(func(ctx *gin.Context) {
		fmt.Println("第二个*****************")
	})

}

func RequestDurationMiddleware(c *gin.Context) {
	startTime := time.Now()

	c.Next()

	// 执行时间
	duration := time.Since(startTime)
	fmt.Printf("请求处理时间: %v\n", duration)
}
func RequestDurationMiddleware2(c *gin.Context) {
	startTime := time.Now()

	c.Next()

	// 执行时间
	duration := time.Since(startTime)
	fmt.Printf("请求处理时间2: %v\n", duration)
}
func HeaderMiddleware(c *gin.Context) {
	fmt.Println("进入了header")

	c.Next()

	// 执行时间
	fmt.Println("离开了header")
}
func TestMiddleware(c *gin.Context) {
	fmt.Println("TestMiddleware")

	c.Next()

	// 执行时间
	fmt.Println("TestMiddleware")
}
func CookieMiddleware(c *gin.Context) {
	fmt.Println("进入了Cookie")

	c.Next()

	// 执行时间
	fmt.Println("离开了Cookie")
}
func InitSvc() *Svc {
	return &Svc{
		name: "111123123123",
	}
}

type Svc struct {
	name string
}

func (this *Svc) String() string {
	fmt.Println(this.name)
	return this.name
}

type Svc2 struct {
	name string
}

func (this *Svc2) String() string {
	return "svc22222222222"
}

type Svc3 struct {
	name string
}

func (this *Svc3) String() string {
	return "33333333333333"
}

type Conf struct {
}

func (this *Conf) M1() *Svc3 {
	return &Svc3{}
}
func (this *Conf) M2() *Svc2 {
	return &Svc2{}
}
