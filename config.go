package ironman

type (
	Error struct {
		Code    int    `json:"code"`
		Key     string `json:"error"`
		Message string `json:"message"`
	}
	Result struct {
		Code    int                    `json:"code"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
)

func ResultOf(code int,message string,data map[string]interface{}) Result  {
	return Result{
		Code:code,
		Data:data,
	}
}
