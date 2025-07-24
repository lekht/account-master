package mock

import (
	"errors"
	"log"
	"reflect"
	"sync"

	"github.com/google/uuid"
	"github.com/lekht/account-master/src/internal/model"
)

var (
	ErrNoUserID   = errors.New("no user with this id")
	ErrUserExists = errors.New("user already exists")
	ErrNoUsername = errors.New("no user with this username")
)

// TODO: add map[username]id (more optimized)
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
		if _, ok := m.users[p.Id]; ok {
			continue
		} else {
			break
		}
	}

	m.users[p.Id] = p
	for _, usr := range m.users {
		if p.Username == usr.Username {
			return ErrUserExists
		}
	}

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
	src := reflect.ValueOf(&p).Elem()
	dst := reflect.ValueOf(&usr).Elem()

	for i := 0; i < src.NumField(); i++ {
		srcField := src.Field(i)
		if !srcField.IsZero() {
			dstField := dst.Field(i)
			if dstField.CanSet() {
				dstField.Set(srcField)
			}
		}
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

	log.Println("==== name ", name)
	for _, user := range m.users {
		log.Printf("==== profile: %+v\n", user)
		if user.Username == name {
			return user, nil
		}
	}

	return model.Profile{}, ErrNoUsername
}
