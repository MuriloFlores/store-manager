package entity

import (
	"errors"

	"github.com/MuriloFlores/order-manager/internal/organization/domain/vo"
	"github.com/google/uuid"
)

type Store struct {
	Name       string
	SchemaName vo.SchemaName
	OwnerID    uuid.UUID
	Status     vo.StoreStatus
}

func NewStore(storeName string, ownerID uuid.UUID) (*Store, error) {
	if storeName == "" {
		return nil, errors.New("store name is required")
	}

	if ownerID == uuid.Nil {
		return nil, errors.New("store owner is required")
	}

	storeSchema, err := vo.NewSchemaName(storeName)
	if err != nil {
		return nil, err
	}

	return &Store{
		Name:       storeName,
		SchemaName: storeSchema,
		OwnerID:    ownerID,
		Status:     vo.StatusPending,
	}, nil
}

func (s *Store) Activate() error {
	if s.Status != vo.StatusPending {
		return errors.New("only pending store can activate")
	}

	s.Status = vo.StatusActive
	return nil
}

func (s *Store) Deactivate() error {
	if s.Status == vo.StatusDeactivated {
		return errors.New("store is already deactivated")
	}

	s.Status = vo.StatusDeactivated
	return nil
}

func (s *Store) Fail() error {
	if s.Status != vo.StatusPending {
		return errors.New("only pending store can activate")
	}

	s.Status = vo.StatusFailed
	return nil
}

func (s *Store) ChangeStoreName(newName string) error {
	if newName == "" {
		return errors.New("new name can't be empty")
	}

	s.Name = newName
	return nil
}
