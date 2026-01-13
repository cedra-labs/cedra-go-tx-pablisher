package cedra

type TransactionPayload struct {
	ModuleAddress [32]byte
	ModuleName    string
	FunctionName  string
	Argumments    [][]byte
}
