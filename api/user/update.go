package user

import (
	"github.com/gin-gonic/gin"

	"lim/pkg/db"
	"lim/pkg/errno"
)

type updateReqModel struct {
	UserID   *string `json:"user_id" binding:"required"`
	Avatar   string  `json:"avatar"`
	Nickname string  `json:"nickname"`
}

func UpdateController(c *gin.Context) {
	var (
		form updateReqModel
		err  error
	)
	if err = c.ShouldBindJSON(&form); err != nil {
		errno.NewFParamInvalid(err.Error()).Reply(c)
		return
	}

	doc := db.UserInfoDoc{}
	err = doc.Update(*form.UserID, form.Nickname, form.Avatar)
	if err != nil {
		errno.NewF(errno.BaseErrMongo, err.Error(), errno.ErrUserUpdateFailed).Reply(c)
		return
	}

	errno.NewS(nil).Reply(c)
}
