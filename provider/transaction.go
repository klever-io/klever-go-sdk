package provider

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/klever-io/klever-go-sdk/core"
	"github.com/klever-io/klever-go-sdk/models"
)

func (kc *kleverChain) Decode(tx *models.Transaction) (*models.TransactionAPI, error) {

	result := struct {
		Data struct {
			Transaction *models.TransactionAPI `json:"tx"`
		} `json:"data"`
	}{}
	result.Data.Transaction = new(models.TransactionAPI)

	body, err := json.Marshal(tx)
	if err != nil {
		return result.Data.Transaction, nil
	}

	err = kc.httpClient.Post(fmt.Sprintf("%s/transaction/decode", kc.networkConfig.GetAPIUri()), string(body), nil, &result)

	return result.Data.Transaction, err
}

func (kc *kleverChain) getPrecision(kda string) (uint32, error) {
	precision := uint32(6)
	isNFT := false
	if strings.Contains(kda, "/") {
		isNFT = true
		precision = 0
	}

	if !isNFT && len(kda) > 0 && kda != core.KLV && kda != core.KFI {
		asset, err := kc.GetAsset(kda)
		if err != nil {
			return 0, err
		}

		precision = asset.Precision
	}
	return precision, nil
}

func (kc *kleverChain) Send(base *models.BaseTX, toAddr string, amount float64, kda string) (*models.Transaction, error) {
	values := []models.ToAmount{{ToAddress: toAddr, Amount: amount}}

	return kc.MultiTransfer(base, kda, values)
}

func (kc *kleverChain) MultiTransfer(base *models.BaseTX, kda string, values []models.ToAmount) (*models.Transaction, error) {
	precision, err := kc.getPrecision(kda)
	if err != nil {
		return nil, err
	}

	contracts := make([]interface{}, 0)
	for _, to := range values {
		parsedAmount := to.Amount * math.Pow10(int(precision))
		contracts = append(contracts, models.TransferTXRequest{
			Receiver: to.ToAddress,
			Amount:   int64(parsedAmount),
			KDA:      kda,
		})
	}

	data, err := kc.buildRequest(models.TXContract_TransferContractType, base, contracts)
	if err != nil {
		return nil, err
	}
	return kc.PrepareTransaction(data)
}

func (kc *kleverChain) buildRequest(
	txType models.TXContract_ContractType,
	base *models.BaseTX,
	contracts []interface{},
) (*models.SendTXRequest, error) {

	if len(contracts) == 0 || len(contracts) > core.MaxLenghtOfContracts {
		return nil, fmt.Errorf("invalid len of contracts to build request: %d", len(contracts))
	}

	var parsedMessage [][]byte
	for _, m := range base.Message {
		parsedMessage = append(parsedMessage, []byte(m))
	}

	var contract interface{}
	if len(contracts) == 1 {
		contract = contracts[0]
	}

	return &models.SendTXRequest{
		Type:      uint32(txType),
		Sender:    base.FromAddress,
		Nonce:     base.Nonce,
		PermID:    base.PermID,
		Data:      parsedMessage,
		Contract:  contract,
		Contracts: contracts,
	}, nil
}

func (kc *kleverChain) PrepareTransaction(request *models.SendTXRequest) (*models.Transaction, error) {
	result := struct {
		Data struct {
			Transaction *models.Transaction `json:"result"`
		} `json:"data"`
	}{}

	result.Data.Transaction = &models.Transaction{}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	err = kc.httpClient.Post(fmt.Sprintf("%s/transaction/send", kc.networkConfig.GetAPIUri()), string(body), nil, &result)

	return result.Data.Transaction, err
}

func (kc *kleverChain) GetTransaction(hash string) (*models.TransactionAPI, error) {
	result := struct {
		Data struct {
			Transaction *models.TransactionAPI `json:"transaction"`
		} `json:"data"`
	}{}

	result.Data.Transaction = &models.TransactionAPI{}

	err := kc.httpClient.Get(fmt.Sprintf("%s/transaction/%s", kc.networkConfig.GetAPIUri(), hash), &result)

	return result.Data.Transaction, err
}
