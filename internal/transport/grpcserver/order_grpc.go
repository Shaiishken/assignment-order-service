package grpcserver

import (
	"context"
	"time"

	"assignment/internal/usecase"
	orderpb "assignment/proto/orderpb"
)

type OrderGRPCServer struct {
	orderpb.UnimplementedOrderServiceServer
	uc *usecase.OrderUsecase
}

func NewOrderGRPCServer(uc *usecase.OrderUsecase) *OrderGRPCServer {
	return &OrderGRPCServer{uc: uc}
}

func (s *OrderGRPCServer) SubscribeToOrderUpdates(
	req *orderpb.OrderRequest,
	stream orderpb.OrderService_SubscribeToOrderUpdatesServer,
) error {
	orderID := req.OrderId
	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		default:
		}

		order, err := s.uc.GetOrder(context.Background(), orderID)
		if err != nil {
			return err
		}

		if order != nil {
			err = stream.Send(&orderpb.OrderStatusUpdate{
				OrderId: order.ID,
				Status:  order.Status,
			})
			if err != nil {
				return err
			}

			if order.Status == "paid" || order.Status == "failed" || order.Status == "cancelled" {
				return nil
			}
		}

		time.Sleep(2 * time.Second)
	}
}
