package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cedra-labs/cedra-go-tx-pablisher"
	"github.com/spf13/cast"
)

const (
	// sender/deployer private key.
	privateKey = "your_private_key"
)

func main() {
	cedraClient := cedra.NewCedraClient(cedra.TestnetChainID)
	sender, err := cedra.NewAccount(privateKey)
	if err != nil {
		panic(err)
	}

	metaBytes, err := os.ReadFile("/path-to-package-build/package-metadata.bcs")
	if err != nil {
		log.Fatalf("failed to read metadata: %v", err)
	}

	files, err := filepath.Glob("/path-to-package-build/bytecode_modules/*.mv")
	if err != nil {
		log.Fatalf("glob modules: %v", err)
	}

	var modules [][]byte

	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			log.Fatalf("read module %s: %v", f, err)
		}

		modules = append(modules, data)
	}

	serializedMeta := cedra.EnncodeToBCSBytes(metaBytes)
	serializedModules := ModuleList(modules).EncodeModule()
	serializedSeed := cedra.EnncodeToBCSBytes([]byte(""))

	payload := cedra.TransactionPayload{
		ModuleAddress: sender.AccountAddress,
		ModuleName:    "deployer",
		FunctionName:  "deploy_derived",
		Argumments: [][]byte{
			serializedMeta,
			serializedModules,
			serializedSeed,
		},
	}

	rawTx, err := cedraClient.NewTransaction(sender, payload)
	if err != nil {
		panic(err)
	}
	encodedTx, auth := rawTx.Sign()

	hash, err := cedraClient.SubmitTransaction(encodedTx, auth)
	if err != nil {
		panic(err)
	}

	fmt.Println(hash)
}

type ModuleList [][]byte

func (m ModuleList) EncodeModule() []byte {
	ser := cedra.NewBCSEncoder()
	ser.EncodeEnum(cast.ToUint64(len(m))) // Length of the outer vector
	for _, module := range m {
		ser.EnncodeBytes(module)
	}

	return ser.GetBytes()
}
