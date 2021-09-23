package transactionx

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"orchestrator/internal/redisx"
	"orchestrator/internal/transactioncache"
	"orchestrator/pkg/constx"
)

type ServiceI interface {
	GetTxByCorrelationID(correlationID string, txRandomID string) (*sql.Tx, error)
	BeginTx(correlationID string) (bool, string, error)
	Commit(correlationID string, txRandomID string) error
	Rollback(correlationID string, txRandomID string) error
}

func NewTransactionSrv(db *sql.DB, redisService redisx.ServiceI, txCacheSrv transactioncache.ServiceI) *Service {
	return &Service{
		DB:    db,
		txCacheSrv: txCacheSrv,
		redisService: redisService,
	}
}

func (t *Service) GetTxByCorrelationID(correlationID string, txRandomID string) (*sql.Tx, error) {
	if val, ok := t.txCacheSrv.Get(correlationID+"_"+txRandomID); ok {
		return val, nil
	}
	return nil, errors.New("Couldn't found the transaction")
}

func (t *Service) BeginTx(correlationID string) (bool, string, error) {
	isRenew := true
	txRandomID := uuid.New().String()

	//Checking correlationID in redis and local cache
	//_, ok := t.txCacheSrv.Get(correlationID+"_"+txRandomID)
	//if  ok {
	//	isRenew = false
	//	return isRenew, nil
	//}
	//_, err := t.redisService.Get(context.Background(), correlationID+"_1234")
	//if  err == nil {
	//	isRenew = false
	//	return isRenew, nil
	//}
	// Doesn't have create new transaction
	tx, err := t.DB.Begin()
	if err != nil {
		return isRenew, txRandomID, err
	}
	//Set to local cache and redis cache
	t.txCacheSrv.Set(correlationID+"_"+txRandomID, tx)
	err = t.redisService.Set(context.Background(), correlationID+"_"+txRandomID, constx.False,  constx.ExpiredMinutes)
	if err != nil {
		return isRenew, txRandomID, err
	}
	go t.redisService.SubscribeTransactionActions(context.Background(), correlationID, txRandomID, tx)

	return  isRenew, txRandomID, nil
}

func (t *Service) Commit(correlationID string, txRandomID string) error {
	topicKey := correlationID+"_"+txRandomID
	return t.redisService.PublishTransactionActions(context.Background(), topicKey, constx.Commit)
}

func (t *Service) Rollback(correlationID string, txRandomID string) error {
	topicKey := correlationID+"_"+txRandomID
	return t.redisService.PublishTransactionActions(context.Background(), topicKey,  constx.RollBack)
}
