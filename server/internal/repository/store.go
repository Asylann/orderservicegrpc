package repository

import (
	"context"
	"github.com/Asylann/orderservicegrpc/server/internal/models"
	"github.com/jmoiron/sqlx"
	"time"
)

type sqlOrderStore interface {
	CreateOrder(UserId int, CartId int, TransportType string, Address string) (time.Time, error)
	GetOrderByUserId(UserId int) (models.Order, error)
	GetItemsOfOrderById(Id int) ([]models.OrderItem, error)
}

type OrderStore struct {
	CreateOrderStmt         *sqlx.Stmt
	GetOrderByUserIdStmt    *sqlx.Stmt
	GetItemsOfOrderByIdStmt *sqlx.Stmt
	PutOneItemToOrder       *sqlx.Stmt
}

func NewOrderStore() (*OrderStore, error) {
	orderStore := OrderStore{}
	var err error
	orderStore.CreateOrderStmt, err = DB.PreparexContext(context.Background(),
		`INSERT INTO orders(cart_id, create_at,delivered_at, transport_type, user_id, address) VALUES($1,$2,$3,$4,$5,$6) RETURNING id`)
	if err != nil {
		return &OrderStore{}, err
	}
	orderStore.GetOrderByUserIdStmt, err = DB.PreparexContext(context.Background(),
		`SELECT * FROM orders WHERE user_id=$1`)
	if err != nil {
		return &OrderStore{}, err
	}
	orderStore.GetItemsOfOrderByIdStmt, err = DB.PreparexContext(context.Background(),
		`SELECT product_id FROM order_items WHERE order_id=$1`)
	if err != nil {
		return &OrderStore{}, err
	}
	orderStore.PutOneItemToOrder, err = DB.PreparexContext(context.Background(),
		`INSERT INTO order_items(order_id, product_id) VALUES($1,$2)`)
	if err != nil {
		return &OrderStore{}, err
	}
	return &orderStore, nil
}
