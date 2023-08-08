package errno

import (
	"github.com/pkg/errors"
)

var (
	BaseErrParam   = errors.New("Param")
	BaseErrMongo   = errors.New("Mongo")
	BaseErrRedis   = errors.New("Redis")
	BaseErrNSQ     = errors.New("NSQ")
	BaseErrOSS     = errors.New("OSS")
	BaseErrTools   = errors.New("Tools")
	BaseErrInvalid = errors.New("Invalid")
	BaseErrUnknown = errors.New("Unknown")
)
