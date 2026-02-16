package cedra

// AccountDTO represents the account information returned from the Cedra node API.
type AccountDTO struct {
	// SequenceNumber is the current sequence number of the account.
	SequenceNumber string `json:"sequence_number"`
	// AuthenticationKey is the authentication key for the account.
	AuthenticationKey string `json:"authentication_key"`
}

// EstimateGasPriceDTO represents the gas price estimates returned from the Cedra node API.
type EstimateGasPriceDTO struct {
	// DeprioritizedGasEstimate is the gas price estimate for deprioritized transactions.
	DeprioritizedGasEstimate uint64 `json:"deprioritized_gas_estimate"`
	// GasEstimate is the standard gas price estimate.
	GasEstimate uint64 `json:"gas_estimate"`
	// PrioritizedGasEstimate is the gas price estimate for prioritized transactions.
	PrioritizedGasEstimate uint64 `json:"prioritized_gas_estimate"`
}

// TransactionDTO represents the transaction response from the Cedra node API.
type TransactionDTO struct {
	// Hash is the transaction hash returned after submission.
	Hash string `json:"hash"`

	VMStatus string `json:"vm_status"`

	TxType string `json:"type"`
}
