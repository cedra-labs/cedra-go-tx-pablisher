package cedra

import (
	"crypto/ed25519"
	"crypto/sha3"

	"github.com/pkg/errors"
)

const (
	// transactionPrefix is the prefix used when signing transactions.
	transactionPrefix = "CEDRA::RawTransaction"
)

type SequenceNumber uint64

func (s SequenceNumber) ToUint64() uint64 {
	return uint64(s)
}

func (s SequenceNumber) ToBCSBytes() []byte {
	return EncodeUintToBCS(s.ToUint64())
}

type MaxGasAmount uint64

func (s MaxGasAmount) ToUint64() uint64 {
	return uint64(s)
}

func (s MaxGasAmount) ToBCSBytes() []byte {
	return EncodeUintToBCS(s.ToUint64())
}

type GasUnitPrice uint64

func (s GasUnitPrice) ToUint64() uint64 {
	return uint64(s)
}

func (s GasUnitPrice) ToBCSBytes() []byte {
	return EncodeUintToBCS(s.ToUint64())
}

// Transaction represents a complete Cedra blockchain transaction.
// Fields are grouped by size for optimal memory alignment.
type Transaction struct {
	// Sender is the account that will sign and submit the transaction.
	Sender Account
	// Payload contains the transaction payload (module, function, arguments).
	Payload TransactionPayload
	// FaAddress is the struct tag for the fee asset (coin type).
	FaAddress StructTag
	// SequenceNumber is the sequence number for the sender account.
	SequenceNumber SequenceNumber
	// MaxGasAmount is the maximum amount of gas units the transaction can consume.
	MaxGasAmount MaxGasAmount
	// GasUnitPrice is the price per gas unit.
	GasUnitPrice GasUnitPrice
	// ExpirationTimestampSeconds is the Unix timestamp when the transaction expires.
	ExpirationTimestampSeconds uint64
	// ChainId identifies the blockchain network.
	ChainId uint8
}

// SetFeeCoin sets the fee coin type for the transaction.
// Currently not implemented and returns an error.
func (tx *Transaction) SetFeeCoin(coin string) error {
	// TODO: Implement
	return errors.New("unimplemeted")
}

// ToBCSBytes encodes the transaction into Binary Canonical Serialization (BCS) format.
// Returns the serialized byte representation of the transaction.
func (tx *Transaction) ToBCSBytes() []byte {
	bcs := NewBCSEncoder()
	defer bcs.buf.Reset()
	bcs.WriteRawBytes(tx.Sender.AccountAddress[:])
	bcs.WriteRawBytes(tx.SequenceNumber.ToBCSBytes())
	bcs.WriteRawBytes(tx.Payload.ToBCSBytes())
	bcs.WriteRawBytes(tx.MaxGasAmount.ToBCSBytes())
	bcs.WriteRawBytes(tx.GasUnitPrice.ToBCSBytes())
	bcs.WriteRawBytes(EncodeUintToBCS(tx.ExpirationTimestampSeconds))
	bcs.WriteRawBytes(EncodeUintToBCS(tx.ChainId))
	bcs.WriteRawBytes(tx.FaAddress.ToBCSBytes())

	return bcs.GetBytes()
}

// Sign signs the transaction using the sender's private key and creates an authenticator.
// Returns the encoded transaction bytes and the authenticator for submission.
// The transaction is signed with the ED25519 private key after hashing with the transaction prefix.
func (tx *Transaction) Sign() ([]byte, CedraAuthenticator) {
	encodedTx := tx.ToBCSBytes()
	txPrefix := sha3.Sum256([]byte(transactionPrefix))

	message := make([]byte, 0, len(txPrefix)+len(encodedTx))
	message = append(message, txPrefix[:]...)
	message = append(message, encodedTx...)

	signature := ed25519.Sign(tx.Sender.PrivateKey, message)
	authenticator := NewCedraAuthenticator(tx.Sender.PublicKey, signature)

	return encodedTx, authenticator
}
