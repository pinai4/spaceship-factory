package v1

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"github.com/pinai4/spaceship-factory/payment/internal/service"
	paymentV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) *api {
	return &api{paymentService: paymentService}
}

func PrinterInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			if r, ok := resp.(*paymentV1.PayOrderResponse); ok && r != nil {
				log.Printf("Payment was successful, %s\n", r.GetTransactionUuid())
				// log.Printf("Request: %#v\n", req)
			}
		}

		return resp, err
	}
}
