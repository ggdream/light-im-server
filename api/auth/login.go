package auth

import (
	"github.com/gin-gonic/gin"

	"lim/config"
	"lim/pkg/cache"
	"lim/pkg/db"
	"lim/pkg/errno"
	"lim/tools/token"
)

type loginReqModel struct {
	UserID *string `json:"user_id" binding:"required"`
}

type loginResModel struct {
	Token    string `json:"token"`
	ExpireAt int64  `json:"expire_at"`
}

func LoginController(c *gin.Context) {
	var (
		form loginReqModel
		err  error
	)
	if err = c.ShouldBindJSON(&form); err != nil {
		errno.NewFParamInvalid(err.Error()).Reply(c)
		return
	}

	doc := db.UserInfoDoc{}
	err = doc.Search(*form.UserID)
	if err != nil {
		errno.NewF(errno.BaseErrMongo, err.Error(), errno.ErrUserSearchFailed).Reply(c)
		return
	}

	tokenKey := config.GetApp().TokenKey
	tk, expireAt, _ := token.Generate(tokenKey, *form.UserID, 0)

	ca := cache.AuthToken{}
	err = ca.Set(*form.UserID, 0, tk, expireAt)
	if err != nil {
		errno.NewF(errno.BaseErrRedis, err.Error(), errno.ErrAuthLoginFailed).Reply(c)
		return
	}

	ret := &loginResModel{
		Token:    tk,
		ExpireAt: expireAt,
	}
	errno.NewS(ret).Reply(c)
}
