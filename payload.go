package cedra

import "github.com/spf13/cast"

const (
	// transactionPayloadVariant is the variant identifier for transaction payloads.
	transactionPayloadVariant = 2
	// txTypedArgsLen is the length of typed arguments (currently always 0).
	txTypedArgsLen = 0
)

// TransactionPayload represents the payload of a Cedra transaction.
// It specifies which module function to call and with what arguments.
type TransactionPayload struct {
	// ModuleAddress is the 32-byte address of the module to call.
	ModuleAddress [32]byte
	// ModuleName is the name of the module containing the function.
	ModuleName string
	// FunctionName is the name of the function to call.
	FunctionName string
	// Arguments is a slice of byte arrays representing the function arguments.
	Arguments [][]byte
}

// ToBCSBytes encodes the transaction payload into Binary Canonical Serialization (BCS) format.
// Returns the serialized byte representation of the payload.
func (p *TransactionPayload) ToBCSBytes() []byte {
	bcs := NewBCSEncoder()
	defer bcs.buf.Reset()
	bcs.EncodeEnum(transactionPayloadVariant)
	bcs.WriteRawBytes(p.ModuleAddress[:])
	bcs.EncodeString(p.ModuleName)
	bcs.EncodeString(p.FunctionName)
	bcs.EncodeEnum(txTypedArgsLen)
	bcs.EncodeEnum(cast.ToUint64(len(p.Arguments)))
	for _, a := range p.Arguments {
		bcs.EncodeEnum(cast.ToUint64(len(a)))
		bcs.WriteRawBytes(a)
	}

	return bcs.GetBytes()
}
