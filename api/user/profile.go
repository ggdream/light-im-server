package user

import (
	"github.com/gin-gonic/gin"

	"lim/pkg/db"
	"lim/pkg/errno"
)

type profileReqModel struct {
	UserID *string `json:"user_id" binding:"required"`
}

type profileResModel struct {
	UserID   string `json:"user_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func ProfileController(c *gin.Context) {
	var (
		form profileReqModel
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

	ret := &profileResModel{
		UserID:   *form.UserID,
		Nickname: doc.Username,
		Avatar:   doc.Avatar,
	}
	errno.NewS(ret).Reply(c)
}
