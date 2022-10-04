package usecase

import (
	"context"

	"github.com/qulaz/gas-price-test/internal/entity"
)

//go:generate go run github.com/golang/mock/mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase

type (
	GasGraph interface {
		Calculate(context.Context) (entity.GasGraphResult, error)
	}

	GasTransactionRepo interface {
		GetGasTransactions(context.Context) ([]*entity.GasTransaction, error)
	}
)
