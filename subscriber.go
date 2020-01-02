package es

import (
	"encoding/json"
	"reflect"
)

// this format has to be generated constructed from EventBridge template
type jsonEvent struct {
	Type string          `json:"t"`
	Data json.RawMessage `json:"d"`
}

// Unmarshaler -
// TODO: add init methods for binding
type Unmarshaler struct {
	eventTypes map[string]reflect.Type
}

// NewUnmarshaler -
func NewUnmarshaler(events ...interface{}) *Unmarshaler {
	result := &Unmarshaler{
		eventTypes: map[string]reflect.Type{},
	}
	result.bind(events...)
	return result
}

func (unmarshaler *Unmarshaler) bind(events ...interface{}) {
	for _, event := range events {
		eventType, t := eventType(event)
		unmarshaler.eventTypes[eventType] = t
	}
}

// Unmarshal -
func (unmarshaler Unmarshaler) Unmarshal(data []byte) (interface{}, error) {
	wrapper := jsonEvent{}
	err := json.Unmarshal(data, &wrapper)
	if err != nil {
		return nil, ErrUnmarshalEvent
	}

	t, ok := unmarshaler.eventTypes[wrapper.Type]
	if !ok {
		return nil, ErrUnknownEventType
	}

	v := reflect.New(t).Interface()
	err = json.Unmarshal(wrapper.Data, v)
	if err != nil {
		return nil, ErrUnmarshalEvent
	}

	return v.(interface{}), nil
}
