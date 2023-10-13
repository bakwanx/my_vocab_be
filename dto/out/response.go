package out

type Response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
