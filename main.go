package main

import (
	coreHttp "github.com/wike2023/wike-go/http"
)

func main() {
	coreHttp.God().
		Mount(coreHttp.RouterInit, coreHttp.RouterInit2).
		Provide(coreHttp.InitSvc).
		Config(&coreHttp.Conf{}).
		GlobalUse(coreHttp.RequestDurationMiddleware).
		GlobalUse(coreHttp.RequestDurationMiddleware2).
		Use("/api", coreHttp.HeaderMiddleware, coreHttp.CookieMiddleware).
		Use("/api/v2", coreHttp.HeaderMiddleware).
		Run()
}
