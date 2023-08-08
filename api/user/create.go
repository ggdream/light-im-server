package user

import (
	"github.com/gin-gonic/gin"

	"lim/pkg/db"
	"lim/pkg/errno"
)

type createReqModel struct {
	UserID   *string `json:"user_id" binding:"required"`
	Avatar   string  `json:"avatar"`
	Nickname string  `json:"nickname"`
}

func CreateController(c *gin.Context) {
	var (
		form createReqModel
		err  error
	)
	if err = c.ShouldBindJSON(&form); err != nil {
		errno.NewFParamInvalid(err.Error()).Reply(c)
		return
	}

	doc := db.UserInfoDoc{}
	err = doc.Create(*form.UserID, form.Nickname, form.Avatar)
	if err != nil {
		errno.NewF(errno.BaseErrMongo, err.Error(), errno.ErrUserCreateFailed).Reply(c)
		return
	}

	errno.NewS(nil).Reply(c)
}
