package orderx

import (
	"orchestrator/internal/transactionx"
)

type (
	GService struct{
		repos RepositoryI
	}

	Repository struct{
		txSrv transactionx.ServiceI
	}

	Header struct {
		CorrelationID string `json:"correlationID"`
	}

	Order struct {
		ID         *string `json:"id"`
		CustomerID string  `json:"customerID"`
		TotalPrice string  `json:"totalPrice"`
	}

	OrderRequest struct {
		Header Header `json:"header"`
		Body   Order  `json:"body"`
	}

	CustomerRequest struct {
		Header Header              `json:"header"`
		Body   CustomerRequestBody `json:"body"`
	}

	CustomerRequestBody struct {
		ID             *string `json:"id"`
		LastTotalPayed string  `json:"totalPayed"`
	}
)

