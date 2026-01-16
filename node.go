package cedra

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

const (
	clientHeader                 = "X-Cedra-Client"
	clientHeaderValue            = "cedra-tx-publisher"
	contentTypeAptosSignedTxnBcs = "application/x.cedra.signed_transaction+bcs"
)

type CedraNode struct {
	chain   Chain
	nodeURL url.URL
}

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
	}
}

func (n CedraNode) SubmitTransaction(tx []byte) (string, error) {
	requstBody := bytes.NewReader(tx)
	requestURL := n.nodeURL.JoinPath("transactions")
	headers := map[string]string{
		"content-type": contentTypeAptosSignedTxnBcs,
	}

	hash, err := makeRequest[TransactionDTO](http.MethodPost, *requestURL, requstBody, headers)
	if err != nil {
		return "", errors.Wrap(err, "can't execute requested transaction")
	}
	return hash.Hash, nil
}

func (n CedraNode) GetEstimateGasPrice() (EstimateGasPriceDTO, error) {
	var requstBody io.Reader
	var headers map[string]string

	requestURL := n.nodeURL.JoinPath("estimate_gas_price")
	estimateGasPrice, err := makeRequest[EstimateGasPriceDTO](http.MethodGet, *requestURL, requstBody, headers)
	if err != nil {
		return estimateGasPrice, errors.Wrap(err, "can't estimate gas price")
	}

	return estimateGasPrice, nil
}

func (n CedraNode) GetSequenceNumber(address string) (uint64, error) {
	var requstBody io.Reader
	var headers map[string]string

	requestURL := n.nodeURL.JoinPath("accounts", address)

	accountInfo, err := makeRequest[AccountDTO](http.MethodGet, *requestURL, requstBody, headers)
	if err != nil {
		return 0, errors.Wrap(err, "can't get account info")
	}

	return cast.ToUint64(accountInfo.SequenceNumber), nil
}

func makeRequest[T any](method string, requetURL url.URL, body io.Reader, headers map[string]string) (T, error) {
	var response T
	req, err := http.NewRequest(method, requetURL.String(), body)
	if err != nil {
		return response, errors.Wrap(err, "can't create a new request")
	}

	req.Header.Set(clientHeader, clientHeaderValue)
	for header, value := range headers {
		req.Header.Set(header, value)
	}

	resp, err := http.DefaultClient.Do(req)
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

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return response, errors.Wrap(err, "can't unarshal response to object")
	}

	return response, nil
}
