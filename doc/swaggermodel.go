package doc

type Meta struct {
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	Timestamp string `json:"timestamp"`
}

type ResponseModel struct {
	Data    any `json:"data"`
	Error   any `json:"error"`
	Message string      `json:"message"`
	Meta    Meta        `json:"meta"`
}

type InsufficientPermissionsResponseModel struct {
	Error   string `json:"error"`
}