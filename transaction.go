package cedra

import (
	"crypto/ed25519"
	"crypto/sha3"
	"fmt"

	"github.com/spf13/cast"
)

const (
	CedraCoin         = "0x1::cedra_coin::CedraCoin"
	transactionPrefix = "CEDRA::RawTransaction"
)

type Transaction struct {
	Sender                     Account
	SequenceNumber             uint64
	Payload                    TransactionPayload
	MaxGasAmount               uint64
	GasUnitPrice               uint64
	FaAddress                  string
	ExpirationTimestampSeconds uint64
	ChainId                    uint8
	Auth                       TxAuthorizer
}

func NewRawTransaction(sender Account, payload TransactionPayload) *Transaction {
	sequenceNumber := uint64(1) // TODO get number

	return &Transaction{
		Sender:         sender,
		SequenceNumber: sequenceNumber,
		Payload:        payload,
		FaAddress:      CedraCoin,
	}
}

func (tx *Transaction) SetFeeCoin(coin string) {
	tx.FaAddress = coin
}

func (tx *Transaction) Sign() {
	bcs := NewBCSEncoder()
	bcs.WriteRawBytes(tx.Sender.AccountAddress[:])
	bcs.WriteRawBytes(EncodeUintToBCS(tx.SequenceNumber))
	bcs.WriteRawBytes(tx.Payload.ModuleAddress[:])
	EncodeToBCSString(tx.Payload.ModuleName, bcs)
	EncodeToBCSString(tx.Payload.FunctionName, bcs)
	// txn.Payload.MarshalBCS(ser) // TODO: ???
	argsLen := cast.ToUint8(len(tx.Payload.Argumments))
	bcs.Uleb128(argsLen)
	for _, a := range tx.Payload.Argumments {
		bcs.WriteRawBytes(a)
	}
	bcs.WriteRawBytes(EncodeUintToBCS(tx.MaxGasAmount))
	bcs.WriteRawBytes(EncodeUintToBCS(tx.GasUnitPrice))
	bcs.WriteRawBytes(EncodeUintToBCS(tx.ExpirationTimestampSeconds))
	bcs.WriteRawBytes(EncodeUintToBCS(tx.ChainId))
	EncodeToBCSString(tx.FaAddress, bcs)

	txPrefix := sha3.Sum256([]byte(transactionPrefix))
	encodedTx := bcs.GetBytes()

	message := make([]byte, len(txPrefix)+len(encodedTx))
	message = append(message, txPrefix[:]...)
	message = append(message[len(txPrefix):], encodedTx...)

	sigBytes := ed25519.Sign(tx.Sender.PrivateKey, message)

	fmt.Println(sigBytes)
}

type TxAuthorizer struct{}
