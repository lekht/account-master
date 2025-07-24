package mock

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/lekht/account-master/src/internal/model"
)

func TestMock_Users(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*Mock)
		want      []model.Profile
		wantErr   bool
	}{
		{
			name: "empty users",
			mockSetup: func(m *Mock) {
			},
			want: []model.Profile{},
		},
		{
			name: "one user",
			mockSetup: func(m *Mock) {
				id := uuid.New()
				m.users[id] = model.Profile{Id: id, Username: "test"}
			},
			want: []model.Profile{{Id: uuid.New(), Username: "test"}},
		},
		{
			name: "multiple users",
			mockSetup: func(m *Mock) {
				id1 := uuid.New()
				id2 := uuid.New()
				m.users[id1] = model.Profile{Id: id1, Username: "test1", Password: "ppp", Email: "test@com", Admin: false}
				m.users[id2] = model.Profile{Id: id2, Username: "test2", Password: "ppp", Email: "test@com", Admin: false}
			},
			want: []model.Profile{{Id: uuid.New(), Username: "test1"}, {Id: uuid.New(), Username: "test2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			if tt.mockSetup != nil {
				tt.mockSetup(m)
			}

			gotUsers, err := m.Users()
			if (err != nil) != tt.wantErr {
				t.Errorf("Users() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(gotUsers) != len(tt.want) {
				t.Errorf("Users() got %d users, want %d", len(gotUsers), len(tt.want))
			}
		})
	}
}

func TestMock_CreateUser(t *testing.T) {
	id := uuid.New()
	tests := []struct {
		name        string
		mockSetup   func(*Mock)
		user        model.Profile
		wantErr     bool
		errExpected error
	}{
		{
			name:    "create user",
			user:    model.Profile{Id: id, Username: "test"},
			wantErr: false,
		},
		{
			name: "create user with existing username",
			mockSetup: func(m *Mock) {
				m.users[id] = model.Profile{Id: id, Username: "test"}
			},
			user:        model.Profile{Id: id, Username: "test"},
			errExpected: ErrUserExists,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			if tt.mockSetup != nil {
				tt.mockSetup(m)
			}

			err := m.CreateUser(tt.user)
			if (err != nil) != tt.wantErr {
				if err != nil {
					t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				} else {
					t.Errorf("CreateUser() no error, but it was expected\n")
				}
			}

			if err != nil {
				if !errors.Is(err, tt.errExpected) {
					t.Errorf("Mock.CreateUser() not expected error\n")
				}
			}
		})
	}
}

func TestMock_UserByID(t *testing.T) {
	id := uuid.New()
	user := model.Profile{Id: id, Username: "test"}

	tests := []struct {
		name      string
		id        uuid.UUID
		mockSetup func(*Mock)
		want      model.Profile
		wantErr   bool
	}{
		{
			name: "user exists",
			id:   id,
			mockSetup: func(m *Mock) {
				m.users[id] = user
			},
			want:    user,
			wantErr: false,
		},
		{
			name:    "user does not exist",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		m := New()
		if tt.mockSetup != nil {
			tt.mockSetup(m)
		}

		got, err := m.UserByID(tt.id)
		if (err != nil) != tt.wantErr {
			t.Errorf("Mock.UserByID() error = %v, wantErr %v", err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Mock.UserByID() = %v, want %v", got, tt.want)
		}
	}
}

func TestMock_UpdateUser(t *testing.T) {
	id := uuid.New()
	existingUser := model.Profile{Id: id, Username: "test", Email: "old@example.com", Admin: false, Password: "oldPassword"}
	m := New()
	m.users[id] = existingUser

	tests := []struct {
		name        string
		id          uuid.UUID
		profile     model.Profile
		wantErr     bool
		errExpected error
		wantUser    model.Profile
	}{
		{
			name:        "update email",
			id:          id,
			profile:     model.Profile{Email: "new@example.com"},
			wantErr:     false,
			errExpected: nil,
			wantUser:    model.Profile{Id: id, Username: "test", Email: "new@example.com", Admin: false, Password: "oldPassword"},
		},
		{
			name:        "update username",
			id:          id,
			profile:     model.Profile{Username: "new_username"},
			wantErr:     false,
			errExpected: nil,
			wantUser:    model.Profile{Id: id, Username: "new_username", Email: "old@example.com", Admin: false, Password: "oldPassword"},
		},
		{
			name:        "update password",
			id:          id,
			profile:     model.Profile{Password: "new_password"},
			wantErr:     false,
			errExpected: nil,
			wantUser:    model.Profile{Id: id, Username: "test", Email: "old@example.com", Admin: false, Password: "new_password"},
		},
		{
			name:        "update admin status",
			id:          id,
			profile:     model.Profile{Admin: true},
			wantErr:     false,
			errExpected: nil,
			wantUser:    model.Profile{Id: id, Username: "test", Email: "old@example.com", Admin: true, Password: "oldPassword"},
		},
		{
			name:        "user does not exist",
			id:          uuid.New(),
			profile:     model.Profile{Email: "new@example.com"},
			errExpected: ErrNoUserID,

			wantErr: true,
		},
	}
	for _, tt := range tests {

		err := m.UpdateUser(tt.id, tt.profile)
		if (err != nil) != tt.wantErr {
			if err != nil {
				t.Errorf("Mock.UpdateUser() error occured, but not expected\n")
				continue
			} else {
				t.Errorf("Mock.UpdateUser() no error, but it was expected\n")
			}
		}

		if err != nil {
			if !errors.Is(err, tt.errExpected) {
				t.Errorf("Mock.UpdateUser() not expected error\n")
			}
		}
	}
}

func TestMock_DeleteUser(t *testing.T) {
	id := uuid.New()
	m := New()
	m.users[id] = model.Profile{Id: id, Username: "test"}

	tests := []struct {
		name        string
		id          uuid.UUID
		wantErr     bool
		errExpected error
		setup       func(*Mock, uuid.UUID)
	}{
		{
			name:    "delete user",
			id:      id,
			wantErr: false,
		},
		{
			name:        "user does not exist",
			id:          uuid.New(),
			wantErr:     true,
			errExpected: ErrNoUserID,
		},
	}

	for _, tt := range tests {
		err := m.DeleteUser(tt.id)
		if (err != nil) != tt.wantErr {
			t.Errorf("Mock.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			continue
		}

		if (err != nil) != tt.wantErr {
			if err != nil {
				t.Errorf("Mock.DeleteUser() error occured, but not expected\n")
				continue
			} else {
				t.Errorf("Mock.DeleteUser() no error, but it was expected\n")
			}
		}

		if err != nil {
			if !errors.Is(err, tt.errExpected) {
				t.Errorf("Mock.DeleteUser() not expected error\n")
			}
		}
	}
}

func TestMock_UserByName(t *testing.T) {
	name := "test"
	id := uuid.New()
	user := model.Profile{Id: id, Username: name}

	tests := []struct {
		name        string
		userName    string
		mockSetup   func(*Mock)
		want        model.Profile
		wantErr     bool
		errExpected error
	}{
		{
			name:     "user exists",
			userName: name,
			mockSetup: func(m *Mock) {
				m.users[id] = user
			},
			want:    user,
			wantErr: false,
		},
		{
			name:        "user does not exist",
			userName:    "nonexistent",
			wantErr:     true,
			errExpected: ErrNoUsername,
		},
	}

	for _, tt := range tests {
		m := New()
		if tt.mockSetup != nil {
			tt.mockSetup(m)
		}
		got, err := m.UserByName(tt.userName)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Mock.UserByName() = %v, want %v", got, tt.want)
		}

		if (err != nil) != tt.wantErr {
			t.Errorf("Mock.UserByName() error = %v, wantErr %v", err, tt.wantErr)
			continue
		}

		if (err != nil) != tt.wantErr {
			if err != nil {
				t.Errorf("Mock.UserByName() error occured, but not expected\n")
				continue
			} else {
				t.Errorf("Mock.UserByName() no error, but it was expected\n")
			}
		}

		if err != nil {
			if !errors.Is(err, tt.errExpected) {
				t.Errorf("Mock.UserByName() not expected error\n")
			}
		}
	}
}
