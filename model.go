package es

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	InitializeEvent(user, org string)
	User() string
	Organization() string
}

// Model -
type Model struct {
	EventID        string    `json:"_id"`
	Version        int       `json:"_v"`
	At             time.Time `json:"_at"`
	UserID         string    `json:"_u"`
	OrganizationID string    `json:"_org"`
}

func (model *Model) InitializeEvent(user, org string) {
	model = &Model{
		EventID:        uuid.New().String(),
		At:             time.Now().UTC(),
		Version:        0,
		UserID:         user,
		OrganizationID: org,
	}
}

func (model *Model) User() string {
	//return model.UserID
	return "test-user"
}

func (model *Model) Organization() string {
	//return model.OrganizationID
	return "test-org"
}
