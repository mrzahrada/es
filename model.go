package es

import (
	"time"

	"github.com/google/uuid"
)

type Identity struct {
	UserID    string `json:"_au"`
	UserOrgID string `json:"_ao"`
	UserRole  string `json:"_ar"`
}

// Event -
type Event interface{}

// Model -
type Model struct {
	EventID  string    `json:"_e"`
	Version  int       `json:"_v"`
	At       time.Time `json:"_t"`
	Identity Identity  `json:"_i"`
}

func NewModel(identity Identity) *Model {
	return &Model{
		EventID:  uuid.New().String(),
		Version:  0,
		At:       time.Now().UTC(),
		Identity: identity,
	}
}
