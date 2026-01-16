package cedra

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

const (
	defaultMaxGasAmount = uint64(100_000)
	CedraAddress        = "0x0000000000000000000000000000000000000000000000000000000000000001"
	CedraCoin           = "0x0000000000000000000000000000000000000000000000000000000000000001::cedra_coin::CedraCoin"
)

type CedraClient struct {
	node    CedraNode
	chainID ChainID
}

func NewCedraClient(chainID ChainID) CedraClient {
	return CedraClient{
		node:    NewCedraNode(chainID),
		chainID: chainID,
	}
}

func (c CedraClient) NewTransaction(sender Account, payload TransactionPayload) (*Transaction, error) {
	wg := &sync.WaitGroup{}
	var err error
	var sequenceNumber uint64
	var gasUnitPrice uint64

	expirationSeconds := cast.ToUint64(time.Now().Unix() + 300)
	wg.Add(1)
	go func() {
		defer wg.Done()
		sequenceNumber, err = c.node.GetSequenceNumber(sender.GetAccountAddressString())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		estimate, err := c.node.GetEstimateGasPrice()
		if err == nil {
			gasUnitPrice = estimate.GasEstimate
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

	}()

	wg.Wait()

	if err != nil {
		return nil, errors.Wrap(err, "can't create new transaction")
	}
	structTag, err := NewStringStructTag(CedraCoin)
	if err != nil {
		return nil, errors.Wrap(err, "can't create new transaction: invalid struct tag")
	}

	return &Transaction{
		Sender:                     sender,
		SequenceNumber:             sequenceNumber,
		Payload:                    payload,
		FaAddress:                  structTag,
		GasUnitPrice:               gasUnitPrice,
		MaxGasAmount:               defaultMaxGasAmount,
		ExpirationTimestampSeconds: expirationSeconds,
		ChainId:                    uint8(c.chainID),
	}, nil
}

func (c CedraClient) SubmitTransaction(tx []byte, auth CedraAuthenticator) (string, error) {
	tx = append(tx, auth.EncodeBSC()...)

	hash, err := c.node.SubmitTransaction(tx)
	if err != nil {
		return "", err
	}
	return hash, nil
}
