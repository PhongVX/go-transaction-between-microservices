package transactionx

import (
	"context"
	"database/sql"
	"errors"
	"orchestrator/internal/redisx"
	"orchestrator/internal/transactioncache"
	"orchestrator/pkg/constx"
	"time"
)

type ServiceI interface {
	GetTxByCorrelationID(correlationID string) (*sql.Tx, error)
	BeginTx(correlationID string) (bool, error)
	Commit(correlationID string) error
	Rollback(correlationID string) error
}

func NewTransactionSrv(db *sql.DB, redisService redisx.ServiceI, txCacheSrv transactioncache.ServiceI) *Service {
	return &Service{
		DB:    db,
		txCacheSrv: txCacheSrv,
		redisService: redisService,
	}
}

func (t *Service) GetTxByCorrelationID(correlationID string) (*sql.Tx, error) {
	if val, ok := t.txCacheSrv.Get(correlationID); ok {
		return val, nil
	}
	return nil, errors.New("Couldn't found the transaction")
}

func (t *Service) BeginTx(correlationID string) (bool, error) {
	isRenew := true
	//Checking correlationID in redis and local cache
	_, ok := t.txCacheSrv.Get(correlationID)
	if  ok {
		isRenew = false
		return isRenew, nil
	}
	_, err := t.redisService.Get(context.Background(), correlationID)
	if  err == nil {
		isRenew = false
		return isRenew, nil
	}
	// Doesn't have create new transaction
	tx, err := t.DB.Begin()
	if err != nil {
		return isRenew, err
	}
	//Set to local cache and redis cache
	t.txCacheSrv.Set(correlationID, tx)
	t.redisService.Set(context.Background(), correlationID, "",  time.Minute * 10)

	go t.redisService.SubscribeTransactionActions(context.Background(), correlationID, tx)

	return  isRenew, nil
}

func (t *Service) Commit(correlationID string) error {
	return t.redisService.PublishTransactionActions(context.Background(), correlationID, constx.Commit)
}

func (t *Service) Rollback(correlationID string) error {
	return t.redisService.PublishTransactionActions(context.Background(), correlationID, constx.RollBack)
}
