package vo

import (
	"strings"
	"testing"
)

func TestNewStoreName(t *testing.T) {
	t.Run("should create new store name from an input", func(t *testing.T) {
		rawStoreName := "Café do Programador #1"

		storeName, err := NewSchemaName(rawStoreName)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if !strings.HasPrefix(storeName.String(), "tenant_cafe_do_programador_1_") {
			t.Errorf("expected formatated name, go %v", storeName.String())
		}
	})

	t.Run("should not create new store name from an invalid input", func(t *testing.T) {
		rawStoreName := ""

		storeName, err := NewSchemaName(rawStoreName)
		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if storeName != "" {
			t.Errorf("expected empty string, got %v", storeName)
		}
	})

	t.Run("should truncate store name longer than 50 characters", func(t *testing.T) {
		rawName := "Essa loja tem um nome absurdamente grande é ridículo de mais para nao ser um teste"

		storeName, err := NewSchemaName(rawName)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if !strings.HasPrefix(storeName.String(), "tenant_essa_loja_tem_um_nome_absurdamente_grande_e_ridicu_") {
			t.Errorf("expected formatted name, got %v", storeName.String())
		}
	})

	t.Run("should restore store name from valid an input", func(t *testing.T) {
		storeNameFromInput := "tenant_cafe_do_programador_abc12"

		storeName, err := RestoreSchemaName(storeNameFromInput)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if storeName.String() != storeNameFromInput {
			t.Errorf("expected equals names, go %v", storeName.String())
		}
	})

	t.Run("should not restore store name from invalid input", func(t *testing.T) {
		storeNameFromInput := "tenant_cafe; DROP TABLE users"

		storeName, err := RestoreSchemaName(storeNameFromInput)
		if err == nil {
			t.Errorf("expected error, got %v", err)
		}

		if storeName != "" {
			t.Errorf("expected empty string, got %v", storeName.String())
		}
	})

	t.Run("shouldn't be duplicates ", func(t *testing.T) {
		rawName := "loja do programador"

		storeName1, err := NewSchemaName(rawName)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		storeName2, err := NewSchemaName(rawName)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if storeName1.String() == storeName2.String() {
			t.Errorf("expected not equals names, go %v", storeName1.String())
		}
	})
}
