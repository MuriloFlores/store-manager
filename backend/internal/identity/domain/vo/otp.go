package vo

import (
	"crypto/rand"
	"errors"
	"math/big"
	"regexp"
)

var (
	ErrInvalidOTPFormat = errors.New("OTP must contain exactly 6 digits")
	ErrGeneratingOTP    = errors.New("failed to generate secure OTP")
)

type OTP string

func NewOTP(value string) (OTP, error) {
	matched, _ := regexp.MatchString(`^\d{6}$`, value)
	if !matched {
		return "", ErrInvalidOTPFormat
	}

	return OTP(value), nil
}

func GenerateOTP() (OTP, error) {
	const otpChars = "0123456789"
	const length = 6
	buffer := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(otpChars))))
		if err != nil {
			return "", ErrGeneratingOTP
		}

		buffer[i] = otpChars[num.Int64()]
	}

	return NewOTP(string(buffer))
}

func (o OTP) String() string {
	return string(o)
}
