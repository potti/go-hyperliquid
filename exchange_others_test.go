package hyperliquid

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/sonirico/vago/ent"
	"github.com/stretchr/testify/require"
)

func TestWithdrawFromBridgeActionMatchesExchangeSchema(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	nonce := int64(123456789)
	action := buildWithdrawFromBridgeAction(
		198.34,
		"0xABCDEFABCDEFABCDEFABCDEFABCDEFABCDEFABCD",
		nonce,
	)

	payloadTypes := []apitypes.Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "destination", Type: "string"},
		{Name: "amount", Type: "string"},
		{Name: "time", Type: "uint64"},
	}
	signature, err := SignUserSignedAction(
		privateKey,
		action,
		payloadTypes,
		"HyperliquidTransaction:Withdraw",
		true,
	)
	require.NoError(t, err)

	body, err := json.Marshal(map[string]any{
		"action":    action,
		"nonce":     nonce,
		"signature": signature,
	})
	require.NoError(t, err)
	require.Contains(t, string(body), `"type":"withdraw3"`)
	require.Contains(t, string(body), `"signatureChainId":"0x66eee"`)
	require.Contains(t, string(body), `"hyperliquidChain":"Mainnet"`)
	require.Contains(t, string(body), `"amount":"198.340000"`)
	require.NotContains(t, string(body), `"amount":198.34`)
}

func setupExchange(t *testing.T) *Exchange {
	t.Helper()
	_ = loadEnvClean(".env.testnet")
	key := ent.Str("HL_PRIVATE_KEY", "")
	exchange, err := newExchange(key, TestnetAPIURL)
	require.NoError(t, err)
	return exchange
}

func TestPerpDeployHaltTrading(t *testing.T) {
	t.Run("halt trading success response", func(t *testing.T) {
		exchange := setupExchange(t)
		initRecorder(t, false, "PerpDeployHaltTrading_Success")

		res, err := exchange.PerpDeployHaltTrading(context.TODO(), "test:BTC", true)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "ok", res.Status)
	})

	t.Run("halt trading error response", func(t *testing.T) {
		exchange := setupExchange(t)
		initRecorder(t, false, "PerpDeployHaltTrading_Error")

		res, err := exchange.PerpDeployHaltTrading(context.TODO(), "test:BTC", true)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "err", res.Status)
		require.NotEmpty(t, res.Response)
	})
}

func TestPerpDeployRegisterAsset(t *testing.T) {
	oracleUpdater := "0xABCDEF1234567890ABCDEF1234567890ABCDEF12"

	t.Run("with schema", func(t *testing.T) {
		exchange := setupExchange(t)
		initRecorder(t, false, "PerpDeployRegisterAsset_WithSchema")

		maxGas := int(1000000000000)
		res, err := exchange.PerpDeployRegisterAsset(
			context.TODO(),
			"test",
			&maxGas,
			AssetRequest{
				Coin:          "test:TEST0",
				SzDecimals:    2,
				OraclePx:      "10.0",
				MarginTableID: 10,
				OnlyIsolated:  false,
			},
			&PerpDexSchemaInput{
				FullName:        "test dex",
				CollateralToken: 0,
				OracleUpdater:   &oracleUpdater,
			},
		)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "ok", res.Status)
	})

	t.Run("without schema", func(t *testing.T) {
		exchange := setupExchange(t)
		initRecorder(t, false, "PerpDeployRegisterAsset_WithoutSchema")

		res, err := exchange.PerpDeployRegisterAsset(
			context.TODO(),
			"test",
			nil,
			AssetRequest{
				Coin:          "test:TEST0",
				SzDecimals:    2,
				OraclePx:      "10.0",
				MarginTableID: 10,
				OnlyIsolated:  false,
			},
			nil,
		)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "ok", res.Status)
	})

	t.Run("error response", func(t *testing.T) {
		exchange := setupExchange(t)
		initRecorder(t, false, "PerpDeployRegisterAsset_Error")

		res, err := exchange.PerpDeployRegisterAsset(
			context.TODO(),
			"test",
			nil,
			AssetRequest{
				Coin:          "test:TEST0",
				SzDecimals:    2,
				OraclePx:      "10.0",
				MarginTableID: 10,
				OnlyIsolated:  false,
			},
			nil,
		)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "err", res.Status)
		require.NotEmpty(t, res.Response)
	})
}

func TestPerpDeployRegisterAsset2(t *testing.T) {
	oracleUpdater := "0xABCDEF1234567890ABCDEF1234567890ABCDEF12"

	t.Run("with schema strictIsolated", func(t *testing.T) {
		exchange := setupExchange(t)
		initRecorder(t, false, "PerpDeployRegisterAsset2_WithSchema")

		maxGas := int(1000000000000)
		res, err := exchange.PerpDeployRegisterAsset2(
			context.TODO(),
			"test",
			&maxGas,
			AssetRequest2{
				Coin:          "test:TEST0",
				SzDecimals:    2,
				OraclePx:      "10.0",
				MarginTableID: 10,
				MarginMode:    "strictIsolated",
			},
			&PerpDexSchemaInput{
				FullName:        "test dex",
				CollateralToken: 0,
				OracleUpdater:   &oracleUpdater,
			},
		)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "ok", res.Status)
	})

	t.Run("without schema noCross", func(t *testing.T) {
		exchange := setupExchange(t)
		initRecorder(t, false, "PerpDeployRegisterAsset2_WithoutSchema")

		res, err := exchange.PerpDeployRegisterAsset2(
			context.TODO(),
			"test",
			nil,
			AssetRequest2{
				Coin:          "test:TEST0",
				SzDecimals:    2,
				OraclePx:      "10.0",
				MarginTableID: 10,
				MarginMode:    "noCross",
			},
			nil,
		)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "ok", res.Status)
	})
}
