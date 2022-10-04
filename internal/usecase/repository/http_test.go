package repository

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/qulaz/gas-price-test/internal/entity"
)

func TestHttpGasTransactionRepo_GetGasTransactions(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_, repo, expectedTransactions, teardown := setupTest(t)
		defer teardown()

		transactions, err := repo.GetGasTransactions(context.Background())
		require.NoError(t, err)
		assert.Equal(t, expectedTransactions, transactions)
	})
	t.Run("url unreached", func(t *testing.T) {
		server, repo, _, teardown := setupTest(t)
		defer teardown()

		server.Close()

		transactions, err := repo.GetGasTransactions(context.Background())
		require.Error(t, err)
		assert.Nil(t, transactions)
	})
}

func setupTest(t *testing.T) (
	*httptest.Server,
	*HttpGasTransactionRepo,
	[]*entity.GasTransaction,
	func(),
) {
	t.Helper()

	f, err := os.Open("testdata/gas_price.json")
	require.NoError(t, err)

	rawGasPrice, err := io.ReadAll(f)
	require.NoError(t, err)

	mux := http.NewServeMux()
	mux.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(rawGasPrice)
		require.NoError(t, err)
	})

	server := httptest.NewServer(mux)

	repo := NewHttpGasTransactionRepo(server.URL + "/transactions")

	var transactions httpGasTransactionResponse
	err = json.Unmarshal(rawGasPrice, &transactions)
	require.NoError(t, err)

	return server, repo, transactions.Ethereum.Transactions, func() {
		server.Close()
	}
}
