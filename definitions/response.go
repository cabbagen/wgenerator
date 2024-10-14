package definitions

var ResponseCodeMap map[string]int = map[string]int{
	"fail":             100500,
	"success":          100200,
	"forbidden":        100403,
	"unauthentication": 100401,
}

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func NewResponse(code int, data interface{}, message string) Response {
	return Response{code, data, message}
}
