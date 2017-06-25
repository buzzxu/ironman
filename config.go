package ironman

import "github.com/labstack/echo"

type (
	// Error 异常结构
	Error struct {
		Code    int         `json:"code"`
		Key     string      `json:"error,omitempty"`
		Message interface{} `json:"message"`
	}
	// Result 返回结果
	Result struct {
		Code int         `json:"code"`
		Data interface{} `json:"data,omitempty"`
	}
)

// JSON 输出json
func JSON(c echo.Context, result *Result) error {
	return c.JSON(result.Code, result)
}

// ResultOf 构造Result
func ResultOf(code int, data interface{}) *Result {
	return &Result{
		Code: code,
		Data: data,
	}
}
