package main

import (
	"context"
	"fmt"

	"github.com/cedra-labs/cedra-go-tx-pablisher"
)

const (
	// sender private key.
	privateKey = "0x1b542690da83a0c3507e9e8c6ac03be689863d5241483b206a8f5ffd1fefd540"
	// receiver account address.
	receiverAddress        = "c745ffa4f97fa9739fae0cb173996f70bb8e4b0310fa781ccca2f7dc13f7db06"
	transferAmount  uint64 = 100000000 // 1 Cedra coin.
)

func main() {
	ctx := context.TODO()
	cedraClient := cedra.NewCedraClient(cedra.TestnetChainID)

	sender, err := cedra.NewAccount(privateKey)
	if err != nil {
		panic(err)
	}

	moduleAddress, err := cedra.NewAccountAddress(cedra.CedraAddress)
	if err != nil {
		panic(err)
	}
	receiverAddr, err := cedra.NewAccountAddress(receiverAddress)
	if err != nil {
		panic(err)
	}

	payload := cedra.TransactionPayload{
		ModuleAddress: moduleAddress,
		ModuleName:    "cedra_account",
		FunctionName:  "transfer",
		Arguments: [][]byte{
			receiverAddr[:],
			cedra.EncodeUintToBCS(transferAmount),
		},
	}

	rawTx, err := cedraClient.NewTransaction(sender, &payload)
	encodedTx, authenticator := rawTx.Sign()

	hash, err := cedraClient.SubmitTransaction(encodedTx, authenticator)
	if err != nil {
		panic(err)
	}

	isExecuted, err := cedraClient.IsTxExecuted(ctx, hash)
	if err != nil {
		panic(err)
	}

	fmt.Println("executed", isExecuted, hash)
}
