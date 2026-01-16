package cedra

import (
	"github.com/pkg/errors"
)

const (
	// DevnetChainID is the chain identifier for the development network.
	DevnetChainID ChainID = 3
	// TestnetChainID is the chain identifier for the test network.
	TestnetChainID ChainID = 2
	// MainnetChainID is the chain identifier for the main network.
	MainnetChainID ChainID = 1
)

// ChainID represents a blockchain network identifier.
type ChainID uint8

// NewLocalnetChainID creates a new chain ID for a local network.
// Panics if the provided ID is 0, as chain IDs must be greater than 0.
func NewLocalnetChainID(id uint8) ChainID {
	if id == 0 {
		panic(errors.New("can't create new chain id: chain id should be greater than 0"))
	}

	return ChainID(id)
}

// Chain represents the configuration for a Cedra blockchain network.
type Chain struct {
	// CedraNodeUrl is the base URL for the Cedra node API.
	CedraNodeUrl string
	// ChainID is the identifier for this chain.
	ChainID ChainID
}

// ChainConfig is a map of chain IDs to their corresponding chain configurations.
type ChainConfig map[ChainID]Chain

// CedraChains contains the predefined chain configurations for devnet, testnet, and mainnet.
var CedraChains = ChainConfig{
	DevnetChainID: {
		ChainID:      DevnetChainID,
		CedraNodeUrl: "https://devnet.cedra.dev/v1/",
	},
	TestnetChainID: {
		ChainID:      TestnetChainID,
		CedraNodeUrl: "https://testnet.cedra.dev/v1/",
	},
	MainnetChainID: {
		ChainID:      MainnetChainID,
		CedraNodeUrl: "https://mainnet.cedra.dev/v1/",
	},
}
