// internal/entity/uuid.go
package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrorEntityIDInvalid = errors.New("invalid id, must be a valid UUID")
	ErrorFailedToScan    = errors.New("failed to scan EntityID")
)

type EntityID struct {
	value uuid.UUID
}

func NewEntityID() EntityID {
	return EntityID{
		value: uuid.New(),
	}
}

func ParseEntityID(s string) (EntityID, error) {
	id, err := uuid.Parse(strings.TrimSpace(s))

	if err != nil {
		return EntityID{}, ErrorEntityIDInvalid
	}

	return EntityID{id}, nil
}

func (e EntityID) String() string {
	return e.value.String()
}

func (e EntityID) Equals(other EntityID) bool {
	return e.value == other.value
}

func (e *EntityID) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		return e.unmarshalString(v)
	case []byte:
		return e.unmarshalString(string(v))
	default:
		return ErrorFailedToScan
	}
}

func (e *EntityID) unmarshalString(s string) error {
	id, err := uuid.Parse(strings.TrimSpace(s))

	if err != nil {
		return ErrorEntityIDInvalid
	}

	e.value = id

	return nil
}

func (e EntityID) Value() (driver.Value, error) {
	return e.String(), nil
}

func (e EntityID) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *EntityID) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	id, err := uuid.Parse(s)

	if err != nil {
		return ErrorEntityIDInvalid
	}

	e.value = id

	return nil
}
