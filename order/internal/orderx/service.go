package orderx

import (
	"context"
	"encoding/json"
	"github.com/PhongVX/micro-protos/order"
	"github.com/PhongVX/micro-protos/transaction"
	"log"
	"order/pkg/http/request"
)

type (
	ServiceI interface {
		InsertOrder(ctx context.Context, o OrderRequest) (*string, error)
	}
)

func NewService(orderC order.OrderClient, txC transaction.TransactionClient) Service{
	return Service{
		gOrderC: orderC,
		gTransactionC: txC,
	}
}

func (s *Service) InsertOrder(ctx context.Context, o OrderRequest) (*string, error) {
	//1. Begin Transaction
	txRequest := &transaction.BeginTxRequest{CorrelationID: o.Header.CorrelationID}
	txRes, err := s.gTransactionC.BeginTx(ctx, txRequest)
	if err != nil {
		log.Println("Failed to begin transaction")
	}
	orderRequest := order.InsertOrderRequest{
		CorrelationID: o.Header.CorrelationID,
		PhoneNumber: o.Body.PhoneNumber,
		Address: o.Body.Address,
		Name :o.Body.Name,
		TotalPrice: o.Body.TotalPrice,
	}
	//2. Insert To Order Table
	orderRes, err := s.gOrderC.InsertOrder(ctx, &orderRequest)
	txInfo := &transaction.CommonTxDoActionRequest{
		CorrelationID: o.Header.CorrelationID,
		BeginTxRes: txRes,
	}
	if err != nil {
		s.gTransactionC.Rollback(ctx, txInfo)
		return nil, err
	}
	//3 Insert Into Order Detail Table
	orderDetails := make([]*order.OrderDetail, 0)
	for _, od := range o.Body.OrderDetails {
		orderDetails = append(orderDetails, &order.OrderDetail{
			ProductID: od.ProductID,
			Quantity: od.Quantity,
			OrderID: orderRes.Id,
			Price: od.Price,
			TotalPrice: od.TotalPrice,
		})
		//3 Call HTTP API Update Product Quantity
		pReq := ProductRequest{
			Header: Header{
				CorrelationID: o.Header.CorrelationID,
			},
			Body: Product{
				ID: od.ProductID,
				Quantity: od.Quantity,
			},
		}
		pBytes, err := json.Marshal(pReq)
		if err != nil {
			return nil, err
		}
		resP, err := request.Put(productAPI, pBytes)
		if err != nil {
			s.gTransactionC.Rollback(ctx, txInfo)
			return nil, err
		}
		if resP.ID == nil {
			s.gTransactionC.Rollback(ctx, txInfo)
			return nil, err
		}

	}
	//4. Insert Order Detail
	_, err = s.gOrderC.InsertOrderDetail(ctx, &order.InsertOrderDetailRequest{
		CorrelationID: o.Header.CorrelationID,
		OrderDetails: orderDetails,
	})
	if err != nil {
		s.gTransactionC.Rollback(ctx, txInfo)
		return nil, err
	}
	s.gTransactionC.Commit(ctx, txInfo)
	return &orderRes.Id, nil
}