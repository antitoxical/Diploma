package web

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

func OkResp(status string, data interface{}) Response {
	return Response{
		Status: status,
		Data:   data,
	}
}

func ErrorResp(status string, data interface{}) Response {
	return Response{
		Status: status,
		Data:   data,
	}
}
