package account

import (
	"github.com/klever-io/klever-go-sdk/core/address"
	"github.com/klever-io/klever-go-sdk/provider"
)

type account struct {
	address address.Address
	balance int64
	nonce   uint64
}

func NewAccount(addr address.Address) (Account, error) {

	return &account{address: addr, balance: 0, nonce: 0}, nil
}

func (a *account) Address() address.Address {
	return a.address
}

func (a *account) Balance() int64 {
	return a.balance
}

func (a *account) Nonce() uint64 {
	return a.nonce
}

func (a *account) IncrementNonce() {
	a.nonce += 1
}

func (a *account) Sync(p provider.KleverChain) error {
	acc, err := p.GetAccount(a.address.Bech32())
	if err != nil {
		return err
	}

	a.balance = acc.Balance
	a.nonce = acc.Nonce

	return nil
}
