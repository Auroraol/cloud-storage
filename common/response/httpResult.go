package response

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

// http返回
func HttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err == nil {
		// 成功返回
		result := Success(resp)
		httpx.WriteJson(w, http.StatusOK, result)
	} else {
		// 错误返回
		errcode := SYSTEM_ERROR
		errmsg := "服务器开小差啦，稍后再来试一试"

		causeErr := errors.Cause(err) // 获取错误原因
		if e, ok := causeErr.(*CodeError); ok {
			// 自定义CodeError
			errcode = e.GetErrCode()
			errmsg = e.GetErrMsg()
		} else if gstatus, ok := status.FromError(causeErr); ok {
			// gRPC 错误
			grpcCode := uint32(gstatus.Code())
			if IsCodeErr(grpcCode) {
				errcode = grpcCode
				errmsg = gstatus.Message()
			} else {
				// 对于非自定义的 gRPC 错误，返回通用的服务器错误
				errcode = SYSTEM_ERROR
				errmsg = "服务器内部错误"
			}
		}

		// 记录错误日志
		logx.WithContext(r.Context()).Errorf("【API-ERR】 : %+v ", err)

		// 返回错误响应
		httpx.WriteJson(w, http.StatusBadRequest, Fail(errcode, errmsg))
	}
}

// 授权的http方法
func AuthHttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {

	if err == nil {
		//成功返回
		r := Success(resp)
		httpx.WriteJson(w, http.StatusOK, r)
	} else {
		//错误返回
		errcode := SYSTEM_ERROR
		errmsg := "服务器开小差啦，稍后再来试一试"

		causeErr := errors.Cause(err)           // err类型
		if e, ok := causeErr.(*CodeError); ok { //自定义错误类型
			//自定义CodeError
			errcode = e.GetErrCode()
			errmsg = e.GetErrMsg()
		} else {
			if gstatus, ok := status.FromError(causeErr); ok { // grpc err错误
				grpcCode := uint32(gstatus.Code())
				if IsCodeErr(grpcCode) { //区分自定义错误跟系统底层、db等错误，底层、db错误不能返回给前端
					errcode = grpcCode
					errmsg = gstatus.Message()
				}
			}
		}

		logx.WithContext(r.Context()).Errorf("【GATEWAY-ERR】 : %+v ", err)

		httpx.WriteJson(w, http.StatusUnauthorized, Fail(errcode, errmsg))
	}
}

// http 参数错误返回
func ParamErrorResult(r *http.Request, w http.ResponseWriter, err error) {
	errMsg := fmt.Sprintf("%s ,%s", MapErrMsg(DATA_PARAM_ERROR), err.Error())
	httpx.WriteJson(w, http.StatusBadRequest, Fail(DATA_PARAM_ERROR, errMsg))
}
