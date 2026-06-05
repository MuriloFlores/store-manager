package vo

import (
	"strings"
	"testing"
)

func TestNewStatusName(t *testing.T) {
	t.Run("should return new status name", func(t *testing.T) {
		rawStatus := "  pending"
		rawStatus2 := "PENDING"

		statusName, err := NewStoreStatus(rawStatus)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if strings.ToUpper(strings.TrimSpace(rawStatus)) != statusName.String() {
			t.Errorf("expected normalized value, got %v ", statusName.String())
		}

		if rawStatus2 != statusName.String() {
			t.Errorf("expected normalized value, got %v ", statusName.String())
		}
	})

	t.Run("should return error if store status is invalid", func(t *testing.T) {
		rawStatus := "  invalid status"

		statusName, err := NewStoreStatus(rawStatus)
		if err == nil {
			t.Errorf("expected error, got none")
		}

		if statusName != "" {
			t.Errorf("expected empty string, got %v", statusName.String())
		}
	})
}
