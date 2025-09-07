package models

import "time"

type Order struct {
	Id            int       `json:"id" sqlx:"id"`
	CartId        int       `json:"cart_Id" sqlx:"cart_id"`
	CreateAt      time.Time `json:"createAt" sqlx:"create_at"`
	DeliveredAt   time.Time `json:"deliveredAt" sqlx:"delivered_at"`
	TransportType string    `json:"transport_Type" sqlx:"transport_type"`
	UserId        int       `json:"userId" sqlx:"user_id"`
	Address       string    `json:"address" sqlx:"address"`
}
