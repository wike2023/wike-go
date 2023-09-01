package coreHttp

import (
	"github.com/gin-gonic/gin"
	"github.com/wike2023/wike-go/config"
	zaplog "github.com/wike2023/wike-go/lib/log"
	"go.uber.org/fx"
	"net/http"
	"reflect"
)

var Module = fx.Module("infra",
	fx.Provide(zaplog.LoggerInit, config.Config),
	fx.Invoke(func(*http.Server) {}),
)

func (this *GCore) Config(cfgs ...interface{}) *GCore {
	for _, cfg := range cfgs {
		t := reflect.TypeOf(cfg)
		if t.Kind() != reflect.Ptr {
			panic("required ptr object") //必须是指针对象
		}
		if t.Elem().Kind() != reflect.Struct {
			continue
		} //处理依赖注入 (new)
		v := reflect.ValueOf(cfg)
		for i := 0; i < t.NumMethod(); i++ {
			method := v.Method(i)
			callRet := method.Call(nil)

			if callRet != nil && len(callRet) == 1 {
				this.supply = append(this.supply, callRet[0].Interface())
			}
		}
	}

	return this
}

func (this *GCore) Provide(list ...interface{}) *GCore {
	this.provides = append(this.provides, list...)
	return this
}
func (this *GCore) Invokes(list ...interface{}) *GCore {
	this.invokes = append(this.invokes, list...)
	return this
}

func (this *GCore) GlobalUse(middleware ...gin.HandlerFunc) *GCore {
	this.gin.Use(middleware...)
	return this
}
func (this *GCore) Use(path string, middleware ...gin.HandlerFunc) *GCore {
	this.middleware[path] = append(this.middleware[path], middleware...)
	return this
}
func (this *GCore) Mount(classList ...interface{}) *GCore {
	for _, item := range classList {
		this.Controller = append(this.Controller, item)
	}
	return this
}
func (this *GCore) Supply(supply ...interface{}) *GCore {
	this.supply = append(this.supply, supply...)
	return this
}
