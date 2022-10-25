package wallet_test

import (
	"encoding/hex"
	"testing"

	"github.com/klever-io/klever-go-sdk/core/wallet"
	"github.com/stretchr/testify/assert"
)

func TestWallet_WalletFromPem(t *testing.T) {
	fileName := tempPemFile()
	wallet, err := wallet.NewWalletFromPEM(fileName)
	assert.Nil(t, err)

	assert.Equal(t, "8734062c1158f26a3ca8a4a0da87b527a7c168653f7f4c77045e5cf571497d9d", hex.EncodeToString(wallet.PrivateKey()))
	assert.Equal(t, "e41b323a571fd955e09cd41660ff4465c3f44693c87f2faea4a0fc408727c8ea", hex.EncodeToString(wallet.PublicKey()))
}

func TestWallet_Mnemonic(t *testing.T) {
	wallet, err := wallet.NewWalletFromMnemonic("abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about")
	assert.Nil(t, err)

	acc, err := wallet.GetAccount()
	assert.Nil(t, err)
	assert.Equal(t, acc.Address().Bech32(), "klv1usdnywjhrlv4tcyu6stxpl6yvhplg35nepljlt4y5r7yppe8er4qujlazy")
}

func TestWallet_PrivateKey(t *testing.T) {
	wallet, err := wallet.NewWalletFroHex("8734062c1158f26a3ca8a4a0da87b527a7c168653f7f4c77045e5cf571497d9d")
	assert.Nil(t, err)

	acc, err := wallet.GetAccount()
	assert.Nil(t, err)
	assert.Equal(t, acc.Address().Bech32(), "klv1usdnywjhrlv4tcyu6stxpl6yvhplg35nepljlt4y5r7yppe8er4qujlazy")
}

func TestWallet_Sign(t *testing.T) {
	wallet, err := wallet.NewWalletFroHex("8734062c1158f26a3ca8a4a0da87b527a7c168653f7f4c77045e5cf571497d9d")
	assert.Nil(t, err)

	signature, err := wallet.Sign(make([]byte, 32))
	assert.Nil(t, err)
	assert.Equal(t, "e0cb52c1b61916594daefe04aed24620058204a77e95bbbe5bb4ab165a76eda7d0f316d14d10ef184738f149bb5c0bc96ff4a194c25342dc1f6735a35bda0708", hex.EncodeToString(signature))

	hexSign, err := wallet.SignHex("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Nil(t, err)
	assert.Equal(t, signature, hexSign)
}
