package errno

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"lim/config"
)

type Response struct {
	data   any
	err    error
	retErr Errno
}

func NewResp(data any, baseErr error, errDesc string, retErr Errno) *Response {
	return &Response{
		data:   data,
		err:    errors.Wrap(baseErr, errDesc),
		retErr: retErr,
	}
}

func NewF(baseErr error, errDesc string, retErr Errno) *Response {
	return NewResp(nil, baseErr, errDesc, retErr)
}

func NewS(data any) *Response {
	return NewResp(data, nil, "", ErrSuccess)
}

func NewFParamInvalid(errDesc string) *Response {
	return NewF(BaseErrParam, errDesc, ErrParamInvalid)
}

func (r *Response) Reply(c *gin.Context) {
	if r.err != nil {
		config.CtxKeyManager.SetError(c, r.err.Error())
	}
	data := retBody(r.data, r.retErr)

	switch r.retErr {
	case ErrSuccess:
		c.JSON(http.StatusOK, data)
	default:
		c.AbortWithStatusJSON(http.StatusOK, data)
	}
}

func retBody(data any, err Errno) gin.H {
	res := gin.H{"code": err.Int(), "desc": err.String()}
	if data != nil {
		res["data"] = data
	}

	return res
}

// Well 正常顺利响应
func Well(c *gin.Context, data any, desc ...string) {
	New(c, ErrSuccess, data, desc...)
}

// Fail 发生错误响应
func Fail(c *gin.Context, errno Errno, desc ...string) {
	msg := errno.String()
	if len(desc) > 0 {
		msg = desc[0]
	}

	obj := gin.H{
		"code": errno.Int(),
		"desc": msg,
	}

	c.AbortWithStatusJSON(http.StatusOK, obj)
}

// New 响应接口
func New(c *gin.Context, errno Errno, data any, desc ...string) {
	msg := errno.String()
	if len(desc) > 0 {
		msg = desc[0]
	}

	obj := gin.H{
		"code": errno.Int(),
		"data": data,
		"desc": msg,
	}

	c.JSON(http.StatusOK, obj)
}

// NewError 自定义创建新错误
func NewError(message string) error {
	return errors.New(message)
}
