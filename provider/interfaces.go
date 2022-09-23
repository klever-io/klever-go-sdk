package provider

import "github.com/klever-io/klever-go-sdk/models"

type KleverChain interface {
	// Query Account data
	GetAccount(address string) (*models.Account, error)
	GetAccountAllowance(address string, kda string) (*models.AccountAllowance, error)
	GetAsset(assetID string) (*models.KDAData, error)
	// Transaction helpers
	Decode(tx *models.Transaction) (*models.TransactionAPI, error)
	GetTransaction(hash string) (*models.TransactionAPI, error)
	// Transfer actions
	Send(base *models.BaseTX, toAddr string, amount float64, kda string) (*models.Transaction, error)
	MultiTransfer(base *models.BaseTX, kda string, values []models.ToAmount) (*models.Transaction, error)
	// Acctount Actions
	SetAccountName(base *models.BaseTX, name string) (*models.Transaction, error)
	SetPermission(base *models.BaseTX, permissions []models.PermissionTXRequest) (*models.Transaction, error)
	// Governance Actions
	Proposal(base *models.BaseTX, description string, parameters map[int32]string, duration uint32) (*models.Transaction, error)
	Vote(base *models.BaseTX, proposalID uint64, amount float64, voteType uint64) (*models.Transaction, error)
	// Market&ITO Actions
	ConfigITO(base *models.BaseTX, kdaID, receiverAddress string, status int32, maxAmount float64, packs []models.ParsedPack) (*models.Transaction, error)
	SetITOPrices(base *models.BaseTX, kdaID string, packs []models.ParsedPack) (*models.Transaction, error)
	CreateMarketplace(base *models.BaseTX, name, referralAddr string, referralPercent float64) (*models.Transaction, error)
	ConfigMarketplace(base *models.BaseTX, id, name, referralAddr string, referralPercent float64) (*models.Transaction, error)
	BuyOrder(base *models.BaseTX, id, currency string, amount float64, buyType int32) (*models.Transaction, error)
	SellOrder(base *models.BaseTX, kdaID, currency, mktID string, price, reservePrice float64, endTime int64, mktType int32, message string) (*models.Transaction, error)
	CancelMarketOrder(base *models.BaseTX, orderID string) (*models.Transaction, error)
	Withdraw(base *models.BaseTX, kda string) (*models.Transaction, error)
	// Staking Action
	Freeze(base *models.BaseTX, amount float64, kda string) (*models.Transaction, error)
	Unfreeze(base *models.BaseTX, bucketId, kda string) (*models.Transaction, error)
	Delegate(base *models.BaseTX, toAddr, bucketId string) (*models.Transaction, error)
	Undelegate(base *models.BaseTX, toAddr, bucketId string) (*models.Transaction, error)
	// Validator Actions
	Unjail(base *models.BaseTX) (*models.Transaction, error)
	Claim(base *models.BaseTX, id string, claimType int32) (*models.Transaction, error)
}
