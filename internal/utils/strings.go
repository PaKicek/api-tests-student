package utils

import "unicode"

func MixCase(text string) string {
	runes := []rune(text)
	for i, r := range runes {
		if i%2 == 0 {
			runes[i] = unicode.ToUpper(r)
			continue
		}
		runes[i] = unicode.ToLower(r)
	}

	return string(runes)
}
