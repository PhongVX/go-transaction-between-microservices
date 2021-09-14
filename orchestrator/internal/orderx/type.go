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
)

