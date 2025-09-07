package service

import (
	"context"
	pb "github.com/Asylann/OrderServiceGRPC/proto"
	"github.com/Asylann/OrderServiceGRPC/server/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	OrderStore *repository.OrderStore
}

func (server *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	deliveredAt, err := server.OrderStore.CreateOrder(int(req.UserId), int(req.CartId), req.TypeOfTransportation, req.Address)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Printf("Order was create by userId=%s", req.UserId)
	return &pb.CreateOrderResponse{DeliveredAt: timestamppb.New(deliveredAt)}, nil
}

func (server *Server) GetOrderByUserId(ctx context.Context, req *pb.GetOrderByUserIdRequest) (*pb.GetOrderByUserIdResponse, error) {
	order, err := server.OrderStore.GetOrderByUserId(int(req.UserId))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.GetOrderByUserIdResponse{Order: &pb.Order{
		UserId:        int32(order.UserId),
		CartId:        int32(order.CartId),
		CreateAt:      timestamppb.New(order.CreateAt),
		DeliveredAt:   timestamppb.New(order.DeliveredAt),
		Id:            int32(order.Id),
		TransportType: order.TransportType,
		Address:       order.Address,
	}}, nil
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
