package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/klever-io/klever-go-sdk/core"
	"github.com/klever-io/klever-go-sdk/models"
	"github.com/klever-io/klever-go-sdk/models/proto"
)

func (kc *kleverChain) DecodeWithContext(ctx context.Context, tx *proto.Transaction) (*models.TransactionAPI, error) {

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

	err = kc.httpClient.Post(ctx, fmt.Sprintf("%s/transaction/decode", kc.networkConfig.GetNodeUri()), string(body), nil, &result)

	return result.Data.Transaction, err
}

func (kc *kleverChain) Decode(tx *proto.Transaction) (*models.TransactionAPI, error) {
	return kc.DecodeWithContext(context.Background(), tx)
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

func (kc *kleverChain) Send(base *models.BaseTX, toAddr string, amount float64, kda string) (*proto.Transaction, error) {
	values := []models.ToAmount{{ToAddress: toAddr, Amount: amount, KDA: kda}}

	return kc.MultiTransfer(base, values)
}

func (kc *kleverChain) MultiSend(base *models.BaseTX, contracts []models.AnyContractRequest) (*proto.Transaction, error) {
	var contractsParsed []interface{}

	if len(contracts) > 20 {
		return nil, errors.New("max contracts reached")
	}

	for _, contract := range contracts {
		c, err := contract.PrepareToSend()
		if err != nil {
			return nil, err
		}

		contractsParsed = append(contractsParsed, c)
	}

	data, err := kc.buildRequest(0, base, contractsParsed)
	if err != nil {
		return nil, err
	}

	return kc.PrepareTransaction(data)
}

func (kc *kleverChain) MultiTransfer(base *models.BaseTX, values []models.ToAmount) (*proto.Transaction, error) {
	contracts := make([]interface{}, 0)
	for _, to := range values {
		precision, err := kc.getPrecision(to.KDA)
		if err != nil {
			return nil, err
		}

		parsedAmount := to.Amount * math.Pow10(int(precision))
		contracts = append(contracts, models.TransferTXRequest{
			Receiver: to.ToAddress,
			Amount:   int64(parsedAmount),
			KDA:      to.KDA,
		})
	}

	data, err := kc.buildRequest(proto.TXContract_TransferContractType, base, contracts)
	if err != nil {
		return nil, err
	}
	return kc.PrepareTransaction(data)
}

func (kc *kleverChain) buildRequest(
	txType proto.TXContract_ContractType,
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
		KDAFee:    base.KdaFee,
	}, nil
}

func (kc *kleverChain) PrepareTransaction(request *models.SendTXRequest) (*proto.Transaction, error) {
	return kc.PrepareTransactionWithContext(context.Background(), request)
}

func (kc *kleverChain) PrepareTransactionWithContext(ctx context.Context, request *models.SendTXRequest) (*proto.Transaction, error) {
	result := struct {
		Data struct {
			Transaction *proto.Transaction `json:"result"`
		} `json:"data"`
	}{}

	result.Data.Transaction = &proto.Transaction{}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	err = kc.httpClient.Post(ctx, fmt.Sprintf("%s/transaction/send", kc.networkConfig.GetNodeUri()), string(body), nil, &result)
	if err == nil {
		hash, err := kc.CalculateHash(result.Data.Transaction.RawData)
		if err == nil {
			result.Data.Transaction.Hash = hash
		}
	}

	return result.Data.Transaction, err
}

// CalculateHash marshalizes the interface and calculates its hash
func (kc *kleverChain) CalculateHash(
	object interface{},
) ([]byte, error) {

	mrsData, err := kc.marshalizer.Marshal(object)
	if err != nil {
		return nil, err
	}

	hash := kc.hasher.Compute(string(mrsData))
	return hash, nil
}

func (kc *kleverChain) GetTransactionWithContext(ctx context.Context, hash string) (*models.TransactionAPI, error) {
	result := struct {
		Data struct {
			Transaction *models.TransactionAPI `json:"transaction"`
		} `json:"data"`
	}{}

	result.Data.Transaction = &models.TransactionAPI{}

	err := kc.httpClient.Get(ctx, fmt.Sprintf("%s/transaction/%s", kc.networkConfig.GetAPIUri(), hash), &result)

	return result.Data.Transaction, err
}

func (kc *kleverChain) GetTransaction(hash string) (*models.TransactionAPI, error) {
	return kc.GetTransactionWithContext(context.Background(), hash)
}

func (kc *kleverChain) BroadcastTransactionWithContext(ctx context.Context, tx *proto.Transaction) (string, error) {
	toBroadcast := struct {
		TX *proto.Transaction `json:"tx"`
	}{
		TX: tx,
	}

	data, err := json.Marshal(toBroadcast)
	if err != nil {
		return "", err
	}

	result := struct {
		Data struct {
			TXCount int    `json:"txCount"`
			TXHash  string `json:"txHash"`
		} `json:"data"`
		Error string `json:"error"`
		Code  string `json:"code"`
	}{}

	err = kc.httpClient.Post(ctx, fmt.Sprintf("%s/transaction/broadcast", kc.networkConfig.GetNodeUri()), string(data), nil, &result)
	if err != nil {
		return "", err
	}

	if len(result.Error) != 0 {
		return "", fmt.Errorf("error broadcasting transcation: %s", result.Error)
	}

	return result.Data.TXHash, err
}

func (kc *kleverChain) BroadcastTransaction(tx *proto.Transaction) (string, error) {
	return kc.BroadcastTransactionWithContext(context.Background(), tx)
}

func (kc *kleverChain) BroadcastTransactions(txs []*proto.Transaction) ([]string, error) {
	return kc.BroadcastTransactionsWithContext(context.Background(), txs)
}

func (kc *kleverChain) BroadcastTransactionsWithContext(ctx context.Context, txs []*proto.Transaction) ([]string, error) {
	toBroadcast := struct {
		TXs []*proto.Transaction `json:"txs"`
	}{
		TXs: txs,
	}

	data, err := json.Marshal(toBroadcast)
	if err != nil {
		return nil, err
	}

	result := struct {
		Data struct {
			TxsHashes []string `json:"txsHashes"`
		} `json:"data"`
		Error string `json:"error"`
		Code  string `json:"code"`
	}{}

	err = kc.httpClient.Post(ctx, fmt.Sprintf("%s/transaction/broadcast", kc.networkConfig.GetNodeUri()), string(data), nil, &result)
	if err != nil {
		return nil, err
	}

	if len(result.Error) != 0 {
		return nil, fmt.Errorf("error broadcasting transcations: %s", result.Error)
	}

	return result.Data.TxsHashes, nil
}
