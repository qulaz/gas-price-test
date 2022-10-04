package repository

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/qulaz/gas-price-test/internal/entity"
)

type httpGasTransactionResponse struct {
	Ethereum struct {
		Transactions []*entity.GasTransaction `json:"transactions"`
	} `json:"ethereum"`
}

type HttpGasTransactionRepo struct {
	client http.Client
	url    string
}

func NewHttpGasTransactionRepo(url string) *HttpGasTransactionRepo {
	return &HttpGasTransactionRepo{
		client: http.Client{ //nolint: exhaustruct
			Transport: http.DefaultTransport,
			Timeout:   time.Second * 5,
		},
		url: url,
	}
}

func (h *HttpGasTransactionRepo) GetGasTransactions(ctx context.Context) ([]*entity.GasTransaction, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	rawResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var httpResp httpGasTransactionResponse
	if err := json.Unmarshal(rawResp, &httpResp); err != nil {
		return nil, err
	}

	return httpResp.Ethereum.Transactions, nil
}
