package main

import (
	"fmt"

	"github.com/cedra-labs/cedra-go-tx-pablisher"
)

const (
	// sender private key.
	privateKey = "sedner_private_key"
	// receiver account address.
	receiverAddress        = "reciver_address"
	transferAmount  uint64 = 100000000 // 1 Cedra coin.
)

func main() {
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

	fmt.Println(hash)
}
