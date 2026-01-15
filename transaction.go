package cedra

import (
	"crypto/ed25519"
	"crypto/sha3"
	"errors"
)

const (
	transactionPrefix = "CEDRA::RawTransaction"
)

type Transaction struct {
	Sender                     Account
	SequenceNumber             uint64
	Payload                    TransactionPayload
	MaxGasAmount               uint64
	GasUnitPrice               uint64
	FaAddress                  StructTag
	ExpirationTimestampSeconds uint64
	ChainId                    uint8
}

func (tx *Transaction) SetFeeCoin(coin string) error {
	// TODO: Implement
	return errors.New("unimplemeted")
}

func (tx Transaction) ToBCSBytes() []byte {
	bcs := NewBCSEncoder()
	defer bcs.buf.Reset()

	bcs.WriteRawBytes(tx.Sender.AccountAddress[:])
	bcs.WriteRawBytes(EncodeUintToBCS(tx.SequenceNumber))
	bcs.WriteRawBytes(tx.Payload.ToBCSBytes())
	bcs.WriteRawBytes(EncodeUintToBCS(tx.MaxGasAmount))
	bcs.WriteRawBytes(EncodeUintToBCS(tx.GasUnitPrice))
	bcs.WriteRawBytes(EncodeUintToBCS(tx.ExpirationTimestampSeconds))
	bcs.WriteRawBytes(EncodeUintToBCS(tx.ChainId))

	bcs.WriteRawBytes(tx.FaAddress.ToBCSBytes())
	encodedTx := bcs.GetBytes()

	return encodedTx
}

func (tx *Transaction) Sign() ([]byte, CedraAuthenticator) {
	encodedTx := tx.ToBCSBytes()
	txPrefix := sha3.Sum256([]byte(transactionPrefix))

	message := []byte{}
	message = append(message, txPrefix[:]...)
	message = append(message, encodedTx...)

	signature := ed25519.Sign(tx.Sender.PrivateKey, message)
	authenticator := NewCedraAuthenticator(tx.Sender.PublicKey, signature)

	return encodedTx, authenticator
}
