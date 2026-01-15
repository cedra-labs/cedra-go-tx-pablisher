package cedra

import "github.com/spf13/cast"

const (
	transactionPayloadVariant = 2
	txTypedArgsLen            = 0
)

type TransactionPayload struct {
	ModuleAddress [32]byte
	ModuleName    string
	FunctionName  string
	Argumments    [][]byte
}

func (p TransactionPayload) ToBCSBytes() []byte {
	bcs := NewBCSEncoder()
	defer bcs.buf.Reset()

	bcs.EncodeEnum(transactionPayloadVariant)
	bcs.WriteRawBytes(p.ModuleAddress[:])
	bcs.EncodeString(p.ModuleName)
	bcs.EncodeString(p.FunctionName)
	bcs.EncodeEnum(txTypedArgsLen)
	argsLen := cast.ToUint64(len(p.Argumments))
	bcs.EncodeEnum(argsLen)
	for _, a := range p.Argumments {
		bcs.EncodeEnum(cast.ToUint64(len(a)))
		bcs.WriteRawBytes(a)
	}

	encodedPayload := bcs.GetBytes()

	return encodedPayload
}
