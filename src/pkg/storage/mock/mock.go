package mock

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/lekht/account-master/src/internal/model"
)

var (
	ErrNoUserID   = errors.New("no user with this id")
	ErrUserExists = errors.New("user already exists")
	ErrNoUsername = errors.New("no user with this username")
)

type Mock struct {
	users map[uuid.UUID]model.Profile

	mu sync.RWMutex
}

func New() (*Mock, error) {
	m := Mock{
		users: make(map[uuid.UUID]model.Profile),
	}

	return &m, nil
}

func (m *Mock) Users() ([]model.Profile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	users := make([]model.Profile, 0, len(m.users))
	for _, usr := range m.users {
		users = append(users, usr)
	}

	if len(users) == 0 {
		return users, nil
	}

	return users, nil
}

func (m *Mock) CreateUser(p model.Profile) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, usr := range m.users {
		if p.Username == usr.Username {
			return ErrUserExists
		}
	}

	for {
		p.Id = uuid.New()
		if _, ok := m.users[p.Id]; !ok {
			break
		}
	}

	m.users[p.Id] = p

	return nil
}

func (m *Mock) UserByID(id uuid.UUID) (model.Profile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.users[id]
	if !exists {
		return model.Profile{}, ErrNoUserID
	}

	return user, nil
}

func (m *Mock) UpdateUser(id uuid.UUID, p model.Profile) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	usr, exists := m.users[id]
	if !exists {
		return ErrNoUserID
	}

	// updates only non default value
	if p.Email != "" {
		usr.Email = p.Email
	}

	if p.Username != "" {
		usr.Username = p.Username
	}

	if p.Password != "" {
		usr.Password = p.Password
	}

	if p.Admin != usr.Admin {
		usr.Admin = p.Admin
	}

	m.users[id] = usr

	return nil
}

func (m *Mock) DeleteUser(id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.users[id]
	if !exists {
		return ErrNoUserID
	}

	delete(m.users, id)

	return nil
}

func (m *Mock) UserByName(name string) (model.Profile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.Username == name {
			return user, nil
		}
	}

	return model.Profile{}, ErrNoUsername
}
