package es

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Identity struct {
	UserID    string `json:"_au"`
	UserOrgID string `json:"_ao"`
	UserRole  string `json:"_ar"`
}

type Event interface {
	EventID() string
	At() time.Time
	Version() int
	Identity() Identity
}

// Model -
type Model struct {
	eventID  string
	version  int
	at       time.Time
	identity Identity
}

func NewModel(identity Identity) *Model {
	return &Model{
		eventID:  uuid.New().String(),
		version:  0,
		at:       time.Now().UTC(),
		identity: identity,
	}
}

func (model *Model) MarshalJSON() ([]byte, error) {

	return json.Marshal(struct {
		Model interface{} `json:"_m"`
	}{
		Model: struct {
			EventID   string `json:"e"`
			Version   int    `json:"v"`
			At        int64  `json:"t"`
			UserID    string `json:"u"`
			UserOrgID string `json:"o"`
			UserRole  string `json:"r"`
		}{
			EventID:   model.eventID,
			Version:   model.version,
			At:        model.at.Unix(),
			UserID:    model.identity.UserID,
			UserOrgID: model.identity.UserOrgID,
			UserRole:  model.identity.UserRole,
		},
	})
}

func (model *Model) UnmarshalJSON(data []byte) error {

	type input struct {
		Model struct {
			EventID   string `json:"e"`
			Version   int    `json:"v"`
			At        int64  `json:"t"`
			UserID    string `json:"u"`
			UserOrgID string `json:"o"`
			UserRole  string `json:"r"`
		} `json:"_m"`
	}
	x := input{}
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	model.eventID = x.Model.EventID
	model.version = x.Model.Version
	model.at = time.Unix(x.Model.At, 0)
	model.identity = Identity{
		UserID:    x.Model.UserID,
		UserOrgID: x.Model.UserOrgID,
		UserRole:  x.Model.UserRole,
	}

	return nil
}

func (model Model) EventID() string {
	return model.eventID
}

func (model Model) At() time.Time {
	return model.at
}
func (model Model) Version() int {
	return model.version
}
func (model Model) Identity() Identity {
	return model.identity
}
