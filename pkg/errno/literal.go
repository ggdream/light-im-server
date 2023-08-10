package errno

var literalMap = map[Errno]string{
	ErrSuccess:             "操作成功",
	ErrParamInvalid:        "参数格式有误",
	ErrUserCreateFailed:    "用户创建失败",
	ErrUserDeleteFailed:    "用户删除失败",
	ErrUserUpdateFailed:    "用户更新失败",
	ErrUserSearchFailed:    "用户查找失败",
	ErrAuthLoginFailed:     "用户登录失败",
	ErrAuthLogoutFailed:    "用户退出登录失败",
	ErrChatMsgSaveFailed:   "聊天消息存储失败",
	ErrChatMsgMarkFailed:   "聊天消息标记已读失败",
	ErrChatMsgUnreadFailed: "获取消息未读数失败",
	ErrChatConvListFailed:  "聊天会话列表获取失败",
	ErrChatConvDelFailed:   "聊天会话删除失败",
	ErrJWT:                 "登录凭证非法",
	ErrMySQL:               "数据库出错啦",
	ErrRedis:               "缓存出错啦",
}
