package models

import pb "github.com/Asylann/OrderService/proto"

type Server struct {
	pb.UnimplementedOrderServiceServer
}
