package swap_test

import (
	"context"
	"fmt"
	"testing"
	"encoding/json"

	"github.com/test-go/testify/require"
	"github.com/zkweb3/jupiter-go/jupiter"
	"github.com/zkweb3/jupiter-go/solana"
)

const (
	testPrivateKey string = "5473ZnvEhn35BdcCcPLKnzsyP6TsgqQrNFpn4i2gFegFiiJLyWginpa9GoFn2cy6Aq2EAuxLt2u2bjFDBPvNY6nw"
	BOME_tokenAddress string = "ukHH6c7mMyiWCf1b9pnWe25TSpkDDt3H5pQZgZ74J82"
)

func TestJupQuoteApi(t *testing.T) {
	wallet, err := solana.NewWalletFromPrivateKeyBase58(testPrivateKey)
	require.NoError(t, err)

	jupClient, err := jupiter.NewClientWithResponses(jupiter.DefaultAPIURL)
	require.Nil(t, err)
	ctx := context.TODO()

	slippageBps := 50	// 0.5% slippage

	// More info: https://station.jup.ag/docs/apis/swap-api
	// Swapping SOL to BOME with input 1 SOL and 0.5% slippage
	quoteResponse, err := jupClient.GetQuoteWithResponse(ctx, &jupiter.GetQuoteParams{
		InputMint:   "So11111111111111111111111111111111111111112",
		OutputMint:  BOME_tokenAddress,
		Amount:      1000000000,	// 1 SOL
		SlippageBps: &slippageBps,	// 0.5% slippage
	})
	require.Nil(t, err)
	require.NotNil(t, quoteResponse.JSON200)
	quote := quoteResponse.JSON200
	fmt.Println("quote.outAmount", quote.OutAmount)

	// More info: https://station.jup.ag/docs/apis/troubleshooting
	prioritizationFeeLamports := jupiter.SwapRequest_PrioritizationFeeLamports{}
	dynamicComputeUnitLimit := true		// allow dynamic compute limit instead of max 1,400,000

	err = json.Unmarshal([]byte("{\"autoMultiplier\":2}"), &prioritizationFeeLamports)	// which will 2x of the auto fees
	require.Nil(t, err)

	// Get instructions for a swap
	swapResponse, err := jupClient.PostSwapWithResponse(ctx, jupiter.PostSwapJSONRequestBody{
		PrioritizationFeeLamports: &prioritizationFeeLamports,
		QuoteResponse:             *quote,
		UserPublicKey:             wallet.PublicKey().String(),
		DynamicComputeUnitLimit:   &dynamicComputeUnitLimit,
	})
	require.Nil(t, err)
	fmt.Println("swapResponse.status", swapResponse.Status())
	require.NotNil(t, swapResponse.JSON200)
	swap := swapResponse.JSON200
	fmt.Println("swapTransaction", swap.SwapTransaction)
}
