package cedra

import "github.com/spf13/cast"

const (
	txVariant uint64 = 0
)

type CedraAuthenticator struct {
	Variant uint64
	Auth    SenderAuth
}

func (a CedraAuthenticator) EncodeBSC() []byte {
	bcs := NewBCSEncoder()
	defer bcs.buf.Reset()
	bcs.EncodeEnum(a.Variant)
	bcs.WriteRawBytes(a.Auth.EncodeBSC())

	return bcs.GetBytes()
}

type SenderAuth struct {
	PKey      []byte
	Signature []byte
}

func (a SenderAuth) EncodeBSC() []byte {
	bcs := NewBCSEncoder()
	defer bcs.buf.Reset()
	length := cast.ToUint64(len(a.PKey))
	bcs.EncodeEnum(length)
	bcs.WriteRawBytes(a.PKey)
	length = cast.ToUint64(len(a.Signature))
	bcs.EncodeEnum(length)
	bcs.WriteRawBytes(a.Signature)

	return bcs.GetBytes()
}

func NewSenderAuth(pKey []byte, signature []byte) SenderAuth {
	return SenderAuth{
		PKey:      pKey,
		Signature: signature,
	}
}

func NewCedraAuthenticator(pKey []byte, signature []byte) CedraAuthenticator {
	return CedraAuthenticator{
		Variant: txVariant,
		Auth:    NewSenderAuth(pKey, signature),
	}
}
