package coreHttp

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	zaplog "github.com/wike2023/wike-go/lib/log"
	"go.uber.org/fx"
	"net"
	"net/http"
)

func (this *GCore) NewHTTPServer(lc fx.Lifecycle, zap *zaplog.Logger, r *gin.Engine, cfg *viper.Viper) *http.Server {
	this.gin = r
	srv := &http.Server{
		Addr:    ":" + cfg.GetString("port"),
		Handler: this.gin,
	}
	this.gin.GET("/healthz", func(c *gin.Context) {
		c.String(200, "ok")
	})
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				zap.Errorln(err.Error())
				return err
			}
			zap.Debugf("Starting HTTP server at %s", srv.Addr)
			go func() {
				if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
					zap.Errorf("HTTP server listen: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
