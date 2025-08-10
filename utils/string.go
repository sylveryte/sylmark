package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

func StringThis(o any) string {
	js, _ := json.MarshalIndent(o, "|", "   ")
	return fmt.Sprintf(string(js))
}

func FindWord(char int, line string) (before string, after string, found bool) {
	if len(line) > 0 {
		if char > 0 {
			// before
			for i := char - 1; ; {
				fmt.Printf("i=%d\n", i)
				if i != 0 && line[i] != ' ' {
					i--
					continue
				} else {
					found = true
					before = line[i:char]
					break
				}
			}

			// no need to waste effort if already determined that it is tag
			if strings.TrimSpace(before)[0] != '#' {
				// after
				for i := char; ; {
					if i != len(line) && line[i] != ' ' {
						i++
						continue
					} else {
						found = true
						after = line[char:i]
						break
					}
				}
			}
		}
	}

	return
}
