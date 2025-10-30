package service

import (
	"context"
	pb "github.com/Asylann/orderservicegrpc/proto"
	"github.com/Asylann/orderservicegrpc/server/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	OrderStore *repository.OrderStore
}

func (server *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	deliveredAt, OrderId, err := server.OrderStore.CreateOrder(int(req.UserId), int(req.CartId), req.TypeOfTransportation, req.Address)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Printf("Order was create by userId=%s", req.UserId)
	return &pb.CreateOrderResponse{DeliveredAt: timestamppb.New(deliveredAt), OrderId: int64(OrderId)}, nil
}

func (server *Server) GetOrdersByUserId(ctx context.Context, req *pb.GetOrdersByUserIdRequest) (*pb.GetOrdersByUserIdResponse, error) {
	log.Println("Was called!")
	orders, err := server.OrderStore.GetOrderByUserId(int(req.UserId))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var ValidOrders []*pb.Order
	for _, v := range orders {
		ValidOrders = append(ValidOrders, &pb.Order{
			UserId:        int32(v.UserId),
			CartId:        int32(v.CartId),
			Id:            int32(v.Id),
			Address:       v.Address,
			TransportType: v.TransportType,
			CreateAt:      timestamppb.New(v.CreateAt),
			DeliveredAt:   timestamppb.New(v.DeliveredAt),
		})
	}
	log.Printf("Orders by id=%v was sent", orders)
	return &pb.GetOrdersByUserIdResponse{Order: ValidOrders}, nil
}

func (server *Server) GetItemsOfOrderById(ctx context.Context, req *pb.GetItemsOfOrderByIdRequest) (*pb.GetItemsOfOrderByIdResponse, error) {
	products_ids, err := server.OrderStore.GetItemsOfOrderById(int(req.OrderId))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var ValidProducts_ids []int32
	for _, v := range products_ids {
		ValidProducts_ids = append(ValidProducts_ids, int32(v))
	}
	return &pb.GetItemsOfOrderByIdResponse{ListOfProductsId: ValidProducts_ids}, nil
}
