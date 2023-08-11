package errno

import "github.com/pkg/errors"

type Errno uint32

const (
	ErrSuccess   Errno = 0
	ErrUniversal Errno = 999

	ErrParamInvalid              Errno = 1
	ErrUserCreateFailed          Errno = 2
	ErrUserDeleteFailed          Errno = 3
	ErrUserUpdateFailed          Errno = 4
	ErrUserSearchFailed          Errno = 5
	ErrAuthLoginFailed           Errno = 6
	ErrAuthLogoutFailed          Errno = 7
	ErrChatMsgSaveFailed         Errno = 8
	ErrChatMsgMarkFailed         Errno = 9
	ErrChatConvListFailed        Errno = 10
	ErrChatConvDelFailed         Errno = 11
	ErrChatMsgUnreadFailed       Errno = 12
	ErrUUIDGenFailed             Errno = 13
	ErrOSSPresignPutURLGenFailed Errno = 14
	ErrJWT                       Errno = 1000
	ErrMySQL                     Errno = 1002
	ErrRedis                     Errno = 1003
	ErrCopier                    Errno = 1004
)

func (e Errno) Int() int {
	return int(e)
}

func (e Errno) String() string {
	return literalMap[e]
}

func (e Errno) Error() error {
	return errors.New(e.String())
}
