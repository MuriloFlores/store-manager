package entity

import (
	"strings"
	"testing"

	"github.com/MuriloFlores/order-manager/internal/organization/domain/vo"
	"github.com/google/uuid"
)

func TestNewStore(t *testing.T) {
	t.Run("should create a valid store with pending status and sanitized slug", func(t *testing.T) {
		// Arrange (Preparação)
		ownerID := uuid.New()
		storeName := "Café com Programação"

		// Act (Ação)
		// Aqui chamamos a função que AINDA NÃO EXISTE.
		// Isso é o TDD: definir a interface antes da implementação.
		store, err := NewStore(storeName, ownerID)

		// Assert (Verificação)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if store.Name != storeName {
			t.Errorf("expected name %s, got %s", storeName, store.Name)
		}

		if store.Status != vo.StatusPending && store.Status.String() != "PENDING" {
			t.Errorf("expected status Pending, got %s", store.Status)
		}

		// Verificação do Slug (a regra que você definiu)
		// Deve ser minúsculo, sem acentos, espaços viram '_' e ter um sufixo
		if !strings.Contains(store.SchemaName.String(), "cafe_com_programacao") {
			t.Errorf("slug %s does not contain sanitized name", store.SchemaName)
		}

		if store.OwnerID != ownerID {
			t.Errorf("expected owner %v, got %v", ownerID, store.OwnerID)
		}
	})

	t.Run("should return error if store name is empty", func(t *testing.T) {
		ownerID := uuid.New()
		_, err := NewStore("", ownerID)

		if err == nil {
			t.Error("expected error for empty name, got nil")
		}
	})
}
