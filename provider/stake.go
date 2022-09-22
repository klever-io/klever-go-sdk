package provider

import (
	"math"

	"github.com/klever-io/klever-go-sdk/models"
)

func (kc *kleverChain) Freeze(base *models.BaseTX, amount float64, kda string) (*models.Transaction, error) {
	precision, err := kc.getPrecision(kda)
	if err != nil {
		return nil, err
	}

	parsedAmount := amount * math.Pow10(int(precision))

	contracts := []interface{}{models.FreezeTXRequest{
		Amount: int64(parsedAmount),
		KDA:    kda,
	}}

	data, err := kc.buildRequest(models.TXContract_FreezeContractType, base, contracts)
	if err != nil {
		return nil, err
	}
	return kc.PrepareTransaction(data)
}

func (kc *kleverChain) Unfreeze(base *models.BaseTX, bucketId, kda string) (*models.Transaction, error) {
	contracts := []interface{}{models.UnfreezeTXRequest{
		BucketID: bucketId,
		KDA:      kda,
	}}

	data, err := kc.buildRequest(models.TXContract_UnfreezeContractType, base, contracts)
	if err != nil {
		return nil, err
	}
	return kc.PrepareTransaction(data)
}

func (kc *kleverChain) Delegate(base *models.BaseTX, toAddr, bucketId string) (*models.Transaction, error) {
	contracts := []interface{}{models.DelegateTXRequest{
		Receiver: toAddr,
		BucketID: bucketId,
	}}

	data, err := kc.buildRequest(models.TXContract_DelegateContractType, base, contracts)
	if err != nil {
		return nil, err
	}
	return kc.PrepareTransaction(data)
}

func (kc *kleverChain) Undelegate(base *models.BaseTX, toAddr, bucketId string) (*models.Transaction, error) {
	contracts := []interface{}{models.UndelegateTXRequest{
		BucketID: bucketId,
	}}

	data, err := kc.buildRequest(models.TXContract_UndelegateContractType, base, contracts)
	if err != nil {
		return nil, err
	}
	return kc.PrepareTransaction(data)
}

func (kc *kleverChain) Withdraw(base *models.BaseTX, kda string) (*models.Transaction, error) {
	contracts := []interface{}{models.WithdrawTXRequest{
		KDA: kda,
	}}

	data, err := kc.buildRequest(models.TXContract_WithdrawContractType, base, contracts)
	if err != nil {
		return nil, err
	}
	return kc.PrepareTransaction(data)
}

func (kc *kleverChain) Claim(base *models.BaseTX, id string, claimType int32) (*models.Transaction, error) {
	contracts := []interface{}{models.ClaimTXRequest{
		ClaimType: claimType,
		ID:        id,
	}}

	data, err := kc.buildRequest(models.TXContract_ClaimContractType, base, contracts)
	if err != nil {
		return nil, err
	}
	return kc.PrepareTransaction(data)
}

func (kc *kleverChain) Unjail(base *models.BaseTX) (*models.Transaction, error) {
	contracts := []interface{}{models.UnjailTXRequest{}}

	data, err := kc.buildRequest(models.TXContract_UnjailContractType, base, contracts)
	if err != nil {
		return nil, err
	}
	return kc.PrepareTransaction(data)
}
