package account

import (
	"time"

	"github.com/klever-io/klever-go-sdk/core/address"
	"github.com/klever-io/klever-go-sdk/models"
	"github.com/klever-io/klever-go-sdk/provider"
)

type Account interface {
	Address() address.Address
	Balance() int64
	Nonce() uint64
	IncrementNonce()
	Sync(provider.KleverChain) error
	LastUpdate() time.Time
	GetInfo() *models.Account
	NewBaseTX() *models.BaseTX
}
