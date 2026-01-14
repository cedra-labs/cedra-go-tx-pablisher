package main

import (
	"fmt"

	"github.com/cedra-labs/cedra-go-tx-pablisher"
)

const (
	// sender private key.
	privateKey = "85f54f983bd8adcf9aae6729d5075e97ce2d4e6cc4c70eb430c7b80892dd8073"
	// receiver account address.
	receiverAddress        = "3c9124028c90111d7cfd47a28fae30612e397d115c7b78f69713fb729347a77e"
	transferAmount  uint64 = 100000000 // 1 Cedra coin.
)

func main() {
	cedraClient := cedra.NewCedraClient(cedra.TestnetChainID)

	sender, err := cedra.NewAccount(privateKey)
	if err != nil {
		panic(err)
	}

	bcs := cedra.NewBCSEncoder()

	payload := cedra.TransactionPayload{
		ModuleAddress: [32]byte{},
		ModuleName:    "cedra_account",
		FunctionName:  "transfer",
		Argumments: [][]byte{
			cedra.EncodeToBCSString(receiverAddress, bcs),
			cedra.EncodeUintToBCS(transferAmount),
		},
	}

	rawTx, err := cedraClient.NewTransaction(sender, payload)
	encodedTX := rawTx.Sign()

	hash, err := cedraClient.SubmitTransaction(encodedTX)
	if err != nil {
		panic(err)
	}

	fmt.Println(hash)
}
