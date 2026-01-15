package cedra

import (
	"crypto/ed25519"
	"crypto/sha3"

	"github.com/spf13/cast"
)

const (
	CedraCoin                       = "0x1::cedra_coin::CedraCoin"
	transactionPrefix               = "CEDRA::RawTransaction"
	txVariant                 uint8 = 0
	typeTagStruct                   = 7
	transactionPayloadVariant       = 2
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
	Auth                       TxAuthorizer
}

// func (tx *Transaction) SetFeeCoin(coin string) {
// 	tx.FaAddress = coin
// }

func (tx *Transaction) Sign() ([]byte, TxAuthorizer) {
	bcs := NewBCSEncoder()
	bcs.WriteRawBytes(tx.Sender.AccountAddress[:])
	bcs.WriteRawBytes(EncodeUintToBCS(tx.SequenceNumber))
	// Encode Tx payload.
	// Encode entry function.
	bcs.Uleb128(transactionPayloadVariant)
	bcs.WriteRawBytes(tx.Payload.ModuleAddress[:])
	EncodeToBCSString(tx.Payload.ModuleName, bcs)
	EncodeToBCSString(tx.Payload.FunctionName, bcs)
	// txn.Payload.MarshalBCS(ser) // TODO: ???
	// Encode tx payload arguments.
	bcs.Uleb128(0) // ArgTypes
	argsLen := cast.ToUint8(len(tx.Payload.Argumments))
	bcs.Uleb128(argsLen)
	for _, a := range tx.Payload.Argumments {
		bcs.Uleb128(cast.ToUint8(len(a)))
		bcs.WriteRawBytes(a)
	}
	// Encode tx params.
	bcs.WriteRawBytes(EncodeUintToBCS(tx.MaxGasAmount))
	bcs.WriteRawBytes(EncodeUintToBCS(tx.GasUnitPrice))
	bcs.WriteRawBytes(EncodeUintToBCS(tx.ExpirationTimestampSeconds))
	bcs.WriteRawBytes(EncodeUintToBCS(tx.ChainId))

	bcs.Uleb128(typeTagStruct)                  // Encode FaFaAddress
	bcs.WriteRawBytes(tx.FaAddress.Address[:])  // Encode FaFaAddress
	EncodeToBCSString(tx.FaAddress.Module, bcs) // Encode FaFaAddress
	EncodeToBCSString(tx.FaAddress.Name, bcs)   // Encode FaFaAddress
	bcs.Uleb128(0)                              // Encode FaFaAddress

	encodedTx := bcs.GetBytes()
	txPrefix := sha3.Sum256([]byte(transactionPrefix))

	message := []byte{}
	message = append(message, txPrefix[:]...)
	message = append(message, encodedTx...)

	sigBytes := ed25519.Sign(tx.Sender.PrivateKey, message)

	authorizer := TxAuthorizer{
		Variant: 0,
		Auth: Auth{
			PKey:      tx.Sender.PublicKey,
			Signature: sigBytes,
		},
	}

	return encodedTx, authorizer
}

type TxAuthorizer struct {
	Variant uint8
	Auth    Auth
}

type Auth struct {
	PKey      []byte
	Signature []byte
}

type StructTag struct {
	Address [32]byte
	Module  string
	Name    string
}
