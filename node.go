package cedra

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

const (
	// clientHeader is the HTTP header name for the client identifier.
	clientHeader = "X-Cedra-Client"
	// clientHeaderValue is the value sent in the client header.
	clientHeaderValue = "cedra-tx-publisher"
	// contentTypeAptosSignedTxnBcs is the content type for signed transaction BCS data.
	contentTypeAptosSignedTxnBcs = "application/x.cedra.signed_transaction+bcs"
	// defaultHTTPTimeout is the default timeout for HTTP requests.
	defaultHTTPTimeout = 30 * time.Second
)

// CedraNode represents a client for communicating with a Cedra blockchain node.
type CedraNode struct {
	// chain contains the chain configuration.
	chain Chain
	// nodeURL is the parsed URL of the Cedra node.
	nodeURL url.URL
	// httpClient is the HTTP client used for making requests to the node.
	httpClient *http.Client
}

// NewCedraNode creates a new CedraNode instance for the specified chain.
// Panics if the chain configuration is invalid or the chain ID doesn't exist.
func NewCedraNode(chainID ChainID) CedraNode {
	if CedraChains == nil {
		panic(errors.New("can't create a new instance of CedraNode: invalid chain config"))
	}

	chain, ok := CedraChains[chainID]
	if !ok {
		panic(errors.New("can't create a new instance of CedraNode: requested chain doesn't exist"))
	}
	nodeURL, err := url.Parse(chain.CedraNodeUrl)
	if err != nil {
		panic(errors.Wrapf(err, "can't create a new instance of CedraNode"))
	}

	return CedraNode{
		chain:   chain,
		nodeURL: *nodeURL,
		httpClient: &http.Client{
			Timeout: defaultHTTPTimeout,
		},
	}
}

// SubmitTransaction submits a signed transaction to the Cedra node.
// Returns the transaction hash if successful, or an error if submission fails.
func (n CedraNode) SubmitTransaction(tx []byte) (string, error) {
	requestBody := bytes.NewReader(tx)
	requestURL := n.nodeURL.JoinPath("transactions")
	headers := map[string]string{
		"content-type": contentTypeAptosSignedTxnBcs,
	}

	hash, err := makeRequest[TransactionDTO](http.MethodPost, requestURL, requestBody, headers, n.httpClient)
	if err != nil {
		return "", errors.Wrap(err, "can't execute requested transaction")
	}

	return hash.Hash, nil
}

// GetEstimateGasPrice retrieves the current gas price estimates from the Cedra node.
// Returns gas price estimates for different priority levels.
func (n CedraNode) GetEstimateGasPrice() (EstimateGasPriceDTO, error) {
	var body io.Reader
	var headers map[string]string
	requestURL := n.nodeURL.JoinPath("estimate_gas_price")
	estimateGasPrice, err := makeRequest[EstimateGasPriceDTO](http.MethodGet, requestURL, body, headers, n.httpClient)
	if err != nil {
		return estimateGasPrice, errors.Wrap(err, "can't estimate gas price")
	}

	return estimateGasPrice, nil
}

// GetSequenceNumber retrieves the current sequence number for the specified account address.
// Returns the sequence number as a uint64, or an error if the request fails.
func (n CedraNode) GetSequenceNumber(address string) (uint64, error) {
	var body io.Reader
	var headers map[string]string
	requestURL := n.nodeURL.JoinPath("accounts", address)
	accountInfo, err := makeRequest[AccountDTO](http.MethodGet, requestURL, body, headers, n.httpClient)
	if err != nil {
		return 0, errors.Wrap(err, "can't get account info")
	}

	return cast.ToUint64(accountInfo.SequenceNumber), nil
}

func (n CedraNode) WaitTxByHash(txHash string) (TransactionDTO, error) {
	var body io.Reader
	var headers map[string]string
	requestURL := n.nodeURL.JoinPath("transactions/wait_by_hash", txHash)

	tx, err := makeRequest[TransactionDTO](http.MethodGet, requestURL, body, headers, n.httpClient)
	if err != nil {
		return TransactionDTO{}, errors.Wrap(err, "can't whait for requested transaction")
	}

	return tx, nil
}

// makeRequest performs an HTTP request to the Cedra node and unmarshals the JSON response.
// It is a generic function that can handle different response types.
// Returns the unmarshaled response or an error if the request fails.
func makeRequest[T any](method string, requestURL *url.URL, body io.Reader, headers map[string]string, client *http.Client) (T, error) {
	var response T
	req, err := http.NewRequest(method, requestURL.String(), body)
	if err != nil {
		return response, errors.Wrap(err, "can't create a new request")
	}

	req.Header.Set(clientHeader, clientHeaderValue)
	for header, value := range headers {
		req.Header.Set(header, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return response, errors.Wrap(err, "can't execute request")
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, errors.Wrap(err, "can't read request response body")
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return response, errors.New(resp.Status + ": " + string(bodyBytes))
	}

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return response, errors.Wrap(err, "can't unmarshal response to object")
	}

	return response, nil
}
