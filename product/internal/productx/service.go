package productx

import (
	"context"

	"github.com/PhongVX/micro-protos/product"
	"github.com/PhongVX/micro-protos/transaction"
)

type (
	ServiceI interface {
		UpdateProduct(ctx context.Context, p UpdateProductRequest) (*int32, error)
	}
)

func NewService(productC product.ProductClient, txC transaction.TransactionClient) Service{
	return Service{
		gProductC: productC,
		gTransactionC: txC,
	}
}

func (s *Service) UpdateProduct(ctx context.Context, p UpdateProductRequest) (*int32, error) {
	//1. Begin Transaction
	txRequest := &transaction.BeginTxRequest{CorrelationID: p.Header.CorrelationID}
	txRes, err := s.gTransactionC.BeginTx(ctx, txRequest)
	if err != nil {
		return nil, err
	}
	productRequest := product.UpdateProductRequest{
		CorrelationID: p.Header.CorrelationID,
		ID: int32(*p.Body.ID),
		Quantity: p.Body.Quantity,
	}
	//2. Insert To Order Table
	productRes, err := s.gProductC.UpdateProduct(ctx, &productRequest)
	txInfo := &transaction.CommonTxDoActionRequest{
		CorrelationID: p.Header.CorrelationID,
		BeginTxRes: txRes,
	}
	if err != nil {
		s.gTransactionC.Rollback(ctx, txInfo)
		return nil, err
	}
	s.gTransactionC.Commit(ctx, txInfo)
	return &productRes.RowAffected, nil
}