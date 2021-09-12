package transactionx

import (
	"database/sql"
	"errors"
)

func NewTransactionSrv(db *sql.DB) *Service{
	mapTx := make(map[string]*sql.Tx)
	return &Service{
		MapTx: mapTx,
		DB: db,
	}
}

func (t *Service) GetTxByCorrelationID(correlationID string) (*sql.Tx, error) {
	if val, ok := t.MapTx[correlationID]; ok {
		return val, nil
	}
	return nil, errors.New("Couldn't found the transaction")
}

func (t *Service) BeginTx(correlationID string) (*sql.Tx, bool, error){
	isRenew := true
	if tx, ok := t.MapTx[correlationID]; ok {
		isRenew = false
		return tx, isRenew, nil
	}
	tx, err := t.DB.Begin()
	if err != nil {
		return nil, isRenew, err
	}
	t.MapTx[correlationID] = tx
	return tx, isRenew, nil
}

func (t *Service) Commit(correlationID string) error{
	if tx, ok := t.MapTx[correlationID]; ok {
		err := tx.Commit()
		if err != nil {
			return err
		}
		delete(t.MapTx, correlationID)
	}
	return nil
}

func (t *Service) Rollback(correlationID string) error{
	if tx, ok := t.MapTx[correlationID]; ok {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		delete(t.MapTx, correlationID)
	}
	return nil
}