package response

type ApiResponseResult struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
type NullJson struct{}

func Success(data interface{}) *ApiResponseResult {
	return &ApiResponseResult{SUCCESS, MapErrMsg(SUCCESS), data}
}

func Fail(errCode uint32, errMsg string) *ApiResponseResult {
	return &ApiResponseResult{Code: errCode, Msg: errMsg}
}
