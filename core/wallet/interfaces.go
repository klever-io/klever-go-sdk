package wallet

import "github.com/klever-io/klever-go-sdk/core/account"

type Wallet interface {
	PrivateKey() []byte
	PublicKey() []byte
	GetAccount() (account.Account, error)
}
