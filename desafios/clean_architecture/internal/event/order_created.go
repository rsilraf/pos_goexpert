package event

type OrderCreated struct{ OrderEventBase }

func NewOrderCreated() *OrderCreated {
	return &OrderCreated{
		OrderEventBase{
			Name: "OrderCreated",
		},
	}
}
