package transactionx

import (
	"context"
	"database/sql"
	"github.com/PhongVX/micro-protos/transaction"
)

type (
	Service struct {
		MapTx map[string]*sql.Tx
		DB *sql.DB
	}

	GService struct{
		txSrv ServiceI
	}

	ServiceI interface {
		GetTxByCorrelationID(correlationID string) (*sql.Tx, error)
		BeginTx(correlationID string) (*sql.Tx, bool, error)
		Commit(correlationID string) error
		Rollback(correlationID string) error
	}

	GServiceI interface {
		BeginTx(context.Context, *transaction.BeginTxRequest) (*transaction.BeginTxResponse, error)
		Commit(context.Context, *transaction.CommonTxDoActionRequest) (*transaction.CommonTxResponse, error)
		Rollback(context.Context, *transaction.CommonTxDoActionRequest) (*transaction.CommonTxResponse, error)
	}
)
