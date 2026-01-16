// Package cedra provides a Go library for creating and submitting transactions to the Cedra blockchain.
// It supports account management, transaction creation, signing, and submission to Cedra network nodes.
package cedra

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

const (
	// defaultMaxGasAmount is the default maximum gas amount for transactions.
	defaultMaxGasAmount = uint64(100_000)
	// CedraAddress is the canonical address of the Cedra system.
	CedraAddress = "0x0000000000000000000000000000000000000000000000000000000000000001"
	// CedraCoin is the full struct tag for the Cedra coin type.
	CedraCoin = "0x0000000000000000000000000000000000000000000000000000000000000001::cedra_coin::CedraCoin"
)

// CedraClient is the main client for interacting with the Cedra blockchain.
// It provides methods for creating and submitting transactions.
type CedraClient struct {
	// node is the Cedra node client used for network communication.
	node CedraNode
	// chainID identifies the blockchain network (devnet, testnet, mainnet).
	chainID ChainID
}

// NewCedraClient creates a new CedraClient instance for the specified chain.
func NewCedraClient(chainID ChainID) CedraClient {
	return CedraClient{
		node:    NewCedraNode(chainID),
		chainID: chainID,
	}
}

// NewTransaction creates a new transaction with the provided sender and payload.
// It concurrently fetches the sequence number and gas price estimate from the network.
// The transaction expiration is set to 5 minutes from creation time.
// Returns an error if the sequence number cannot be fetched or if the struct tag is invalid.
func (c CedraClient) NewTransaction(sender Account, payload *TransactionPayload) (*Transaction, error) {
	type seqResult struct {
		value uint64
		err   error
	}
	type gasResult struct {
		value uint64
		err   error
	}

	expirationSeconds := cast.ToUint64(time.Now().Unix() + 300)
	seqChan := make(chan seqResult, 1)
	gasChan := make(chan gasResult, 1)

	// Fetch sequence number
	go func() {
		seqNum, err := c.node.GetSequenceNumber(sender.GetAccountAddressString())
		seqChan <- seqResult{value: seqNum, err: err}
	}()

	// Fetch gas price estimate (non-critical, can fail)
	go func() {
		estimate, err := c.node.GetEstimateGasPrice()
		gasPrice := uint64(0)
		if err == nil {
			gasPrice = estimate.GasEstimate
		}
		gasChan <- gasResult{value: gasPrice, err: err}
	}()

	seqRes := <-seqChan
	if seqRes.err != nil {
		return nil, errors.Wrap(seqRes.err, "can't create new transaction: failed to get sequence number")
	}

	gasRes := <-gasChan
	// Gas price estimation failure is non-critical, we can proceed with 0

	structTag, err := NewStringStructTag(CedraCoin)
	if err != nil {
		return nil, errors.Wrap(err, "can't create new transaction: invalid struct tag")
	}

	return &Transaction{
		Sender:                     sender,
		SequenceNumber:             seqRes.value,
		Payload:                    *payload,
		FaAddress:                  structTag,
		GasUnitPrice:               gasRes.value,
		MaxGasAmount:               defaultMaxGasAmount,
		ExpirationTimestampSeconds: expirationSeconds,
		ChainId:                    uint8(c.chainID),
	}, nil
}

// SubmitTransaction submits a signed transaction to the Cedra network.
// The transaction bytes and authenticator are combined and sent to the node.
// Returns the transaction hash if successful, or an error if submission fails.
func (c CedraClient) SubmitTransaction(tx []byte, auth CedraAuthenticator) (string, error) {
	authBytes := auth.EncodeBSC()
	signedTx := make([]byte, 0, len(tx)+len(authBytes))
	signedTx = append(signedTx, tx...)
	signedTx = append(signedTx, authBytes...)

	hash, err := c.node.SubmitTransaction(signedTx)
	if err != nil {
		return "", errors.Wrap(err, "can't submit transaction")
	}

	return hash, nil
}
