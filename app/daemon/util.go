package daemon

import (
	"strconv"
	"strings"
	"time"
)

func getTick(t time.Time) string {
	d := strconv.Itoa(t.Day())
	h := strconv.Itoa(t.Hour())
	tk := strconv.Itoa(t.Minute() / 15)
	return strings.Join([]string{d, h, tk}, "-")
}
