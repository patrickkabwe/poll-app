package views

import "context"

type contextKey string

const (
	FieldValue contextKey = "field_value"
	UserKey    contextKey = "user"
)

func GetFieldValue(c context.Context, key string) string {
	if value, ok := c.Value(FieldValue).(map[string]string); ok {
		return value[key]
	}
	return ""
}

func IsAuthenticated(c context.Context) bool {
	if value, ok := c.Value(UserKey).(bool); ok {
		return value
	}
	return false
}
