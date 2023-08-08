package log

type (
	LogType = string
)

const (
	LogTypeMiddleWare LogType = "middleware"
	LogTypeRollback   LogType = "rollback"
)

type BusinessModel struct {
	Type   string `json:"type"`
	Method string `json:"method"`
	Path   string `json:"path"`
	Time   string `json:"time"`
	// Form   any    `json:"form"`
	Cost   int64  `json:"cost"`
	UserID string `json:"user_id"`
	UserIP string `json:"user_ip"`
	Role   uint   `json:"role"`
	Error  string `json:"error"`
}

func (m *BusinessModel) Encode() map[string]any {
	return map[string]any{
		"type":   m.Type,
		"method": m.Method,
		"path":   m.Path,
		"time":   m.Time,
		// "form":    m.Form,
		"cost":    m.Cost,
		"error":   m.Error,
		"user_id": m.UserID,
		"user_ip": m.UserIP,
		"role":    m.Role,
	}
}

type RollbackModel struct {
	// Type   string `json:"type"`
	Scene  string `json:"scene"`
	Path   string `json:"path"`
	UserID uint   `json:"user_id"`
	Role   uint8  `json:"role"`
	Error  string `json:"error"`
	Extra  string `json:"extra"`
}

func (m *RollbackModel) Encode() map[string]any {
	return map[string]any{
		"type":    LogTypeRollback,
		"scene":   m.Scene,
		"path":    m.Path,
		"error":   m.Error,
		"user_id": m.UserID,
		"role":    m.Role,
		"extra":   m.Extra,
	}
}
