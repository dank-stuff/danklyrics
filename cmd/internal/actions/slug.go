package actions

import (
	"regexp"
	"strings"

	"github.com/mozillazg/go-unidecode"
)

func Slugify(str string) string {
	escapedStr := unidecode.Unidecode(strings.ToLower(str))

	slugBuilder := new(strings.Builder)
	for _, ch := range escapedStr {
		if ch >= 'a' && ch <= 'z' || ch == '-' {
			slugBuilder.WriteRune(ch)
		} else {
			slugBuilder.WriteRune('-')
		}
	}

	return strings.Trim(
		regexp.MustCompile(`[-]+`).
			ReplaceAllString(slugBuilder.String(), "-"),
		"-",
	)
}
