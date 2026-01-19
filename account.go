package cedra

import (
	"crypto/ed25519"
	"crypto/sha3"
	"encoding/hex"
	"strings"

	"github.com/pkg/errors"
)

const (
	// privateKeyPrefix is the prefix used for ED25519 private keys.
	privateKeyPrefix = "ed25519-priv-"
	// keyPrefix is the hexadecimal prefix used for addresses and keys.
	keyPrefix = "0x"

	deriveResourceAccountSchema = 0xFF
)

// Account represents a Cedra blockchain account with its cryptographic keys and address.
type Account struct {
	// AccountAddress is the 32-byte account address derived from the public key.
	AccountAddress [32]byte
	// PrivateKey is the ED25519 private key used for signing transactions.
	PrivateKey ed25519.PrivateKey
	// PublicKey is the ED25519 public key associated with the account.
	PublicKey ed25519.PublicKey
}

// NewAccount creates a new Account from a hexadecimal private key string.
// The hexKey can optionally include the "ed25519-priv-" prefix and/or "0x" prefix.
// Returns an error if the key format is invalid or cannot be parsed.
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
		return Account{}, errors.New("can't extract account public key from account private key")
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

// GetAccountAddressString returns the hexadecimal string representation of the account address.
func (a Account) GetAccountAddressString() string {
	return hex.EncodeToString(a.AccountAddress[:])
}

// NewAccountAddress parses a hexadecimal address string and returns a 32-byte address.
// The address can optionally include the "0x" prefix.
// Returns an error if the address format is invalid or exceeds 32 bytes.
func NewAccountAddress(address string) ([32]byte, error) {
	address = strings.TrimPrefix(address, keyPrefix)
	bytes, err := hex.DecodeString(address)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "can't decode account address")
	}
	if len(bytes) > 32 {
		return [32]byte{}, errors.New("account address too long")
	}
	var buf [32]byte
	copy((buf)[32-len(bytes):], bytes)

	return buf, nil
}

func NewResourceAccount(accountAddress [32]byte, seed []byte) ([32]byte, error) {
	var buf [32]byte

	if len(accountAddress) != 32 {
		return buf, errors.New("account address must be 32 bytes")
	}

	data := make([]byte, 0, 32+len(seed)+1)
	data = append(data, accountAddress[:]...)
	data = append(data, seed...)
	data = append(data, deriveResourceAccountSchema)

	hash := sha3.New256()
	hash.Write(data)
	copy(buf[:], hash.Sum(nil))

	return buf, nil
}
