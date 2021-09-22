package transactionx

import (
	"database/sql"
	"orchestrator/internal/redisx"
	"orchestrator/internal/transactioncache"
)

type (
	Service struct {
		MapTx        map[string]*sql.Tx
		DB           *sql.DB
		redisService redisx.ServiceI
		txCacheSrv   transactioncache.ServiceI
	}

	GService struct {
		txSrv ServiceI
	}

)
