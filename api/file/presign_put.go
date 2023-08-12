package file

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"

	"lim/config"
	"lim/pkg/errno"
	"lim/pkg/oss"
	"lim/tools/fext"
)

type presignPutURLReqModel struct {
	Type        uint8   `json:"type"`
	Name        string  `json:"name"`
	Size        int64   `json:"size"`
	ContentType *string `json:"content_type" binding:"required"`
	Extension   string  `json:"ext"`
}

type presignPutURLResModel struct {
	PresignURL string `json:"presign_url"`
	URL        string `json:"url"`
}

func PresignPutURLController(c *gin.Context) {
	var (
		form presignPutURLReqModel
		err  error
	)
	if err = c.ShouldBindJSON(&form); err != nil {
		errno.NewFParamInvalid(err.Error()).Reply(c)
		return
	}

	u, err := uuid.NewV4()
	if err != nil {
		errno.NewF(errno.BaseErrTools, err.Error(), errno.ErrUUIDGenFailed).Reply(c)
		return
	}

	userId := config.CtxKeyManager.GetUserID(c)
	name := fmt.Sprintf("%s/%s%s", userId, u.String(), fext.MustMimeToExt(*form.ContentType))
	presignUrl, url, err := oss.Client().PresignPutURL(name, time.Minute*5)
	if err != nil {
		errno.NewF(errno.BaseErrOSS, err.Error(), errno.ErrOSSPresignPutURLGenFailed).Reply(c)
		return
	}

	ret := &presignPutURLResModel{
		PresignURL: presignUrl,
		URL:        url,
	}
	errno.NewS(ret).Reply(c)
}
