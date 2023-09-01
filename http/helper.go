package coreHttp

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func (this *GCore) newServeMux(list []Controller) *gin.Engine {
	for _, route := range list {
		group := this.gin.Group(route.Path())
		v, ok := this.middleware[route.Path()]
		if ok {
			group.Use(v...)

		}
		route.Build(group)
	}
	return this.gin
}

func asRoute(f []any) []any {
	res := make([]any, 0)
	for _, v := range f {
		res = append(res, fx.Annotate(
			v,
			fx.As(new(Controller)),
			fx.ResultTags(`group:"routes"`),
		))
	}
	return res
}
func (this *GCore) getController() []any {
	fn := make([]any, 0)
	for _, item := range this.Controller {
		fn = append(fn, item)
	}
	return fn
}
