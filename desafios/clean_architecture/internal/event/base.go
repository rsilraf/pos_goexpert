package event

import "time"

type OrderEventBase struct {
	Name    string
	Payload interface{}
}

func (e *OrderEventBase) GetName() string {
	return e.Name
}

func (e *OrderEventBase) GetPayload() interface{} {
	return e.Payload
}

func (e *OrderEventBase) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *OrderEventBase) GetDateTime() time.Time {
	return time.Now()
}
