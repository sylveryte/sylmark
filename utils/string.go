package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

func StringThis(o any) string {
	js, _ := json.MarshalIndent(o, "|", "   ")
	return fmt.Sprintf(string(js))
}

func FormatDate(t time.Time) string {
	return t.Format(time.DateOnly)
}
