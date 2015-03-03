package gof

func BoolString(b bool, t, f string) string {
	if b {
		return t
	} else {
		return f
	}
}
