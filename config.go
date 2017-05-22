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
