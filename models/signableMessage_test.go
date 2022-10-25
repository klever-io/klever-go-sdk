package models_test

import (
	"encoding/hex"
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

	assert.Equal(t, "11c33430efd7288d306c9f69473157f4dcd0074bdc29d136b4f12e58ba051b9da0e5d72869fc673a871a9b81ae70a61f1e44e89ccb82bd2ce26486499355bf05", hex.EncodeToString(sm.GetSignature()))
}

func TestSignableMessage_LoadJson(t *testing.T) {
	org := models.NewSM("klv18d4z00xwk6jz6c4r4rgz5mcdwdjny9thrh3y8f36cpy2rz6emg5scllvd2", []byte("Klever"))
	message := org.ToJSON()

	sm := models.NewSignableMessage()
	err := sm.LoadJSON(message)
	assert.Nil(t, err)

	assert.Equal(t, org.ToJSON(), sm.ToJSON())
}

func TestSignableMessage_LoadJsonFail(t *testing.T) {
	message := "{\"address\":\"klv123\",\"message\":\"0x4b6c6576657247616d65732054657374204d657373616765\",\"signature\":\"0x11c33430efd7288d306c9f69473157f4dcd0074bdc29d136b4f12e58ba051b9da0e5d72869fc673a871a9b81ae70a61f1e44e89ccb82bd2ce26486499355bf05\",\"version\":1,\"signer\":\"KleverGO\"}"

	sm := models.NewSignableMessage()
	err := sm.LoadJSON(message)
	assert.Contains(t, err.Error(), "invalid bech32 string")

	err = sm.LoadJSON("")
	assert.Contains(t, err.Error(), "unexpected end of JSON input")
}
