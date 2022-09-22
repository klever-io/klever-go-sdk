package models

type BaseTX struct {
	FromAddress string
	Nonce       uint64
	PermID      int32
	Message     []string
}

// SendTXRequest -
type SendTXRequest struct {
	Type      uint32        `form:"type" json:"type"`
	Sender    string        `form:"sender" json:"sender"`
	Nonce     uint64        `form:"nonce" json:"nonce"`
	PermID    int32         `form:"permID" json:"permID"`
	Data      [][]byte      `form:"data" json:"data"`
	Contract  interface{}   `form:"contract" json:"contract"`
	Contracts []interface{} `form:"contracts" json:"contracts"`
}

// TransferTXRequest -
type TransferTXRequest struct {
	Receiver string `form:"receiver" json:"receiver"`
	Amount   int64  `form:"amount" json:"amount"`
	KDA      string `form:"kda" json:"kda"`
}

// CreateAssetTXRequest -
type CreateAssetTXRequest struct {
	Type          uint32            `form:"type" json:"type"`
	Name          string            `form:"name" json:"name"`
	Ticker        string            `form:"ticker" json:"ticker"`
	OwnerAddress  string            `form:"ownerAddress" json:"ownerAddress"`
	Logo          string            `form:"logo" json:"logo"`
	URIs          map[string]string `form:"uris" json:"uris"`
	Precision     uint32            `form:"precision" json:"precision"`
	InitialSupply int64             `form:"initialSupply" json:"initialSupply"`
	MaxSupply     int64             `form:"maxSupply" json:"maxSupply"`
	Royalties     *RoyaltiesInfo    `form:"royalties" json:"royalties"`
	Properties    *PropertiesInfo   `form:"properties" json:"properties"`
	Attributes    *AttributesInfo   `form:"attributes" json:"attributes"`
	Staking       *StakingInfo      `form:"staking" json:"staking"`
	Roles         []*RolesInfo      `form:"roles" json:"roles"`
}

// RoyaltyDataTX -
type RoyaltyDataTX struct {
	Amount     float64 `form:"amount" json:"amount"`
	Percentage float64 `form:"percentage" json:"percentage"`
}

// StakingInfo -
type StakingInfo struct {
	APR                 uint32 `form:"apr" json:"apr"`
	MinEpochsToClaim    uint32 `form:"minEpochsToClaim" json:"minEpochsToClaim"`
	MinEpochsToUnstake  uint32 `form:"minEpochsToUnstake" json:"minEpochsToUnstake"`
	MinEpochsToWithdraw uint32 `form:"minEpochsToWithdraw" json:"minEpochsToWithdraw"`
}

// AssetTriggerTXRequest -
type AssetTriggerTXRequest struct {
	TriggerType uint32            `form:"triggerType" json:"triggerType"`
	AssetID     string            `form:"assetId" json:"assetId"`
	Receiver    string            `form:"receiver" json:"receiver"`
	Amount      int64             `form:"amount" json:"amount"`
	MIME        string            `form:"mime" json:"mime"`
	Logo        string            `form:"logo" json:"logo"`
	URIs        map[string]string `form:"uris" json:"uris"`
	Role        *RolesInfo        `form:"role" json:"role"`
	Staking     *StakingInfo      `form:"staking" json:"staking"`
}

// CreateValidatorTXRequest -
type CreateValidatorTXRequest struct {
	BLSPublicKey        string            `form:"blsPublicKey" json:"blsPublicKey"`
	OwnerAddress        string            `form:"ownerAddress" json:"ownerAddress"`
	RewardAddress       string            `form:"rewardAddress" json:"rewardAddress"`
	CanDelegate         bool              `form:"canDelegate" json:"canDelegate"`
	Commission          uint32            `form:"commission" json:"commission"`
	MaxDelegationAmount int64             `form:"maxDelegationAmount" json:"maxDelegationAmount"`
	Logo                string            `form:"logo" json:"logo"`
	URIs                map[string]string `form:"uris" json:"uris"`
	Name                string            `form:"name" json:"name"`
}

// ValidatorConfigTXRequest -
type ValidatorConfigTXRequest struct {
	BLSPublicKey        string            `form:"blsPublicKey" json:"blsPublicKey"`
	RewardAddress       string            `form:"rewardAddress" json:"rewardAddress"`
	CanDelegate         bool              `form:"canDelegate" json:"canDelegate"`
	Commission          uint32            `form:"commission" json:"commission"`
	MaxDelegationAmount int64             `form:"maxDelegationAmount" json:"maxDelegationAmount"`
	Logo                string            `form:"logo" json:"logo"`
	URIs                map[string]string `form:"uris" json:"uris"`
	Name                string            `form:"name" json:"name"`
}

// FreezeTXRequest -
type FreezeTXRequest struct {
	Amount int64  `form:"amount" json:"amount"`
	KDA    string `form:"kda" json:"kda"`
}

// UnfreezeTXRequest -
type UnfreezeTXRequest struct {
	KDA      string `form:"kda" json:"kda"`
	BucketID string `form:"bucketId" json:"bucketId"`
}

// DelegateTXRequest -
type DelegateTXRequest struct {
	Receiver string `form:"receiver" json:"receiver"`
	BucketID string `form:"bucketId" json:"bucketId"`
}

// UndelegateTXRequest -
type UndelegateTXRequest struct {
	BucketID string `form:"bucketId" json:"bucketId"`
}

// WithdrawTXRequest -
type WithdrawTXRequest struct {
	KDA string `form:"kda" json:"kda"`
}

// ClaimTXRequest -
type ClaimTXRequest struct {
	ClaimType int32  `form:"claimType" json:"claimType"`
	ID        string `form:"id" json:"id"`
}

// UnjailTXRequest -
type UnjailTXRequest struct{}

// SetAccountNameTXRequest -
type SetAccountNameTXRequest struct {
	Name string `form:"name" json:"name"`
}

// ProposalTXRequest -
type ProposalTXRequest struct {
	Parameters     map[int32]string `form:"parameters" json:"parameters"`
	Description    string           `form:"description" json:"description"`
	EpochsDuration uint32           `form:"epochsDuration" json:"epochsDuration"`
}

// VoteTXRequest -
type VoteTXRequest struct {
	Type       uint32 `form:"type" json:"type"`
	ProposalID uint64 `form:"proposalId" json:"proposalId"`
	Amount     int64  `form:"amount" json:"amount"`
}

// ConfigITOTXRequest -
type ConfigITOTXRequest struct {
	KDA             string                     `form:"kda" json:"kda"`
	ReceiverAddress string                     `form:"receiverAddress" json:"receiverAddress"`
	Status          int32                      `form:"status" json:"status"`
	MaxAmount       int64                      `form:"maxAmount" json:"maxAmount"`
	PackInfo        map[string]PackInfoRequest `form:"packInfo" json:"packInfo"`
}

type SetITOPricesTXRequest struct {
	KDA      string                     `form:"kda" json:"kda"`
	PackInfo map[string]PackInfoRequest `form:"packInfo" json:"packInfo"`
}

// PackInfoRequest -
type PackInfoRequest struct {
	Packs []PackItemRequest `form:"packs" json:"packs"`
}

// PackItemRequest -
type PackItemRequest struct {
	Amount int64 `form:"amount" json:"amount"`
	Price  int64 `form:"price" json:"price"`
}

// SetITOPricesTXRequest -

// BuyTXRequest -
type BuyTXRequest struct {
	BuyType    int32  `form:"buyType" json:"buyType"`
	ID         string `form:"id" json:"id"`
	CurrencyID string `form:"currencyId" json:"currencyId"`
	Amount     int64  `form:"amount" json:"amount"`
}

// SellTXRequest -
type SellTXRequest struct {
	MarketType    int32  `form:"marketType" json:"marketType"`
	MarketplaceID string `form:"marketplaceId" json:"marketplaceId"`
	AssetID       string `form:"assetId" json:"assetId"`
	CurrencyID    string `form:"currencyId" json:"currencyId"`
	Price         int64  `form:"price" json:"price"`
	ReservePrice  int64  `form:"reservePrice" json:"reservePrice"`
	EndTime       int64  `form:"endTime" json:"endTime"`
}

// CancelMarketOrderTXRequest -
type CancelMarketOrderTXRequest struct {
	OrderID string `form:"orderId" json:"orderId"`
}

// CreateMarketplaceTXRequest -
type CreateMarketplaceTXRequest struct {
	Name               string `form:"name" json:"name"`
	ReferralAddress    string `form:"referralAddress" json:"referralAddress"`
	ReferralPercentage uint32 `form:"referralPercentage" json:"referralPercentage"`
}

// ConfigMarketplaceTXRequest -
type ConfigMarketplaceTXRequest struct {
	MarketplaceID      string `form:"marketplaceId" json:"marketplaceId"`
	Name               string `form:"name" json:"name"`
	ReferralAddress    string `form:"referralAddress" json:"referralAddress"`
	ReferralPercentage uint32 `form:"referralPercentage" json:"referralPercentage"`
}

// SignerTXRequest -
type SignerTXRequest struct {
	Address string `form:"address" json:"address"`
	Weight  int64  `form:"weight" json:"weight"`
}

// PermissionTXRequest -
type PermissionTXRequest struct {
	Type           int32             `form:"type" json:"type"`
	PermissionName string            `form:"permissionName" json:"permissionName"`
	Threshold      int64             `form:"threshold" json:"threshold"`
	Operations     string            `form:"operations" json:"operations"`
	Signers        []SignerTXRequest `form:"signers" json:"signers"`
}

// UpdateAccountPermissionTXRequest -
type UpdateAccountPermissionTXRequest struct {
	Permissions []PermissionTXRequest `form:"permissions" json:"permissions"`
}
