package cedra

type AccountDTO struct {
	SequenceNumber    string `json:"sequence_number"`
	AuthenticationKey string `json:"authentication_key"`
}

type EstimateGasPriceDTO struct {
	DeprioritizedGasEstimate uint64 `json:"deprioritized_gas_estimate"`
	GasEstimate              uint64 `json:"gas_estimate"`
	PrioritizedGasEstimate   uint64 `json:"prioritized_gas_estimate"`
}

type TransactionDTO struct {
	Hash string
}
