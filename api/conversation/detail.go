package conversation

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"lim/config"
	"lim/pkg/db"
	"lim/pkg/errno"
)

type detailReqModel struct {
	UserID *string `json:"user_id" binding:"required"`
}

type detailResModel struct {
	UserID         string `json:"user_id"`
	ConversationID string `json:"conversation_id"`
	Nickname       string `json:"nickname"`
	Avatar         string `json:"avatar"`
}

func DetailController(c *gin.Context) {
	var (
		form detailReqModel
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

	userId := config.CtxKeyManager.GetUserID(c)
	ret := &detailResModel{
		UserID:         *form.UserID,
		ConversationID: genCID(userId, *form.UserID),
		Nickname:       doc.Username,
		Avatar:         doc.Avatar,
	}
	errno.NewS(ret).Reply(c)
}

func genCID(a, b string) string {
	if strings.Compare(a, b) == 1 {
		a, b = b, a
	}

	return fmt.Sprintf("c_%s_%s", a, b)
}
