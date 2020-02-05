package main

import (
	"encoding/json"
	"fmt"

	"github.com/mrzahrada/es"
)

type Event struct {
	*es.Model
	Msg string `json:"msg"`
}

func main() {

	event := Event{
		Msg: "hello world",
		Model: es.NewModel(es.Identity{
			UserID:    "user-1",
			UserOrgID: "org-1",
			UserRole:  "role-1",
		}),
	}

	data, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", string(data))

}
