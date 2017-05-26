package ironman

import "github.com/labstack/echo"

type (
	Error struct {
		Code    int    `json:"code"`
		Key     string `json:"error"`
		Message interface{} `json:"message"`
	}
	Result struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
	}

)

// JSON输出
func JSON(c echo.Context, result *Result) error {
	return c.JSON(result.Code, result)
}
func ResultOf(code int, data interface{}) *Result {
	return &Result{
		Code: code,
		Data: data,
	}
}
