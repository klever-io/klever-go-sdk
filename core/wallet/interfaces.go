package wallet

import "github.com/klever-io/klever-go-sdk/core/account"

type Wallet interface {
	PrivateKey() []byte
	PublicKey() []byte
	GetAccount() (account.Account, error)
	Sign(msg []byte) ([]byte, error)
	SignHex(msg string) ([]byte, error)
}
