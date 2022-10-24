package models_test

import (
	"testing"

	"github.com/klever-io/klever-go-sdk/core/wallet"
	"github.com/klever-io/klever-go-sdk/models"
	"github.com/stretchr/testify/assert"
)

func TestSignableMessage_SignVerify(t *testing.T) {
	w, _ := wallet.NewWalletFroHex("0000000000000000000000000000000000000000000000000000000000000000")
	acc, _ := w.GetAccount()

	sm := models.NewSM(acc.Address().Bech32(), []byte("KleverGames Test Message"))

	assert.False(t, sm.Verify())

	err := sm.Sign(w)
	assert.Nil(t, err)
	assert.True(t, sm.Verify())
}
