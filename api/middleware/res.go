package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// LangResponse 多语言消息映射
var LangResponse = map[string]map[int]string{
	"zh": {
		0:    "成功",
		400:  "请求参数错误",
		401:  "未授权",
		403:  "禁止访问",
		404:  "资源不存在",
		500:  "服务器内部错误",
		1000: "业务错误",
	},
	"en": {
		0:    "Success",
		400:  "Bad request",
		401:  "Unauthorized",
		403:  "Forbidden",
		404:  "Not found",
		500:  "Internal server error",
		1000: "Business error",
	},
}

// GetLang 从cookie获取用户语言设置
func GetLang(c *gin.Context) string {
	lang, err := c.Cookie("lang")
	if err != nil || (lang != "zh" && lang != "en") {
		return "zh" // 默认中文
	}
	return lang
}

func GetMsg(c *gin.Context, code int) string {
	lang := GetLang(c)
	msg := LangResponse[lang][code]
	if msg == "" {
		msg = LangResponse[lang][1000]
	}
	return msg
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	lang := GetLang(c)
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: LangResponse[lang][0],
		Data:    data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, data interface{}) {
	lang := GetLang(c)
	msg := LangResponse[lang][code]
	if msg == "" {
		msg = LangResponse[lang][1000] // 默认业务错误消息
	}

	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

// FailWithHTTPStatus 带HTTP状态码的失败响应
func FailWithHTTPStatus(c *gin.Context, httpCode, code int, data interface{}) {
	lang := GetLang(c)
	msg := LangResponse[lang][code]
	if msg == "" {
		msg = LangResponse[lang][1000]
	}

	c.JSON(httpCode, Response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}
