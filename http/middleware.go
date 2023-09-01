package coreHttp

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wike2023/wike-go/lib/controller"
	zaplog "github.com/wike2023/wike-go/lib/log"
	"go.uber.org/zap"
	"io"
	"time"
)

func AddTrace() gin.HandlerFunc {
	return func(context *gin.Context) {
		traceId := context.Request.Header.Get("trace_id")
		if traceId == "" {
			traceId = uuid.NewString()
		}
		context.Set("trace_id", traceId)
		context.Next()
	}
}

func AccessLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		var reqBody []byte
		if ctx.Request.Body != nil {
			reqBody, _ = io.ReadAll(ctx.Request.Body)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}
		ctx.Next()
		latency := time.Now().Sub(start)
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		statusCode := ctx.Writer.Status()
		bodySize := ctx.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		zaplog.GetLogger().GetLoggerWithGin(ctx).Infow("接口访问日志",
			zap.String("path", path),
			zap.String("method", method),
			zap.String("http_host", ctx.Request.Host),
			zap.String("ua", ctx.Request.UserAgent()),
			zap.String("remote_addr", ctx.Request.RemoteAddr),
			zap.Int("status_code", statusCode),
			zap.Int("error_code", ctx.GetInt("error_code")),
			zap.String("error_msg", ctx.GetString("error_code")),
			zap.Int("body_size", bodySize),
			zap.String("client_ip", clientIP),
			zap.Duration("latency", latency),
		)
	}
}
func CustomRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				check, ok := err.(*controller.StatusError)
				if ok {
					zaplog.GetLogger().GetLoggerWithGin(c).Warnw("主动抛出的错误", zap.String("path", c.Request.URL.Path), zap.String("error", check.Msg), zap.Int("code", check.Code))
					c.JSON(check.Code, gin.H{"message": check.Msg, "code": check.Code})
					return
				}
				zaplog.GetLogger().GetLoggerWithGin(c).Errorw("接口错误", zap.String("path", c.Request.URL.Path), zap.String("error", check.Msg), zap.Int("code", check.Code))
				c.JSON(500, gin.H{"message": "Internal Server Error", "code": 500})
				return
			}
		}()
		c.Next()
	}
}
