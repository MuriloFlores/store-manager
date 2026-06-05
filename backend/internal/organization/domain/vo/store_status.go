package vo

import (
	"errors"
	"strings"
)

type StoreStatus string

const (
	StatusPending     StoreStatus = "PENDING"
	StatusActive      StoreStatus = "ACTIVE"
	StatusFailed      StoreStatus = "FAILED"
	StatusDeactivated StoreStatus = "DEACTIVATED"
)

func NewStoreStatus(value string) (StoreStatus, error) {
	normalizedValue := StoreStatus(strings.ToUpper(strings.TrimSpace(value)))

	switch normalizedValue {
	case StatusPending, StatusActive, StatusFailed, StatusDeactivated:
		return normalizedValue, nil
	default:
		return "", errors.New("invalid store status")
	}
}

func (s StoreStatus) String() string {
	return string(s)
}
