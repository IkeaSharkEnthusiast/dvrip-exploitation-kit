package stringh

import (
	"fmt"
	"strings"
)

func Between(input, start, end string) (string, error) {
	startIndex := strings.Index(input, start)
	if startIndex == -1 {
		return "", fmt.Errorf("start string not found")
	}

	startIndex += len(start)
	
	endIndex := strings.Index(input[startIndex:], end)
	if endIndex == -1 {
		return "", fmt.Errorf("end string not found")
	}

	return input[startIndex : startIndex+endIndex], nil
}
