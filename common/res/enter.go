// package res: 统一响应封装
package res

import (
	"github.com/gin-gonic/gin"
)

type Code int

const (
	SuccessCode     Code = 0
	FailValidCode   Code = 1001 // 参数校验错误
	FailAuthCode    Code = 1002 // 权限不足
	FailServiceCode Code = 1003 // 服务错误
	FailTokenCode   Code = 1004 // token不合法
)

func (c Code) String() string {
	switch c {
	case SuccessCode:
		return "成功"
	case FailValidCode:
		return "参数校验错误"
	case FailAuthCode:
		return "权限不足"
	case FailTokenCode:
		return "token不合法"
	case FailServiceCode:
		return "服务异常"
	default:
		return "未知错误"
	}
}

// 空对象，用于避免返回 nil
var Empty = map[string]any{}

// response 统一响应结构
type response struct {
	Code  int    `json:"code"`
	Data  any    `json:"data"`
	Msg   string `json:"msg"`
	LogId string `json:"logId"`
}

// JSON 输出响应
func (r response) JSON(c *gin.Context) {
	// 可选：中断后续中间件执行
	// c.Abort()
	r.LogId = c.GetString("logId")
	c.JSON(200, r)
}

// Success 成功响应（自定义消息 + 数据）
func Success(c *gin.Context, msg string, data any) {
	response{Code: int(SuccessCode), Data: data, Msg: msg}.JSON(c)
}

// SuccessWithData 仅返回数据
func SuccessWithData(c *gin.Context, data any) {
	response{Code: int(SuccessCode), Data: data, Msg: SuccessCode.String()}.JSON(c)
}

// SuccessWithMsg 仅返回消息
func SuccessWithMsg(c *gin.Context, msg string) {
	response{Code: int(SuccessCode), Data: Empty, Msg: msg}.JSON(c)
}

// SuccessWithList 分页列表响应
func SuccessWithList(c *gin.Context, list any, count int64) {
	response{
		Code: int(SuccessCode),
		Data: map[string]any{
			"list":  list,
			"count": count,
		},
		Msg: SuccessCode.String(),
	}.JSON(c)
}

// FailWithMsg 通用失败（消息）
func FailWithMsg(c *gin.Context, msg string) {
	response{Code: int(FailServiceCode), Data: Empty, Msg: msg}.JSON(c)
}

// FailWithData 失败 + 自定义数据（慎用）
func FailWithData(c *gin.Context, msg string, data any) {
	response{Code: int(FailServiceCode), Data: data, Msg: msg}.JSON(c)
}

// FailWithCode 按错误码返回（推荐）
func FailWithCode(c *gin.Context, code Code) {
	response{Code: int(code), Data: Empty, Msg: code.String()}.JSON(c)
}

// FailWithMsgAndCode 自定义消息 + 错误码
func FailWithMsgAndCode(c *gin.Context, code Code, msg string) {
	response{Code: int(code), Data: Empty, Msg: msg}.JSON(c)
}

// FailWithError
func FailWithError(c *gin.Context, err error) {
	FailWithData(c, err.Error(), Empty)
}

// 便捷方法
func FailAuth(c *gin.Context) {
	FailWithCode(c, FailAuthCode)
}

func FailToken(c *gin.Context) {
	FailWithCode(c, FailTokenCode)
}

func FailService(c *gin.Context) {
	FailWithCode(c, FailServiceCode)
}

func FailValid(c *gin.Context, msg string) {
	FailWithMsgAndCode(c, FailValidCode, msg)
}
