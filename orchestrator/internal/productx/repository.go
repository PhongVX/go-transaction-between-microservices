package productx

import (
	"context"
	"github.com/PhongVX/micro-protos/order"
	"github.com/PhongVX/micro-protos/product"
	"orchestrator/internal/transactionx"
)

type RepositoryI interface {
	UpdateProduct(context.Context, *product.UpdateProductRequest) (*product.UpdateProductResponse, error)
}

func NewRepository(txSrv transactionx.ServiceI) *Repository{
	return &Repository{
		txSrv: txSrv,
	}
}

func (r *Repository)UpdateProduct(ctx context.Context, p *product.UpdateProductRequest) (*product.UpdateProductResponse, error){
	tx, err := r.txSrv.GetTxByCorrelationID(p.CorrelationID)
	if err != nil {
		return nil, err
	}
	queryString := `UPDATE product SET quantity=(quantity - $1) WHERE id=$2`
	rs, err := tx.Exec(queryString, p.Quantity, p.ID)
	if val, err := rs.RowsAffected(); err == nil && val > 0 {
		return &product.UpdateProductResponse{RowAffected: int32(val)}, nil
	}
	return &product.UpdateProductResponse{RowAffected: int32(0)}, err
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