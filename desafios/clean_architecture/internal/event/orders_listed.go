package event

type OrdersListed struct{ OrderEventBase }

func NewOrdersListed() *OrdersListed {
	return &OrdersListed{
		OrderEventBase{
			Name: "OrdersListed",
		},
	}
}
