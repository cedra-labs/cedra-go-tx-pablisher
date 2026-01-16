package cedra

import (
	"crypto/ed25519"
	"crypto/sha3"
	"encoding/hex"
	"errors"
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
		return Account{}, errors.New("can't extract account piblick key from account private key")
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

func (a Account) GetAccountAddressString() string {
	return hex.EncodeToString(a.AccountAddress[:])
}

func NewAccountAddress(address string) [32]byte {
	address = strings.TrimPrefix(address, keyPrefix)
	bytes, _ := hex.DecodeString(address)
	var buf [32]byte
	copy((buf)[32-len(bytes):], bytes)

	return buf
}
