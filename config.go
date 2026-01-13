package cedra

import "errors"

const (
	DevnetChainID  ChainID = 3
	TestnetChainID ChainID = 2
	MainnetChainID ChainID = 1
)

type ChainID uint8

func NewLocalnetChainID(id uint8) ChainID {
	if id == 0 {
		panic(errors.New("")) // TODO:
	}

	return ChainID(id)
}

type Chain struct {
	CedraNodeUrl string
	ChainID      ChainID
}

type ChainConfig map[ChainID]Chain

var CedraChains = ChainConfig{
	DevnetChainID: {
		ChainID:      DevnetChainID,
		CedraNodeUrl: "",
	},
	TestnetChainID: {
		ChainID:      TestnetChainID,
		CedraNodeUrl: "",
	},
	MainnetChainID: {
		ChainID:      MainnetChainID,
		CedraNodeUrl: "",
	},
}
