package models

type OrderItem struct {
	Id        int `json:"id"`
	OrderId   int `json:"orderId"`
	ProductId int `json:"productId"`
}
