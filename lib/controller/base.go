package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type errorCode int

const (
	Success      errorCode = 200
	Failed       errorCode = 500
	ParamError   errorCode = 400
	NotFound     errorCode = 404
	UnAuthorized errorCode = 401
)

var codeMsg = map[errorCode]string{
	Success:      "正常",
	Failed:       "系统异常",
	ParamError:   "参数错误",
	NotFound:     "记录不存在",
	UnAuthorized: "未授权",
}

type Controller struct {
	Page     int
	PageSize int
}
type Common interface {
	Controller
}
type Pagination struct {
	List       interface{} `json:"list"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalCount int64       `json:"total_count"`
}

var BaseController = &Controller{}

func (r *Controller) ParsePage(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 1000 {
		pageSize = 1000
	}
	r.Page = page
	r.PageSize = pageSize
}

func (*Controller) Success(ctx *gin.Context, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":     Success,
		"msg":      msg,
		"data":     data,
		"trace_id": ctx.GetString("trace_id"),
	})
}
func (r *Controller) List(ctx *gin.Context, list interface{}, totalCount int64) {
	if list == nil {
		list = make([]interface{}, 0)
	}
	r.Success(ctx, "获取列表成功", Pagination{
		List:       list,
		Page:       r.Page,
		PageSize:   r.PageSize,
		TotalCount: totalCount,
	})
}

func (*Controller) Failed(ctx *gin.Context, code errorCode, msg string) {
	errMsg := codeMsg[code] + ": " + msg
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":     code,
		"msg":      errMsg,
		"data":     nil,
		"trace_id": ctx.GetString("trace_id"),
	})
	if code != Success {
		ctx.Set("error_code", int(code))
		ctx.Set("error_msg", msg)
	}
}

type StatusError struct {
	Code int
	Msg  string
}

func (this *StatusError) Error() string {
	return this.Msg
}

func Error(code int, msg string) {
	panic(&StatusError{code, msg})
}
