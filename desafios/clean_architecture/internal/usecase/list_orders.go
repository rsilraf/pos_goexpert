package usecase

import (
	"github.com/rsilraf/pos_goexpert/desafios/clean_architecture/internal/entity"
	"github.com/rsilraf/pos_goexpert/desafios/clean_architecture/pkg/events"
)

type OrderDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type ListOrdersDTO []OrderDTO

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrdersListed    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewListOrdersUseCase(
	repo entity.OrderRepositoryInterface,
	event events.EventInterface,
	dispatcher events.EventDispatcherInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: repo,
		OrdersListed:    event,
		EventDispatcher: dispatcher,
	}
}

func (lo *ListOrdersUseCase) Execute() (ListOrdersDTO, error) {
	orders, err := lo.OrderRepository.ListAll()
	if err != nil {
		return nil, err
	}
	// converte de entity.Order para DTO
	ordersListing := ListOrdersDTO{}
	for _, o := range orders {
		ordersListing = append(ordersListing, OrderDTO{
			ID:         o.ID,
			Price:      o.Price,
			Tax:        o.Tax,
			FinalPrice: o.FinalPrice,
		})
	}

	lo.OrdersListed.SetPayload(ordersListing)
	lo.EventDispatcher.Dispatch(lo.OrdersListed)

	return ordersListing, nil
}
