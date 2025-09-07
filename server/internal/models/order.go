package models

import "time"

type Order struct {
	Id            int       `json:"id" db:"id"`
	CartId        int       `json:"cart_id" db:"cart_id"`
	CreateAt      time.Time `json:"createAt" db:"create_at"`
	DeliveredAt   time.Time `json:"deliveredAt" db:"delivered_at"`
	TransportType string    `json:"transport_Type" db:"transport_type"`
	UserId        int       `json:"userId" db:"user_id"`
	Address       string    `json:"address" db:"address"`
}
