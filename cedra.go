package cedra

import "errors"

type CedraClient struct {
	chain Chain
}

func NewCedraClient(chainID ChainID) CedraClient {
	if CedraChains == nil {
		panic(errors.New("")) // TODO:
	}

	chain, ok := CedraChains[chainID]
	if !ok {
		panic(errors.New("")) // TODO:
	}

	return CedraClient{
		chain: chain,
	}
}

func (c CedraClient) SubmitTransaction(tx Transaction) error {
	return nil
}
