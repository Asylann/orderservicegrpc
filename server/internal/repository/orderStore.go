package repository

import (
	"context"
	"errors"
	"github.com/Asylann/OrderServiceGRPC/server/internal/models"
	cpb "github.com/Asylann/gRPC_Demo/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

var CartClient cpb.CartServiceClient

func InitCartServiceConn() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Println("Can not connect to Cart service")
		return
	}
	CartClient = cpb.NewCartServiceClient(conn)
}

func (orderStore *OrderStore) CreateOrder(UserId int, CartId int, TransportType string, Address string) (time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	DeliveredAt, err := DefineTimeOfDelivery(TransportType)
	if err != nil {
		log.Println(err.Error())
		return DeliveredAt, err
	}
	var OrderId int
	err = orderStore.CreateOrderStmt.QueryRowxContext(ctx, CartId, time.Now(), DeliveredAt, TransportType, UserId, Address).Scan(&OrderId)
	if err != nil {
		return DeliveredAt, err
	}

	r, err := CartClient.GetItemsOfCartById(ctx, &cpb.GetItemsOfCartByIdRequest{Id: int32(CartId)})
	if err != nil {
		log.Println("Error during get cart items from CartService")
		return DeliveredAt, err
	}
	ProductsIds := r.GetProduct()
	var ValidProducts_Id []int
	for _, v := range ProductsIds {
		ValidProducts_Id = append(ValidProducts_Id, int(v.Id))
	}

	err = orderStore.PutItemsToOrder(OrderId, ValidProducts_Id)
	if err != nil {
		log.Printf("Cant put items to Order: %s \n", err.Error())
		return DeliveredAt, err
	}

	return DeliveredAt, nil
}

func (orderStore *OrderStore) GetOrderByUserId(UserId int) (models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	var Order models.Order
	if err := orderStore.GetOrderByUserIdStmt.GetContext(ctx, &Order, UserId); err != nil {
		return Order, err
	}

	return Order, nil
}

func (orderStore *OrderStore) GetItemsOfOrderById(Id int) ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	var products_id []int
	if err := orderStore.GetItemsOfOrderByIdStmt.SelectContext(ctx, &products_id, Id); err != nil {
		return nil, err
	}

	return products_id, nil
}

func (orderStore *OrderStore) PutItemsToOrder(orderId int, productIds []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, v := range productIds {
		if _, err := orderStore.PutOneItemToOrder.ExecContext(ctx, orderId, v); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	return nil
}

func DefineTimeOfDelivery(TransportType string) (time.Time, error) {
	if TransportType == "Ordinary" {
		log.Println("Ordinary was")
		return time.Now().Add(24 * time.Hour * 4), nil //4 days
	} else if TransportType == "Medium" {
		log.Println("Medium was")
		return time.Now().Add(24 * time.Hour * 2), nil //2 days
	} else if TransportType == "Fastest" {
		log.Println("fastest was")
		return time.Now().Add(24 * time.Hour), nil
	} else {
		log.Println("Bag")
		return time.Now(), errors.New("Invalid TransportType")
	}
}
