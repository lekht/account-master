package mock

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/lekht/account-master/src/internal/model"
)

var (
	ERROR_NO_USER_ID  = fmt.Errorf("there is no user with this id")
	ERROR_USER_EXISTS = fmt.Errorf("this user already exists")
)

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

// TODO: dont return pass
func (m *Mock) Users() ([]model.Profile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	users := make([]model.Profile, 0, len(m.users))
	for _, usr := range m.users {
		users = append(users, usr)
	}

	return users, nil
}

func (m *Mock) CreateUser(p model.Profile) (int, error) {
	m.mu.RLock()

	var exists bool
	for _, usr := range m.users {
		if p.Username == usr.Username {
			exists = true
			break
		}
	}

	var currentId int = m.nextId

	m.mu.RUnlock()

	if exists {
		return -1, ERROR_USER_EXISTS
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	p.Id = currentId
	m.users[currentId] = p

	m.nextId++

	return currentId, nil
}

func (m *Mock) UserByID(id int) (model.Profile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.users[id]
	if !exists {
		return model.Profile{}, ERROR_NO_USER_ID
	}

	return user, nil
}

func (m *Mock) UpdateUser(id int, p model.Profile) error {
	m.mu.RLock()
	usr, exists := m.users[id]
	if !exists {
		return ERROR_NO_USER_ID
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

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
	m.mu.RLock()
	_, exists := m.users[id]
	if !exists {
		return ERROR_NO_USER_ID
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.users, id)

	return nil
}
