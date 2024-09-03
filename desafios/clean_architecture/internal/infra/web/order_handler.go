package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rsilraf/pos_goexpert/desafios/clean_architecture/internal/entity"
	"github.com/rsilraf/pos_goexpert/desafios/clean_architecture/internal/event"
	"github.com/rsilraf/pos_goexpert/desafios/clean_architecture/internal/usecase"
	"github.com/rsilraf/pos_goexpert/desafios/clean_architecture/pkg/events"
)

type WebOrderHandler struct {
	// USE CASES
	uc struct {
		listOrders  *usecase.ListOrdersUseCase
		createOrder *usecase.CreateOrderUseCase
	}
}

func NewWebOrderHandler(
	EventDispatcher events.EventDispatcherInterface,
	OrderRepository entity.OrderRepositoryInterface,
) *WebOrderHandler {
	return &WebOrderHandler{
		uc: struct {
			listOrders  *usecase.ListOrdersUseCase
			createOrder *usecase.CreateOrderUseCase
		}{
			listOrders: usecase.NewListOrdersUseCase(
				OrderRepository,
				event.NewOrdersListed(),
				EventDispatcher,
			),
			createOrder: usecase.NewCreateOrderUseCase(
				OrderRepository,
				event.NewOrderCreated(),
				EventDispatcher,
			),
		},
	}
}

func (h *WebOrderHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	orders, err := h.uc.listOrders.Execute()

	if fail(err, w, http.StatusInternalServerError, "error executing UC list orders") {
		return
	}

	err = json.NewEncoder(w).Encode(orders)

	if fail(err, w, http.StatusInternalServerError, "error encoding list orders output") {
		return
	}
}

func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)

	if fail(err, w, http.StatusBadRequest, "error decoding input") {
		return
	}

	output, err := h.uc.createOrder.Execute(dto)
	if fail(err, w, http.StatusInternalServerError, "error executing UC create order") {
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if fail(err, w, http.StatusInternalServerError, "error encoding create order output") {
		return
	}
}

func fail(err error, w http.ResponseWriter, status int, msg string) bool {
	if err != nil {
		log.Println(msg)
		http.Error(w, err.Error(), status)
	}
	return err != nil
}
