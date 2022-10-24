package models

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/klever-io/klever-go-sdk/core/address"
	"golang.org/x/crypto/sha3"
)

const (
	MESSAGE_PREFIX = "\x17Klever Signed Message:\n"
)

type SignableMessage interface {
	SetAddress(address string) error
	SetMessage(message []byte)

	GetSignature() []byte
	SetSignature([]byte)

	Verify() bool
	Sign(Signer) error

	ToJSON() string
}

type signableMessage struct {
	Address   address.Address
	Message   []byte
	Signature []byte
	Version   int
	Signer    string
}

func NewSignableMessageFrom(address string) SignableMessage {
	sm := NewSignableMessage()
	sm.SetAddress(address)
	return sm
}

func NewSM(address string, message []byte) SignableMessage {
	sm := NewSignableMessageFrom(address)
	sm.SetMessage(message)
	return sm
}

func NewSignableMessage() SignableMessage {
	return &signableMessage{
		Version: 1,
		Signer:  "KleverGO",
	}
}

func (sm *signableMessage) SetAddress(addr string) error {
	encoded, err := address.NewAddress(addr)
	if err != nil {
		return err
	}

	sm.Address = encoded

	return nil
}

func (sm *signableMessage) SetMessage(message []byte) {
	sm.Message = make([]byte, len(message))
	copy(sm.Message, message)
}

func (sm *signableMessage) serializeForSigning() []byte {
	messageSize := fmt.Sprintf("%d", len(sm.Message))
	signableMessage := append([]byte(messageSize), sm.Message...)
	bytesToHash := append([]byte(MESSAGE_PREFIX), signableMessage...)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(bytesToHash)

	return hash.Sum(nil)
}

func (sm *signableMessage) GetSignature() []byte {
	return sm.Signature
}

func (sm *signableMessage) SetSignature(signature []byte) {
	sm.Signature = make([]byte, len(signature))
	copy(sm.Signature, signature)
}

func (sm *signableMessage) Verify() bool {
	hash := sm.serializeForSigning()

	return ed25519.Verify(sm.Address.Bytes(), hash, sm.Signature)
}

func (sm *signableMessage) Sign(wallet Signer) error {
	data := sm.serializeForSigning()
	siganture, err := wallet.Sign(data)
	if err != nil {
		return err
	}

	sm.SetSignature(siganture)

	return nil
}

func (sm *signableMessage) ToJSON() string {
	result := struct {
		Address   string `json:"address"`
		Message   string `json:"message"`
		Signature string `json:"signature"`
		Version   int    `json:"version"`
		Signer    string `json:"signer"`
	}{
		Address:   sm.Address.Bech32(),
		Message:   "0x" + hex.EncodeToString(sm.Message),
		Signature: "0x" + hex.EncodeToString(sm.Signature),
		Version:   sm.Version,
		Signer:    sm.Signer,
	}

	data, _ := json.Marshal(result)
	return string(data)
}
