package cedra

const (
	txVariant uint8 = 0
)

type CedraAuthenticator struct {
	Variant uint8
	Auth    SenderAuth
}

type SenderAuth struct {
	PKey      []byte
	Signature []byte
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
