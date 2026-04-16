package usecase

import (
	"context"

	paymentpb "assignment/payment-service/proto/paymentpb"

	"google.golang.org/grpc"
)

type PaymentGRPC struct {
	client paymentpb.PaymentServiceClient
}

func NewPaymentGRPC(conn *grpc.ClientConn) *PaymentGRPC {
	return &PaymentGRPC{
		client: paymentpb.NewPaymentServiceClient(conn),
	}
}

func (p *PaymentGRPC) ProcessPayment(ctx context.Context, orderID string, amount int64) (string, error) {

	res, err := p.client.ProcessPayment(ctx, &paymentpb.PaymentRequest{
		OrderId: orderID,
		Amount:  amount,
	})

	if err != nil {
		return "", err
	}

	return res.Status, nil
}
