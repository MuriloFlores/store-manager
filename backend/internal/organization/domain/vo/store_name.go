package vo

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	"github.com/lithammer/shortuuid/v4"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	schemaName           = regexp.MustCompile(`^tenant_[a-z0-9_]+$`)
	nonAlphaNumericRegex = regexp.MustCompile("[^a-z0-9_]+")
	alphabet             = "0123456789abcdefghijklmnopqrstuvwxyz"
)

type SchemaName string

func RestoreSchemaName(value string) (SchemaName, error) {
	if !schemaName.MatchString(value) {
		return "", errors.New("invalid store name")
	}

	return SchemaName(value), nil
}

func NewSchemaName(value string) (SchemaName, error) {
	slug := strings.ToLower(value)

	slug, err := removeAccents(slug)
	if err != nil {
		return "", err
	}

	slug = nonAlphaNumericRegex.ReplaceAllString(slug, "_")
	slug = strings.Trim(slug, "_")

	if len(slug) > 50 {
		slug = slug[:50]
	}

	if slug == "" {
		return "", errors.New("invalid store name")
	}

	suffix := shortuuid.NewWithAlphabet(alphabet)

	finalName := "tenant_" + slug + "_" + suffix

	return SchemaName(finalName), nil
}

func removeAccents(value string) (string, error) {
	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC)

	result, _, err := transform.String(t, value)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s *SchemaName) String() string {
	return string(*s)
}
