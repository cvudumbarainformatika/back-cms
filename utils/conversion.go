package utils

import "time"

// StringToPtr returns a pointer to the string if it's not empty, otherwise nil
func StringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Int64ToPtr returns a pointer to the int64 value
func Int64ToPtr(i int64) *int64 {
	return &i
}

// TimeToPtr returns a pointer to the time value
func TimeToPtr(t time.Time) *time.Time {
	return &t
}

// BoolToPtr returns a pointer to the bool value
func BoolToPtr(b bool) *bool {
	return &b
}

// PtrToString safely dereferences a string pointer
func PtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
