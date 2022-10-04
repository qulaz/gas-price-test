package entity

import (
	"github.com/shopspring/decimal"
)

type GasTransaction struct {
	Time           TransactionTime `json:"time"`
	GasPrice       decimal.Decimal `json:"gasPrice"`
	GasValue       decimal.Decimal `json:"gasValue"`
	Average        decimal.Decimal `json:"average"`
	MaxGasPrice    decimal.Decimal `json:"maxGasPrice"`
	MedianGasPrice decimal.Decimal `json:"medianGasPrice"`
}

type (
	GasMonthAmountResult struct {
		Date   TransactionTime `json:"date" swaggertype:"primitive,string" example:"22-01-01 00:00"`
		Amount decimal.Decimal `json:"amount"`
	}

	GasDayAverageResult struct {
		Date    TransactionTime `json:"date" swaggertype:"primitive,string" example:"22-01-01 00:00"`
		Average decimal.Decimal `json:"average"`
	}

	GasHourFreqResult struct {
		Hour  int             `json:"hour" minimum:"0" maximum:"24"`
		Value decimal.Decimal `json:"value"`
	}

	GasGraphResult struct {
		GasPerMonthAmount []GasMonthAmountResult `json:"gasPerMonthAmount"`
		GasPerDayAverage  []GasDayAverageResult  `json:"gasPerDayAverage"`
		GasHourFreq       []GasHourFreqResult    `json:"gasHourFreq"`
		GasSpentTotal     decimal.Decimal        `json:"gasSpentTotal"`
	}
)
