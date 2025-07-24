package controllers

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lekht/account-master/src/internal/model"
)

var (
	ErrNillReq     = errors.New("nill request")
	ErrNillProfile = errors.New("nil profile")
)

type AccountRequest struct {
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Admin    *bool   `json:"admin"`
}

type AccountResponse struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Admin    bool      `json:"admin"`
}

func requestToProfile(req *AccountRequest) (*model.Profile, error) {
	if req == nil {
		return nil, ErrNillReq
	}

	p := model.Profile{}

	if req.Email != nil {
		p.Email = *req.Email
	}

	if req.Username != nil {
		p.Username = *req.Username
	}

	if req.Password != nil {
		p.Password = *req.Password
	}

	if req.Admin != nil {
		p.Admin = *req.Admin
	}

	return &p, nil
}

func profileToResponse(p *model.Profile) (*AccountResponse, error) {
	if p == nil {
		return nil, ErrNillProfile
	}
	a := AccountResponse{}

	a.Id = p.Id
	a.Email = p.Email
	a.Username = p.Username
	a.Admin = p.Admin

	return &a, nil
}
