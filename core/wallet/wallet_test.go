package wallet_test

import (
	"testing"

	"github.com/klever-io/klever-go-sdk/core/wallet"
	"github.com/stretchr/testify/assert"
)

func TestWallet_Mnemonic(t *testing.T) {
	wallet, err := wallet.NewWalletFromMnemonic("abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about")
	assert.Nil(t, err)

	acc, err := wallet.GetAccount()
	assert.Nil(t, err)
	assert.Equal(t, acc.Address().Bech32(), "klv1usdnywjhrlv4tcyu6stxpl6yvhplg35nepljlt4y5r7yppe8er4qujlazy")
}
