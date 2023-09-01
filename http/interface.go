package coreHttp

import "github.com/gin-gonic/gin"

type Controller interface {
	Build(r *gin.RouterGroup)
	Path() string
	Get(ctx *gin.Context)
	Init()
}
