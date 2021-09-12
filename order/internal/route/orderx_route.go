package route

import (
	"github.com/PhongVX/micro-protos/order"
	"github.com/PhongVX/micro-protos/transaction"
	"order/internal/orderx"
)

func newOrderXHandler(orderC order.OrderClient, txC transaction.TransactionClient) (*orderx.Handler, error) {
	srv := orderx.NewService(orderC, txC)
	handler := orderx.NewHandler(&srv)
	return handler, nil
}
