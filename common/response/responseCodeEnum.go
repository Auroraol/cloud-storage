package response

// 定义错误码和消息
const (
	// 成功
	SUCCESS uint32 = 200000

	// 后端服务异常以500开头
	SYSTEM_ERROR    uint32 = 500000 // 服务异常
	OPERATION_ERROR uint32 = 500001 // 操作失败，请稍后再试

	// 后端服务异常以400开头
	DATA_PARAM_ERROR              uint32 = 400000 // 传入参数错误
	ACCOUNT_ERROR                 uint32 = 400001 // 用户名或密码错误
	ACCOUNT_NOT_FOUND             uint32 = 400002 // 账号不存在
	ACCOUNT_LOCK                  uint32 = 400003 // 账号已锁定，请联系管理员解锁
	ACCOUNT_DISABLE               uint32 = 400004 // 账号已禁用
	ACCOUNT_EXPIRED               uint32 = 400005 // 账号已过期
	ACCOUNT_LOCKED                uint32 = 400006 // 账号已锁定
	CREDENTIALS_EXPIRED           uint32 = 400007 // 凭证已过期
	ACCESS_DENIED                 uint32 = 400008 // 不允许访问
	PERMISSION_DENIED             uint32 = 400009 // 无权限访问
	CREDENTIALS_INVALID           uint32 = 400010 // 凭证无效或已过期
	REFRESH_CREDENTIALS_INVALID   uint32 = 400011 // 刷新凭证无效或已过期
	INVALID_REQUEST               uint32 = 400012 // 无效请求
	REQUEST_LIMIT                 uint32 = 400013 // 接口限流
	ACCOUNT_REGISTERED_ERROR      uint32 = 400014 // 用户名已被注册
	ACCOUNT_UPDATETNICKNAME_ERROR uint32 = 400015 // 昵称更新失败
	DB_ERROR                      uint32 = 400016 // 数据库错误
	DB_UPDATE_AFFECTED_ZERO_ERROR uint32 = 400017 // 更新数据受影响行数为0

	FILE_TOO_LARGE_ERROR uint32 = 400018 // 文件过大
)

// 定义错误消息映射
var message = map[uint32]string{
	SUCCESS:                       "SUCCESS",
	SYSTEM_ERROR:                  "服务异常",
	OPERATION_ERROR:               "操作失败，请稍后再试",
	DATA_PARAM_ERROR:              "传入参数错误",
	ACCOUNT_ERROR:                 "用户名或密码错误",
	ACCOUNT_NOT_FOUND:             "账号不存在",
	ACCOUNT_LOCK:                  "账号已锁定，请联系管理员解锁",
	ACCOUNT_DISABLE:               "账号已禁用",
	ACCOUNT_EXPIRED:               "账号已过期",
	ACCOUNT_LOCKED:                "账号已锁定",
	CREDENTIALS_EXPIRED:           "凭证已过期",
	ACCESS_DENIED:                 "授权失败, 不允许访问",
	PERMISSION_DENIED:             "认证失败, 无权限访问, 请重新登录",
	CREDENTIALS_INVALID:           "凭证无效或已过期",
	REFRESH_CREDENTIALS_INVALID:   "刷新凭证无效或已过期",
	INVALID_REQUEST:               "无效请求或请求不接受",
	REQUEST_LIMIT:                 "接口访问次数限制",
	ACCOUNT_REGISTERED_ERROR:      "用户名已被注册",
	ACCOUNT_UPDATETNICKNAME_ERROR: "昵称更新失败",
	DB_ERROR:                      "数据库错误",
	DB_UPDATE_AFFECTED_ZERO_ERROR: "更新数据受影响行数为0",
	FILE_TOO_LARGE_ERROR:          "文件过大",
}

func MapErrMsg(errCode uint32) string {
	if msg, ok := message[errCode]; ok {
		return msg
	} else {
		return "服务器开小差啦,稍后再来试一试"
	}
}

func IsCodeErr(errCode uint32) bool {
	if _, ok := message[errCode]; ok {
		return true
	} else {
		return false
	}
}
