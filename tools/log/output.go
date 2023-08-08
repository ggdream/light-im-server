package log

func Debugf(format string, args ...any) {
	log.Debugf(format, args...)
}

func Infof(format string, args ...any) {
	log.Infof(format, args...)
}

func Warnf(format string, args ...any) {
	log.Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	log.Errorf(format, args...)
}

func Panicf(format string, args ...any) {
	log.Panicf(format, args...)
}

func DebugWf(format string, fields Map, args ...any) {
	log.DebugWf(format, fields, args...)
}

func InfoWf(format string, fields Map, args ...any) {
	log.InfoWf(format, fields, args...)
}

func WarnWf(format string, fields Map, args ...any) {
	log.WarnWf(format, fields, args...)
}

func ErrorWf(format string, fields Map, args ...any) {
	log.ErrorWf(format, fields, args...)
}

func PanicWf(format string, fields Map, args ...any) {
	log.PanicWf(format, fields, args...)
}
