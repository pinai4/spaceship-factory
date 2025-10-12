package v1

import (
	def "github.com/pinai4/spaceship-factory/order/internal/client"
	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	generatedClient inventoryV1.InventoryServiceClient
}

func NewClient(inventoryClient inventoryV1.InventoryServiceClient) *client {
	return &client{generatedClient: inventoryClient}
}
