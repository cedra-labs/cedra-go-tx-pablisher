package cedra

import (
	"encoding/hex"
	"sync"
	"time"

	"github.com/spf13/cast"
)

const (
	defaultMaxGasAmount = uint64(100_000)
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
		sequenceNumber, err = c.node.GetSequenceNumber(sender.GetAccounAddressString())
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
		return nil, err // TODO:
	}
	bytes, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		panic(err)
	}

	var mAdd [32]byte
	copy((mAdd)[32-len(bytes):], bytes)

	return &Transaction{
		Sender:         sender,
		SequenceNumber: sequenceNumber,
		Payload:        payload,
		FaAddress: StructTag{
			Address: mAdd,
			Module:  "cedra_coin",
			Name:    "CedraCoin",
		},
		GasUnitPrice:               gasUnitPrice,
		MaxGasAmount:               defaultMaxGasAmount,
		ExpirationTimestampSeconds: expirationSeconds,
		ChainId:                    uint8(c.chainID),
	}, nil
}

func (c CedraClient) SubmitTransaction(tx []byte, auth TxAuthorizer) (string, error) {
	bcs := NewBCSEncoder()
	bcs.Uleb128(txVariant)
	//pub
	length := cast.ToUint8(len(auth.Auth.PKey))
	bcs.Uleb128(length)
	bcs.WriteRawBytes(auth.Auth.PKey)
	//sig
	length = cast.ToUint8(len(auth.Auth.Signature))
	bcs.Uleb128(length)
	bcs.WriteRawBytes(auth.Auth.Signature)

	tx = append(tx, bcs.GetBytes()...)

	hash, err := c.node.SubmitTransaction(tx)
	if err != nil {
		return "", err
	}
	return hash, nil
}
