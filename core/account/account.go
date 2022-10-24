package account

import (
	"time"

	"github.com/klever-io/klever-go-sdk/core/address"
	"github.com/klever-io/klever-go-sdk/models"
	"github.com/klever-io/klever-go-sdk/provider"
)

type account struct {
	address address.Address
	info    *models.Account

	lastUpdate time.Time
}

func NewAccount(addr address.Address) (Account, error) {

	return &account{address: addr, info: &models.Account{}}, nil
}

func (a *account) Address() address.Address {
	return a.address
}

func (a *account) Balance() int64 {
	if a.info != nil {
		return a.info.Balance
	}

	return 0
}

func (a *account) Nonce() uint64 {
	if a.info != nil {
		return a.info.Nonce
	}

	return 0
}

func (a *account) IncrementNonce() {
	a.info.Nonce += 1
}

func (a *account) Sync(p provider.KleverChain) error {
	acc, err := p.GetAccount(a.address.Bech32())
	if err != nil {
		return err
	}

	a.info = acc
	a.lastUpdate = time.Now()

	return nil
}

func (a *account) LastUpdate() time.Time {
	return a.lastUpdate
}

func (a *account) GetInfo() *models.Account {
	return a.info
}

func (a *account) NewBaseTX() *models.BaseTX {
	return &models.BaseTX{
		FromAddress: a.address.Bech32(),
		Nonce:       a.Nonce(),
		PermID:      0, // Default owner
		Message:     make([]string, 0),
	}
}
