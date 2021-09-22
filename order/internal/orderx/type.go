package orderx

import (
	"github.com/PhongVX/micro-protos/order"
	"github.com/PhongVX/micro-protos/transaction"
)

type (
	Header struct {
		CorrelationID string `json:"correlationID"`
	}

	OrderDetail struct {
		ProductID int32 `json:"productID"`
		Price  float64 `json:"price"`
		TotalPrice float64 `json:"totalPrice"`
		Quantity int32 `json:"quantity"`
	}

	Order struct {
		ID         *string `json:"id"`
		PhoneNumber string  `json:"phoneNumber"`
		Name string `json:"name"`
		Address string `json:"address"`
		OrderDetails []*OrderDetail `json:"orderDetails"`
		TotalPrice float64  `json:"totalPrice"`
	}

	OrderRequest struct {
		Header Header `json:"header"`
		Body   Order  `json:"body"`
	}

	ProductRequest struct {
		Header Header `json:"header"`
		Body   Product  `json:"body"`
	}

	Product struct {
		ID int32 `json:"id"`
		Quantity int32 `json:"quantity"`
	}

	Handler struct {
		srv ServiceI
	}

	Service struct {
		gOrderC order.OrderClient
		gTransactionC transaction.TransactionClient
	}

	GService struct {
		orderC order.OrderClient
	}
)

