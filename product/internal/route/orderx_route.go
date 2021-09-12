package route

import (
	"github.com/PhongVX/micro-protos/product"
	"github.com/PhongVX/micro-protos/transaction"
	"product/internal/productx"
)

func newOrderXHandler(productC product.ProductClient, txC transaction.TransactionClient) (*productx.Handler, error) {
	srv := productx.NewService(productC, txC)
	handler := productx.NewHandler(&srv)
	return handler, nil
}
