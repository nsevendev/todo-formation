package doc

type ErrorModel struct {
	Data string `json:"data"`
	Error    string `json:"error"`
	Message  string `json:"message"`
}