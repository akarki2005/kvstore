package protocol

import (
	"fmt"
	"strings"
)

func ParseCommand(command string) (op, key, value string, err error) {
	words := strings.Fields(strings.TrimSpace(command))
	if len(words) == 0 {
		return "", "", "", fmt.Errorf("Empty command.")
	}

	op = strings.ToUpper(words[0])

	switch op {
	case "GET", "DELETE":
		if len(words) != 2 {
			return "", "", "", fmt.Errorf("Usage: %s [key]", op)
		}
		key = words[1]
	case "SET":
		if len(words) != 3 {
			return "", "", "", fmt.Errorf("Usage: SET [key] [value]")
		}
		key = words[1]
		value = words[2]
	default:
		return "", "", "", fmt.Errorf("Unknown operation: %s", op)
	}

	return op, key, value, nil
}