package v1

import (
	"context"
	"log"
	"net/http"

	orderV1 "github.com/pinai4/spaceship-factory/shared/pkg/openapi/order/v1"
)

func (a *api) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	log.Printf("api.NewError error: %v\n", err)
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString("Internal Server Error"),
		},
	}
}
