package config

import "github.com/gin-gonic/gin"

const (
	ctxKeyUserID   = "app:uid"
	ctxKeyAppError = "app:err"
)

var (
	CtxKeyManager = &ctxKeyCtrl{}
)

type ctxKeyCtrl struct{}

func (ctxKeyCtrl) SetUserID(ctx *gin.Context, userId string) {
	ctx.Set(ctxKeyUserID, userId)
}

func (ctxKeyCtrl) GetUserID(ctx *gin.Context) string {
	return ctx.GetString(ctxKeyUserID)
}

func (ctxKeyCtrl) SetError(ctx *gin.Context, err string) {
	ctx.Set(ctxKeyAppError, err)
}

func (ctxKeyCtrl) GetError(ctx *gin.Context) string {
	return ctx.GetString(ctxKeyAppError)
}
