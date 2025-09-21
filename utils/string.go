package utils

import (
	"encoding/json"
	"fmt"
)

func StringThis(o any) string {
	js, _ := json.MarshalIndent(o, "|", "   ")
	return fmt.Sprintf(string(js))
}
