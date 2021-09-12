package transactionx

import (
	"context"
	"github.com/PhongVX/micro-protos/transaction"
)

func NewGService(txSrv ServiceI) *GService {
	return &GService{
		txSrv: txSrv,
	}
}

func (s *GService) BeginTx(ctx context.Context, in *transaction.BeginTxRequest) (*transaction.BeginTxResponse, error){
	_, isRenew, err := s.txSrv.BeginTx(in.CorrelationID)
	return &transaction.BeginTxResponse{IsRenew: isRenew}, err
}

func (s *GService) Commit(ctx context.Context, in *transaction.CommonTxDoActionRequest) (*transaction.CommonTxResponse, error){
	ok := true
	if in.BeginTxRes.IsRenew {
		err := s.txSrv.Commit(in.CorrelationID)
		if err != nil {
			ok = false
			return &transaction.CommonTxResponse{Ok: ok}, err
		}
	}
	return &transaction.CommonTxResponse{Ok: ok}, nil
}

func  (s *GService) Rollback(ctx context.Context, in *transaction.CommonTxDoActionRequest) (*transaction.CommonTxResponse, error){
	ok := true
	if in.BeginTxRes.IsRenew {
		err := s.txSrv.Commit(in.CorrelationID)
		if err != nil {
			ok = false
			return &transaction.CommonTxResponse{Ok: ok}, err
		}
	}
	return &transaction.CommonTxResponse{Ok: ok}, nil
}

