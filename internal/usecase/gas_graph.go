package usecase

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"github.com/qulaz/gas-price-test/internal/entity"
	"github.com/qulaz/gas-price-test/pkg/cache"
	"github.com/qulaz/gas-price-test/pkg/logging"
)

const cacheKey = "gasGraph"

type gasCountSum struct {
	count int
	sum   decimal.Decimal
}

type GasGraphUseCase struct {
	repo     GasTransactionRepo
	cache    cache.ExpiredCache[string, entity.GasGraphResult]
	cacheTtl time.Duration
	logger   logging.ContextLogger
}

func NewGasGraphUseCase(
	repo GasTransactionRepo,
	cache cache.ExpiredCache[string, entity.GasGraphResult],
	cacheTtl time.Duration,
	logger logging.ContextLogger,
) *GasGraphUseCase {
	return &GasGraphUseCase{
		repo:     repo,
		cache:    cache,
		cacheTtl: cacheTtl,
		logger:   logger,
	}
}

func (g *GasGraphUseCase) Calculate(ctx context.Context) (entity.GasGraphResult, error) {
	ctx, logger := g.logger.FromContext(ctx, "usecase", "gasGraph", "method", "Calculate")

	res, err := g.cache.Get(cacheKey)
	if err == nil {
		return res, nil
	}

	res, err = g.calculate(ctx)
	if err != nil {
		return entity.GasGraphResult{}, err
	}

	err = g.cache.Upsert(cacheKey, res, g.cacheTtl)
	if err != nil {
		logger.Warnw("error while set transactions in cache", "err", err)
	}

	return res, nil
}

func (g *GasGraphUseCase) calculate(ctx context.Context) (entity.GasGraphResult, error) {
	transactions, err := g.repo.GetGasTransactions(ctx)
	if err != nil {
		return entity.GasGraphResult{}, err
	}

	monthAmountResChan := make(chan []entity.GasMonthAmountResult)
	dayAverageResChan := make(chan []entity.GasDayAverageResult)
	hourFreqResChan := make(chan []entity.GasHourFreqResult)
	totalPaidResChan := make(chan decimal.Decimal)

	go calculateMonthGasAmount(transactions, monthAmountResChan)
	go calculateDayAverageGasPrice(transactions, dayAverageResChan)
	go calculateHourFreqGasPrice(transactions, hourFreqResChan)
	go calculateTotalGasPaid(transactions, totalPaidResChan)

	return entity.GasGraphResult{
		GasPerMonthAmount: <-monthAmountResChan,
		GasPerDayAverage:  <-dayAverageResChan,
		GasHourFreq:       <-hourFreqResChan,
		GasSpentTotal:     <-totalPaidResChan,
	}, nil
}

func calculateMonthGasAmount(transactions []*entity.GasTransaction, resChan chan<- []entity.GasMonthAmountResult) {
	monthAmount := make(map[entity.TransactionTime]decimal.Decimal)

	for i := range transactions {
		month := transactions[i].Time.TruncToMonth()
		monthAmount[month] = monthAmount[month].Add(transactions[i].GasValue)
	}

	res := make([]entity.GasMonthAmountResult, 0, len(monthAmount))

	for month, gasAmount := range monthAmount {
		res = append(res, entity.GasMonthAmountResult{
			Date:   month,
			Amount: gasAmount,
		})
	}

	resChan <- res
}

func calculateDayAverageGasPrice(transactions []*entity.GasTransaction, resChan chan<- []entity.GasDayAverageResult) {
	dayGasInfo := make(map[entity.TransactionTime]gasCountSum)

	for i := range transactions {
		day := transactions[i].Time.TruncToDay()

		gasInfo, ok := dayGasInfo[day]
		if ok {
			gasInfo.count++
			gasInfo.sum = dayGasInfo[day].sum.Add(transactions[i].GasPrice)
		} else {
			dayGasInfo[day] = gasCountSum{count: 1, sum: transactions[i].GasPrice}
		}
	}

	res := make([]entity.GasDayAverageResult, 0, len(dayGasInfo))

	for day, gasInfo := range dayGasInfo {
		res = append(res, entity.GasDayAverageResult{
			Date:    day,
			Average: gasInfo.sum.Div(decimal.NewFromInt(int64(gasInfo.count))),
		})
	}

	resChan <- res
}

func calculateHourFreqGasPrice(transactions []*entity.GasTransaction, resChan chan<- []entity.GasHourFreqResult) {
	hourGasInfo := make(map[int]gasCountSum)

	for i := range transactions {
		hour := transactions[i].Time.Hour()

		gasInfo, ok := hourGasInfo[hour]
		if ok {
			gasInfo.count++
			gasInfo.sum = hourGasInfo[hour].sum.Add(transactions[i].GasPrice)
		} else {
			hourGasInfo[hour] = gasCountSum{count: 1, sum: transactions[i].GasPrice}
		}
	}

	res := make([]entity.GasHourFreqResult, 0, len(hourGasInfo))

	for hour, gasInfo := range hourGasInfo {
		res = append(res, entity.GasHourFreqResult{
			Hour:  hour,
			Value: gasInfo.sum.Div(decimal.NewFromInt(int64(gasInfo.count))),
		})
	}

	resChan <- res
}

func calculateTotalGasPaid(transactions []*entity.GasTransaction, resChan chan<- decimal.Decimal) {
	var res decimal.Decimal

	for i := range transactions {
		paid := transactions[i].GasValue.Mul(transactions[i].GasPrice)
		res = res.Add(paid)
	}

	resChan <- res
}
