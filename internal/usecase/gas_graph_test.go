package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/qulaz/gas-price-test/internal/entity"
	"github.com/qulaz/gas-price-test/pkg/cache"
	"github.com/qulaz/gas-price-test/pkg/cache/cachemock"
	"github.com/qulaz/gas-price-test/pkg/logging"
)

func setupTest(t *testing.T) (
	*GasGraphUseCase,
	*MockGasTransactionRepo,
	*cachemock.MockExpiredCache[string, entity.GasGraphResult],
	[]*entity.GasTransaction,
	entity.GasGraphResult,
	func(),
) {
	t.Helper()

	ctrl := gomock.NewController(t)
	repo := NewMockGasTransactionRepo(ctrl)
	mockCache := cachemock.NewMockExpiredCache[string, entity.GasGraphResult](ctrl)

	transactions := getTransactions(t)
	rawGasGraph := loadExpectedGasGraph(t)

	var gasGraph entity.GasGraphResult

	err := json.Unmarshal(rawGasGraph, &gasGraph)
	require.NoError(t, err)

	usecase := NewGasGraphUseCase(
		repo,
		mockCache,
		time.Second,
		logging.NewDummyLogger(),
	)

	return usecase, repo, mockCache, transactions, gasGraph, func() {
		ctrl.Finish()
	}
}

func TestGasGraphUseCase_Calculate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc, _, mockCache, _, expectedGasGraph, teardown := setupTest(t)
		defer teardown()

		mockCache.EXPECT().Get(cacheKey).Return(expectedGasGraph, nil)

		gasGraph, err := uc.Calculate(context.Background())
		require.NoError(t, err)
		requireExpectedGasGraph(t, gasGraph)
	})
	t.Run("cache miss", func(t *testing.T) {
		uc, repo, mockCache, transactions, _, teardown := setupTest(t)
		defer teardown()

		mockCache.EXPECT().Get(cacheKey).Return(entity.GasGraphResult{}, cache.ErrKeyNotFound) //nolint: exhaustruct
		repo.EXPECT().GetGasTransactions(gomock.Any()).Return(transactions, nil)
		mockCache.EXPECT().Upsert(cacheKey, gomock.Any(), uc.cacheTtl).Return(nil)

		gasGraph, err := uc.Calculate(context.Background())
		require.NoError(t, err)
		requireExpectedGasGraph(t, gasGraph)
	})
	t.Run("cache upsert error", func(t *testing.T) {
		uc, repo, mockCache, transactions, _, teardown := setupTest(t)
		defer teardown()

		mockCache.EXPECT().Get(cacheKey).Return(entity.GasGraphResult{}, cache.ErrKeyNotFound) //nolint: exhaustruct
		repo.EXPECT().GetGasTransactions(gomock.Any()).Return(transactions, nil)
		mockCache.EXPECT().Upsert(cacheKey, gomock.Any(), uc.cacheTtl).Return(errors.New(""))

		gasGraph, err := uc.Calculate(context.Background())
		require.NoError(t, err)
		requireExpectedGasGraph(t, gasGraph)
	})
	t.Run("repo error", func(t *testing.T) {
		uc, repo, mockCache, _, _, teardown := setupTest(t)
		defer teardown()

		mockCache.EXPECT().Get(cacheKey).Return(entity.GasGraphResult{}, cache.ErrKeyNotFound) //nolint: exhaustruct
		repo.EXPECT().GetGasTransactions(gomock.Any()).Return(nil, errors.New(""))

		gasGraph, err := uc.Calculate(context.Background())
		require.Error(t, err)
		assert.Equal(t, entity.GasGraphResult{}, gasGraph) //nolint: exhaustruct
	})
}

func requireExpectedGasGraph(t *testing.T, res entity.GasGraphResult) {
	t.Helper()

	expectedRes := loadExpectedGasGraph(t)

	sort.Slice(res.GasPerMonthAmount, func(i, j int) bool {
		return res.GasPerMonthAmount[i].Date.Unix() < res.GasPerMonthAmount[j].Date.Unix()
	})
	sort.Slice(res.GasPerDayAverage, func(i, j int) bool {
		return res.GasPerDayAverage[i].Date.Unix() < res.GasPerDayAverage[j].Date.Unix()
	})
	sort.Slice(res.GasHourFreq, func(i, j int) bool {
		return res.GasHourFreq[i].Hour < res.GasHourFreq[j].Hour
	})

	rawRes, err := json.MarshalIndent(res, "", "  ")
	require.NoError(t, err)
	require.Equal(t, expectedRes, rawRes)
}

func loadExpectedGasGraph(t *testing.T) []byte {
	t.Helper()

	f, err := os.Open("testdata/expected_result.json")
	require.NoError(t, err)

	rawTransactions, err := io.ReadAll(f)
	require.NoError(t, err)

	return rawTransactions
}

func getTransactions(t *testing.T) []*entity.GasTransaction {
	t.Helper()

	f, err := os.Open("testdata/transactions.json")
	require.NoError(t, err)

	defer f.Close()

	rawTransactions, err := io.ReadAll(f)
	require.NoError(t, err)

	var transactions []*entity.GasTransaction

	err = json.Unmarshal(rawTransactions, &transactions)
	require.NoError(t, err)

	return transactions
}
