package helpers

import (
	"strings"

	"github.com/samber/lo"
)

func TrimLeft(text string, word string) string {
	runes := lo.Uniq([]rune(word))
	return strings.TrimLeft(text, string(runes))
}

func Foo() string {
	return "foo"
}
