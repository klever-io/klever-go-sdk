package provider

import (
	"github.com/klever-io/klever-go-sdk/models"
	"github.com/klever-io/klever-go-sdk/models/proto"
	"github.com/klever-io/klever-go-sdk/provider/tools/hasher"
	"github.com/klever-io/klever-go-sdk/provider/tools/marshal"
)

type KleverChain interface {
	// Query Account data
	GetAccount(address string) (*models.Account, error)
	GetAccountAllowance(address string, kda string) (*models.AccountAllowance, error)
	GetAsset(assetID string) (*proto.KDAData, error)
	// Transaction helpers
	Decode(tx *proto.Transaction) (*models.TransactionAPI, error)
	GetTransaction(hash string) (*models.TransactionAPI, error)
	GetHasher() hasher.Hasher
	GetMarshalizer() marshal.Marshalizer
	// Transfer actions
	Send(base *models.BaseTX, toAddr string, amount float64, kda string) (*proto.Transaction, error)
	MultiTransfer(base *models.BaseTX, values []models.ToAmount) (*proto.Transaction, error)
	// Asset Actions
	CreateKDA(base *models.BaseTX, kdaType proto.KDAData_EnumAssetType, op *models.KDAOptions) (*proto.Transaction, error)
	AssetTrigger(base *models.BaseTX, kdaID string, triggerType AssetTriggerType, op *models.AssetTriggerOptions) (*proto.Transaction, error)
	Deposit(base *models.BaseTX, op *models.DepositOptions) (*proto.Transaction, error)
	Withdraw(base *models.BaseTX, op *models.WithdrawOptions) (*proto.Transaction, error)
	// Acctount Actions
	SetAccountName(base *models.BaseTX, name string) (*proto.Transaction, error)
	SetPermission(base *models.BaseTX, permissions []models.PermissionTXRequest) (*proto.Transaction, error)
	// Governance Actions
	Proposal(base *models.BaseTX, description string, parameters map[int32]string, duration uint32) (*proto.Transaction, error)
	Vote(base *models.BaseTX, proposalID uint64, amount float64, voteType uint64) (*proto.Transaction, error)
	// Market&ITO Actions
	ConfigITO(base *models.BaseTX, kdaID, receiverAddress string, status int32, maxAmount float64, packs []models.ParsedPack) (*proto.Transaction, error)
	SetITOPrices(base *models.BaseTX, kdaID string, packs []models.ParsedPack) (*proto.Transaction, error)
	ITOTrigger(base *models.BaseTX, kdaID string, triggerType ITOTriggerType, op *models.ITOTriggerOptions) (*proto.Transaction, error)
	CreateMarketplace(base *models.BaseTX, name, referralAddr string, referralPercent float64) (*proto.Transaction, error)
	ConfigMarketplace(base *models.BaseTX, id, name, referralAddr string, referralPercent float64) (*proto.Transaction, error)
	BuyOrder(base *models.BaseTX, id, currency string, amount float64, buyType int32) (*proto.Transaction, error)
	SellOrder(base *models.BaseTX, kdaID, currency, mktID string, price, reservePrice float64, endTime int64, mktType int32, message string) (*proto.Transaction, error)
	CancelMarketOrder(base *models.BaseTX, orderID string) (*proto.Transaction, error)
	// Staking Action
	Freeze(base *models.BaseTX, amount float64, kda string) (*proto.Transaction, error)
	Unfreeze(base *models.BaseTX, bucketId, kda string) (*proto.Transaction, error)
	Delegate(base *models.BaseTX, toAddr, bucketId string) (*proto.Transaction, error)
	Undelegate(base *models.BaseTX, toAddr, bucketId string) (*proto.Transaction, error)
	// Validator Actions
	Unjail(base *models.BaseTX) (*proto.Transaction, error)
	Claim(base *models.BaseTX, id string, claimType int32) (*proto.Transaction, error)
	// Multi contract Action
	MultiSend(base *models.BaseTX, contracts ...models.AnyContractRequest) (*proto.Transaction, error)
	// Network Broadcast
	BroadcastTransaction(tx *proto.Transaction) (string, error)
	BroadcastTransactions(txs ...*proto.Transaction) ([]string, error)
}
