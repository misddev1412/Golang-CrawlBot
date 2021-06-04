package packages

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

var re = regexp.MustCompile("[^a-z0-9]+")
var t = transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

func GenerateSlug(s string) string {

	s = FixUnicode(s)
	result, _, _ := transform.String(t, s)

	return strings.Trim(re.ReplaceAllString(strings.ToLower(result), "-"), "-")
}

/**
 * Short the string by word count first then do slug
 */

func GenerateShortSlug(s string, word_count int) string {
	arr_words := strings.Split(s, " ")

	if len(arr_words) < word_count {
		return GenerateSlug(s)
	}

	arr_words = arr_words[:word_count]

	return GenerateSlug(strings.Join(arr_words, " "))

}

/**
 * Logic some cases that unicode lib can not resolve
 */
func FixUnicode(s string) string {
	r := strings.NewReplacer("đ", "d")
	return r.Replace(s)
}
