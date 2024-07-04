package boolutil

func BoolToString(value bool) string {
	if value {
		return "1"
	}
	return "0"
}
