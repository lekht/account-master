package storage

import (
	"fmt"
	"log"

	"github.com/lekht/account-master/src/internal/model"
	"github.com/lekht/account-master/src/pkg/storage/mock"
)

type DB interface {
	Users() ([]model.Profile, error)
	UserByID(int) (model.Profile, error)
	CreateUser(model.Profile) (int, error)
	UpdateUser(int, model.Profile) error
	DeleteUser(int) error
}

type Storage struct {
	db DB
}

func New() (*Storage, error) {
	m, err := mock.New()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}

	s := Storage{
		db: m,
	}

	log.Println("New storage is created")

	return &s, nil
}

func (s *Storage) Users() ([]model.Profile, error) {
	users, err := s.db.Users()
	if err != nil {
		return nil, fmt.Errorf("failed to get users list from storage: %v", err)
	}

	return users, nil
}

func (s *Storage) CreateUser(p model.Profile) (int, error) {
	id, err := s.db.CreateUser(p)
	if err != nil {
		return -1, fmt.Errorf("failed to create new user in storage: %v", err)
	}

	return id, nil
}

func (s *Storage) UserByID(id int) (model.Profile, error) {
	usr, err := s.db.UserByID(id)
	if err != nil {
		return usr, fmt.Errorf("failed to get user by id from storage: %v", err)
	}

	return usr, nil
}

func (s *Storage) UpdateUser(id int, p model.Profile) error {
	err := s.db.UpdateUser(id, p)
	if err != nil {
		return fmt.Errorf("failed to update user in storage: %v", err)
	}

	return nil
}

func (s *Storage) DeleteUser(id int) error {
	err := s.db.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("failed to delete user in storage: %v", err)
	}

	return nil
}
