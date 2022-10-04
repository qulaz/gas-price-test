package config

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T, env map[string]string) func() {
	t.Helper()

	for key, value := range env {
		t.Setenv(key, value)
	}

	// teardown test
	return func() {
		once = sync.Once{}
		config = nil
	}
}

func TestGetConfig_Defaults(t *testing.T) {
	defer setupTest(t, map[string]string{})()

	c, err := GetConfig()
	require.NoError(t, err)

	assert.Equal(t, true, c.API.Debug)
	assert.Equal(t, "0.0.0.0", c.API.Host)
	assert.Equal(t, "8000", c.API.Port)
	assert.Equal(t, "", c.API.Domain)

	assert.Equal(
		t,
		"https://github.com/CryptoRStar/GasPriceTestTask/raw/main/gas_price.json",
		c.API.TransactionsUrl,
	)
	assert.Equal(t, time.Second*30, c.API.GasGraphTtl)
}

func TestGetConfig_RewriteDefaults(t *testing.T) {
	env := map[string]string{
		"DEBUG":               "true",
		"HOST":                "localhost",
		"PORT":                "5555",
		"DOMAIN":              "google.com",
		"TRANSACTIONS_URL":    "https://etherium.io/transactions.json",
		"GAS_GRAPH_CACHE_TTL": "15m",
	}

	defer setupTest(t, env)()

	c, err := GetConfig()
	require.NoError(t, err)

	assert.Equal(t, true, c.API.Debug)
	assert.Equal(t, env["HOST"], c.API.Host)
	assert.Equal(t, env["PORT"], c.API.Port)
	assert.Equal(t, env["DOMAIN"], c.API.Domain)

	assert.Equal(t, env["TRANSACTIONS_URL"], c.API.TransactionsUrl)
	assert.Equal(t, time.Minute*15, c.API.GasGraphTtl)
}
