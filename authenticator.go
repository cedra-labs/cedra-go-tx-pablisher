package cedra

import "github.com/spf13/cast"

const (
	// txVariant is the transaction variant identifier for Cedra transactions.
	txVariant uint64 = 0
)

// CedraAuthenticator represents the authentication information for a Cedra transaction.
type CedraAuthenticator struct {
	// Variant specifies the transaction variant type.
	Variant uint64
	// Auth contains the sender authentication data (public key and signature).
	Auth SenderAuth
}

// EncodeBSC encodes the authenticator into Binary Canonical Serialization (BCS) format.
// Returns the serialized byte representation of the authenticator.
func (a CedraAuthenticator) EncodeBSC() []byte {
	bcs := NewBCSEncoder()
	defer bcs.buf.Reset()
	bcs.EncodeEnum(a.Variant)
	bcs.WriteRawBytes(a.Auth.EncodeBSC())

	return bcs.GetBytes()
}

// SenderAuth contains the authentication data for a transaction sender.
type SenderAuth struct {
	// PKey is the public key of the transaction sender.
	PKey []byte
	// Signature is the ED25519 signature of the transaction.
	Signature []byte
}

// EncodeBSC encodes the sender authentication into Binary Canonical Serialization (BCS) format.
// Returns the serialized byte representation of the authentication data.
func (a SenderAuth) EncodeBSC() []byte {
	bcs := NewBCSEncoder()
	defer bcs.buf.Reset()
	bcs.EncodeEnum(cast.ToUint64(len(a.PKey)))
	bcs.WriteRawBytes(a.PKey)
	bcs.EncodeEnum(cast.ToUint64(len(a.Signature)))
	bcs.WriteRawBytes(a.Signature)

	return bcs.GetBytes()
}

// NewSenderAuth creates a new SenderAuth instance with the provided public key and signature.
func NewSenderAuth(pKey []byte, signature []byte) SenderAuth {
	return SenderAuth{
		PKey:      pKey,
		Signature: signature,
	}
}

// NewCedraAuthenticator creates a new CedraAuthenticator with the provided public key and signature.
func NewCedraAuthenticator(pKey []byte, signature []byte) CedraAuthenticator {
	return CedraAuthenticator{
		Variant: txVariant,
		Auth:    NewSenderAuth(pKey, signature),
	}
}
