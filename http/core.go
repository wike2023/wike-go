package coreHttp

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type GCore struct {
	gin        *gin.Engine
	app        *fx.App
	provides   []interface{}
	supply     []interface{}
	invokes    []interface{}
	Controller []interface{}
	middleware map[string][]gin.HandlerFunc
	port       string
}

func God() *GCore {
	r := gin.New()
	r.Use(AddTrace(), CustomRecover(), AccessLog())
	return &GCore{
		gin:        r,
		Controller: make([]interface{}, 0),
		middleware: map[string][]gin.HandlerFunc{},
		provides:   make([]interface{}, 0),
		invokes:    make([]interface{}, 0),
		supply:     make([]interface{}, 0),
		port:       "8080",
	}
}
func (this *GCore) Run() {

	this.app = fx.New(
		Module,
		fx.NopLogger,
		fx.Provide(this.NewHTTPServer),
		fx.Provide(fx.Annotate(
			this.newServeMux,
			fx.ParamTags(`group:"routes"`),
		)),
		fx.Supply(this.supply...),
		fx.Provide(this.provides...),
		fx.Invoke(this.invokes...),
		fx.Provide(asRoute(this.getController())...),
	)
	this.app.Run()
}
