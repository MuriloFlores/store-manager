package vo

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	timeCost  = 3
	memory    = 64 * 1024
	threads   = 4
	keyLength = 32
)

var (
	ErrPasswordTooShort      = errors.New("password too short")
	ErrInternalError         = errors.New("internal error")
	ErrEmptyPassword         = errors.New("empty password")
	ErrLowPasswordComplexity = errors.New("password with low complexity")
)

type Password string

func NewPassword(plainText, pepper string) (Password, error) {
	if len(plainText) < 8 {
		return "", ErrPasswordTooShort
	}

	if !validateComplexity(plainText) {
		return "", ErrLowPasswordComplexity
	}

	saltedPlain := plainText + pepper

	salt := make([]byte, keyLength/2)
	if _, err := rand.Read(salt); err != nil {
		return "", ErrInternalError
	}

	hashKey := argon2.IDKey([]byte(saltedPlain), salt, timeCost, memory, threads, keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hashKey)

	encodeHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, timeCost, threads, b64Salt, b64Hash)

	return Password(encodeHash), nil
}

func RestorePassword(hash string) (Password, error) {
	if hash == "" {
		return "", ErrEmptyPassword
	}

	return Password(hash), nil
}

func (p Password) Matches(plainText, pepper string) bool {
	parts := strings.Split(string(p), "$")
	if len(parts) != 6 {
		return false
	}

	var memory, time uint32
	var threads uint8

	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	compareHash := argon2.IDKey([]byte(plainText+pepper), salt, time, memory, threads, uint32(len(decodedHash)))

	return subtle.ConstantTimeCompare(decodedHash, compareHash) == 1
}

func (p Password) String() string {
	return string(p)
}

func validateComplexity(p string) bool {
	hasUpper := regexp.MustCompile("[A-Z]").MatchString(p)
	hasLower := regexp.MustCompile("[a-z]").MatchString(p)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(p)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(p)

	return hasUpper && hasNumber && hasSpecial && hasLower
}
