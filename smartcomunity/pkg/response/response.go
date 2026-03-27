package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 基础响应结构体
type Response struct {
	Code int         `json:"code"` // 业务状态码 (0或200表示成功，其他表示失败)
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 数据载体
}

// const 定义业务状态码 (根据需要扩展)
const (
	CodeSuccess      = 200 // 成功
	CodeServerBusy   = 500 // 系统繁忙
	CodeInvalidParam = 400 // 参数错误
	CodeAuthFail     = 401 // 鉴权失败
)

// Result 基础响应方法
func Result(c *gin.Context, httpStatus int, code int, msg string, data interface{}) {
	c.JSON(httpStatus, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// Success 成功响应 (code默认200)
func Success(c *gin.Context, data interface{}) {
	Result(c, http.StatusOK, CodeSuccess, "success", data)
}

// Fail 失败响应 (code默认400，也可以自定义)
func Fail(c *gin.Context, msg string) {
	Result(c, http.StatusOK, CodeInvalidParam, msg, nil)
}

// FailWithCode 带自定义业务码的失败响应
func FailWithCode(c *gin.Context, code int, msg string) {
	Result(c, http.StatusOK, code, msg, nil)
}
