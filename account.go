package cedra

import (
	"crypto/ed25519"
	"crypto/sha3"
	"encoding/hex"
	"strings"
)

const (
	privateKeyPrefix = "ed25519-priv-"
	keyPrefix        = "0x"
)

type Account struct {
	AccountAddress [32]byte
	PrivateKey     ed25519.PrivateKey
	PublicKey      ed25519.PublicKey
}

func NewAccount(hexKey string) (Account, error) {
	hexKey = strings.TrimPrefix(hexKey, privateKeyPrefix)
	hexKey = strings.TrimPrefix(hexKey, keyPrefix)

	privBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return Account{}, err
	}

	privateKey := ed25519.NewKeyFromSeed(privBytes)
	publicKey, ok := privateKey.Public().(ed25519.PublicKey)
	if !ok {
		return Account{}, nil // TODO:
	}

	hasher := sha3.New256()
	for _, b := range [][]byte{publicKey, {0}} {
		hasher.Write(b)
	}
	accountAddress := hasher.Sum([]byte{})

	return Account{
		PrivateKey:     privateKey,
		PublicKey:      publicKey,
		AccountAddress: [32]byte(accountAddress),
	}, nil
}

func (a Account) GetAccounAddressString() string {
	return hex.EncodeToString(a.AccountAddress[:])
}
