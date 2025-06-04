package testsCommon

import (
	"context"

	"github.com/klever-io/klever-go-sdk/core/address"
	"github.com/klever-io/klever-go-sdk/models"
	"github.com/klever-io/klever-go/data/transaction"
)

// ProxyStub -
type ProxyStub struct {
	GetNetworkConfigCalled func() (*models.NetworkConfig, error)
	SendTransactionCalled  func(transaction *transaction.Transaction) (string, error)
	SendTransactionsCalled func(txs []*transaction.Transaction) ([]string, error)
	GetAccountCalled       func(address address.Address) (*models.Account, error)
}

// GetNetworkConfig -
func (eps *ProxyStub) GetNetworkConfig(_ context.Context) (*models.NetworkConfig, error) {
	if eps.GetNetworkConfigCalled != nil {
		return eps.GetNetworkConfigCalled()
	}

	return &models.NetworkConfig{}, nil
}

// SendTransaction -
func (eps *ProxyStub) SendTransaction(_ context.Context, transaction *transaction.Transaction) (string, error) {
	if eps.SendTransactionCalled != nil {
		return eps.SendTransactionCalled(transaction)
	}

	return "", nil
}

// SendTransactions -
func (eps *ProxyStub) SendTransactions(_ context.Context, txs []*transaction.Transaction) ([]string, error) {
	if eps.SendTransactionCalled != nil {
		return eps.SendTransactionsCalled(txs)
	}

	return make([]string, 0), nil
}

// GetAccount -
func (eps *ProxyStub) GetAccount(_ context.Context, address address.Address) (*models.Account, error) {
	if eps.GetAccountCalled != nil {
		return eps.GetAccountCalled(address)
	}

	return &models.Account{}, nil
}

// IsInterfaceNil -
func (eps *ProxyStub) IsInterfaceNil() bool {
	return eps == nil
}
