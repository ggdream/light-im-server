package user

import (
	"github.com/gin-gonic/gin"

	"lim/pkg/cache"
	"lim/pkg/db"
	"lim/pkg/errno"
)

type deleteReqModel struct {
	UserID *string `json:"user_id" binding:"required"`
}

func DeleteController(c *gin.Context) {
	var (
		form deleteReqModel
		err  error
	)
	if err = c.ShouldBindJSON(&form); err != nil {
		errno.NewFParamInvalid(err.Error()).Reply(c)
		return
	}

	doc := db.UserInfoDoc{}
	err = doc.Delete(*form.UserID)
	if err != nil {
		errno.NewF(errno.BaseErrMongo, err.Error(), errno.ErrUserDeleteFailed).Reply(c)
		return
	}

	ca := cache.AuthToken{}
	_ = ca.Del(*form.UserID)

	errno.NewS(nil).Reply(c)
}
