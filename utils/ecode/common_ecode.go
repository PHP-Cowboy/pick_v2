package ecode

var (
	Success                 = New(200, "ok") // 正确
	RequestError            = New(400, "请求错误")
	IllegalRequest          = New(401, "非法请求")
	RequestNotFound         = New(404, "请求不存在")
	ServerErr               = New(500, "服务器错误")
	DataInsertError         = New(501, "数据添加失败")
	DataDeleteError         = New(502, "数据删除失败")
	DataSaveError           = New(503, "数据保存失败")
	DataQueryError          = New(504, "数据查询失败")
	DataNotExist            = New(505, "数据不存在")
	DataQueryTimeOutOfRange = New(506, "超出查询时间范围")

	CommunalSignInvalid    = New(10, "sign参数异常")
	ParamInvalid           = New(11, "参数不合法")
	CommunalSessionInvalid = New(12, "session参数异常")
	CommunalParamInvalid   = New(13, "公共参数异常")
	UserNotLogin           = New(14, "用户未登录")

	UserNotFound            = New(1000, "用户未找到")
	RoleNotFound            = New(1001, "角色未找到")
	DataTransformationError = New(1002, "数据转换出错")
	PasswordCheckFailed     = New(1003, "密码校验有误")
)
