package productx

import (
	"github.com/PhongVX/micro-protos/product"
	"github.com/PhongVX/micro-protos/transaction"
)

type (
	Header struct {
		CorrelationID string `json:"correlationID"`
	}

	Product struct {
		ID         *int `json:"id"`
		Quantity int32  `json:"quantity"`
	}

	UpdateProductRequest struct {
		Header Header `json:"header"`
		Body   Product  `json:"body"`
	}

	Handler struct {
		srv ServiceI
	}

	Service struct {
		gProductC product.ProductClient
		gTransactionC transaction.TransactionClient
	}

	GService struct {
		productC product.ProductClient
	}
)

