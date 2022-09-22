package provider

import "github.com/klever-io/klever-go-sdk/models"

type KleverChain interface {
	GetAccount(address string) (*models.Account, error)
	GetAsset(assetID string) (*models.KDAData, error)
	Decode(tx *models.Transaction) (*models.TransactionAPI, error)
}
