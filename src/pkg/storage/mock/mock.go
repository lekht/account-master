package mock

import "github.com/lekht/account-master/src/internal/model"

type Mock struct {
	users map[string]model.Profile
}
