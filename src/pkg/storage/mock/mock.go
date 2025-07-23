package mock

import (
	"errors"
	"log"
	"reflect"
	"sort"
	"sync"

	"github.com/lekht/account-master/src/internal/model"
)

var (
	ErrNoUserID   = errors.New("no user with this id")
	ErrUserExists = errors.New("user already exists")
	ErrNoUsername = errors.New("no user with this username")
)

// TODO: add map[username]id
type Mock struct {
	users  map[int]model.Profile
	nextId int

	mu sync.RWMutex
}

func New() (*Mock, error) {
	m := Mock{
		users:  make(map[int]model.Profile),
		nextId: 0,
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

	sort.Slice(users, func(i, j int) bool {
		return users[i].Id < users[j].Id
	})

	return users, nil
}

func (m *Mock) CreateUser(p model.Profile) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, usr := range m.users {
		if p.Username == usr.Username {
			var currentId int = m.nextId

			p.Id = currentId
			m.users[currentId] = p

			m.nextId++

			return currentId, nil
		}
	}

	return 0, ErrUserExists
}

func (m *Mock) UserByID(id int) (model.Profile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.users[id]
	if !exists {
		return model.Profile{}, ErrNoUserID
	}

	return user, nil
}

func (m *Mock) UpdateUser(id int, p model.Profile) error {
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

func (m *Mock) DeleteUser(id int) error {
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
