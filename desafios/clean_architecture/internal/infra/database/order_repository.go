package database

import (
	"context"
	"database/sql"

	"github.com/rsilraf/pos_goexpert/desafios/clean_architecture/internal/entity"
	db "github.com/rsilraf/pos_goexpert/desafios/clean_architecture/internal/infra/sqlc"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) ListAll() ([]*entity.Order, error) {
	ctx := context.Background()
	db := db.New(r.Db)
	dbOrders, err := db.ListOrders(ctx)
	if err != nil {
		return nil, err
	}

	// converte de db.Order para entity.Order
	orders := []*entity.Order{}
	for _, o := range dbOrders {

		orders = append(orders, &entity.Order{
			ID:         o.ID,
			Price:      o.Price,
			Tax:        o.Tax,
			FinalPrice: o.FinalPrice,
		})
	}
	return orders, nil
}
