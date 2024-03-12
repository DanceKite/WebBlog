package controller

type ResCode int64

const (
	// CodeSuccess 成功的code
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeNeedLogin
	CodeInvalidToken
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户已经存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务端繁忙",

	CodeNeedLogin:    "需要登录",
	CodeInvalidToken: "无效的Token",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if ok {
		return msg
	}
	return codeMsgMap[CodeServerBusy]
}
