package orderx

import (
	"context"
	"github.com/PhongVX/micro-protos/order"
	"github.com/google/uuid"
	"orchestrator/internal/transactionx"
)

type RepositoryI interface {
	InsertOrder(ctx context.Context, in *order.InsertOrderRequest) (*order.InsertOrderResponse, error)
	InsertOrderDetail(context.Context, *order.InsertOrderDetailRequest) (*order.InsertOrderDetailResponse, error)
}

func NewRepository(txSrv transactionx.ServiceI) *Repository{
	return &Repository{
		txSrv: txSrv,
	}
}

func (r *Repository)InsertOrder(ctx context.Context, o *order.InsertOrderRequest) (*order.InsertOrderResponse, error) {
	tx, err := r.txSrv.GetTxByCorrelationID(o.CorrelationID)
	if err != nil {
		return nil, err
	}
	id := uuid.New().String()
	queryString := `INSERT INTO orders(id, phone_number, address, name, total_price) VALUES($1, $2, $3, $4, $5)`
	rs, err := tx.Exec(queryString, id, o.PhoneNumber, o.Address, o.Name, o.TotalPrice)
	if val, err := rs.RowsAffected(); err == nil && val > 0 {
		return &order.InsertOrderResponse{Id: id}, nil
	}
	return nil, err
}

func (r *Repository) InsertOrderDetail(ctx context.Context, o *order.InsertOrderDetailRequest) (*order.InsertOrderDetailResponse, error){
	tx, err := r.txSrv.GetTxByCorrelationID(o.CorrelationID)
	if err != nil {
		return nil, err
	}
	numRowAffected := int64(0)
	for _, od := range o.OrderDetails {
		queryString := `INSERT INTO order_detail(order_id, product_id, quantity, price, total_price) VALUES($1, $2, $3, $4, $5)`
		rs, err := tx.Exec(queryString, od.OrderID, od.ProductID, od.Quantity, od.Price, od.TotalPrice)
		rowAffected, err := rs.RowsAffected()
		if err != nil {
			return nil, err
		}
		numRowAffected += rowAffected
	}
	return &order.InsertOrderDetailResponse{RowAffected: int32(numRowAffected)}, nil
}