package models

import (
	"encoding/json"
	"time"

	"github.com/klever-io/klever-go-sdk/models/proto"
)

type ToAmount struct {
	ToAddress string
	Amount    float64
}

type URI struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type TransactionAPI struct {
	Hash         string                   `json:"hash"`
	BlockNum     uint64                   `json:"blockNum,omitempty"`
	Sender       string                   `json:"sender"`
	Nonce        uint64                   `json:"nonce"`
	PermissionID int32                    `json:"permissionID,omitempty"`
	Data         []string                 `json:"data,omitempty"`
	Timestamp    time.Duration            `json:"timestamp,omitempty"`
	KAppFee      int64                    `json:"kAppFee"`
	BandwidthFee int64                    `json:"bandwidthFee"`
	Status       string                   `json:"status"`
	ResultCode   string                   `json:"resultCode,omitempty"`
	Version      uint32                   `json:"version,omitempty"`
	ChainID      string                   `json:"chainID,omitempty"`
	Signature    []string                 `json:"signature,omitempty"`
	SearchOrder  uint32                   `json:"searchOrder"`
	Receipts     []map[string]interface{} `json:"receipts"`
	Contracts    []*proto.TXContract      `json:"contract"`
}

func (t *TransactionAPI) String() string {
	result, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		result = make([]byte, 0)
	}

	return string(result)
}

type TXContractAPI struct {
	Type       proto.TXContract_ContractType `json:"type"`
	TypeString string                        `json:"typeString"`
	Parameter  interface{}                   `json:"parameter,omitempty"`
}

//-- TransferContract

type TransferContract struct {
	AssetID   string `json:"assetId,omitempty"`
	ToAddress string `json:"toAddress,omitempty"`
	Amount    int64  `json:"amount,omitempty"`
}

// -- CreateAssetContract
type CreateAssetContract struct {
	Type              string          `json:"type"`
	Name              string          `json:"name"`
	Ticker            string          `json:"ticker"`
	Logo              string          `json:"logo"`
	OwnerAddress      string          `json:"ownerAddress"`
	URIs              []*URI          `json:"uris"`
	Precision         uint32          `json:"precision"`
	InitialSupply     int64           `json:"initialSupply"`
	CirculatingSupply int64           `json:"circulatingSupply"`
	MaxSupply         int64           `json:"maxSupply"`
	MintedValue       int64           `json:"mintedValue"`
	BurnedValue       int64           `json:"burnedValue"`
	IssueDate         int64           `json:"issueDate"`
	Royalties         *RoyaltiesInfo  `json:"royalties"`
	Properties        *PropertiesInfo `json:"properties"`
	Attributes        *AttributesInfo `json:"attributes"`
	Staking           *Staking        `json:"staking"`
	Roles             []*RolesInfo    `json:"roles"`
}

type RolesInfo struct {
	Address             string `json:"address"`
	HasRoleMint         bool   `json:"hasRoleMint"`
	HasRoleSetITOPrices bool   `json:"hasRoleSetITOPrices"`
}

type RoyaltiesInfo struct {
	Address            string             `json:"address,omitempty"`
	TransferPercentage []*RoyaltyDataInfo `json:"transferPercentage,omitempty"`
	TransferFixed      int64              `json:"transferFixed,omitempty"`
	MarketPercentage   uint32             `json:"marketPercentage,omitempty"`
	MarketFixed        int64              `json:"marketFixed,omitempty"`
}

type RoyaltyDataInfo struct {
	Amount     int64  `json:"amount,omitempty"`
	Percentage uint32 `json:"percentage,omitempty"`
}

type PropertiesInfo struct {
	CanFreeze      bool `json:"canFreeze"`
	CanWipe        bool `json:"canWipe"`
	CanPause       bool `json:"canPause"`
	CanMint        bool `json:"canMint"`
	CanBurn        bool `json:"canBurn"`
	CanChangeOwner bool `json:"canChangeOwner"`
	CanAddRoles    bool `json:"canAddRoles"`
}

type AttributesInfo struct {
	IsPaused         bool `json:"isPaused"`
	IsNFTMintStopped bool `json:"isNFTMintStopped"`
}

type Staking struct {
	Type                string `json:"type"`
	APR                 uint32 `protobuf:"varint,2,opt,name=APR,proto3" json:"apr"`
	MinEpochsToClaim    uint32 `protobuf:"varint,3,opt,name=MinEpochsToClaim,proto3" json:"minEpochsToClaim"`
	MinEpochsToUnstake  uint32 `protobuf:"varint,4,opt,name=MinEpochsToUnstake,proto3" json:"minEpochsToUnstake"`
	MinEpochsToWithdraw uint32 `protobuf:"varint,5,opt,name=MinEpochsToWithdraw,proto3" json:"minEpochsToWithdraw"`
}

type StakingData struct {
	InterestType        string     `json:"interestType"`
	APR                 []*APRData `json:"apr"`
	FPR                 []*FPRData `json:"fpr"`
	TotalStaked         int64      `json:"totalStaked"`
	CurrentFPRAmount    int64      `json:"currentFPRAmount"`
	MinEpochsToClaim    uint32     `json:"minEpochsToClaim"`
	MinEpochsToUnstake  uint32     `json:"minEpochsToUnstake"`
	MinEpochsToWithdraw uint32     `json:"minEpochsToWithdraw"`
}

type APRData struct {
	Timestamp int64  `json:"timestamp"`
	Epoch     uint32 `json:"epoch"`
	Value     uint32 `json:"value"`
}

type FPRData struct {
	TotalAmount  int64  `json:"totalAmount"`
	TotalStaked  int64  `json:"totalStaked"`
	Epoch        uint32 `json:"epoch"`
	TotalClaimed int64  `json:"TotalClaimed"`
}

// -- CreateValidatorContract
type CreateValidatorContract struct {
	OwnerAddress string          `json:"ownerAddress"`
	Config       ValidatorConfig `json:"config"`
}

// -- ValidatorConfigContract
type ValidatorConfigContract struct {
	Config ValidatorConfig `json:"config,omitempty"`
}

// -- ValidatorConfig
type ValidatorConfig struct {
	BLSPublicKey        string `json:"blsPublicKey,omitempty"`
	RewardAddress       string `json:"rewardAddress,omitempty"`
	CanDelegate         bool   `json:"canDelegate"`
	Commission          uint32 `json:"commission"`
	MaxDelegationAmount int64  `json:"maxDelegationAmount"`
	Logo                string `json:"logo"`
	URIs                []*URI `json:"uris"`
	Name                string `json:"name"`
}

// -- FreezeContract
type FreezeContract struct {
	Amount  int64  `json:"amount,omitempty"`
	AssetID string `json:"assetId,omitempty"`
}

// -- UnfreezeContract
type UnfreezeContract struct {
	BucketID string `json:"bucketID,omitempty"`
	AssetID  string `json:"assetId,omitempty"`
}

// -- DelegateContract
type DelegateContract struct {
	ToAddress string `json:"toAddress,omitempty"`
	BucketID  string `json:"bucketID,omitempty"`
}

// -- UndelegateContract
type UndelegateContract struct {
	BucketID string `json:"bucketID,omitempty"`
}

// --  WithdrawContract
type WithdrawContract struct {
	AssetID string `json:"assetId,omitempty"`
}

// -- UnjailContract
type UnjailContract struct {
}

// -- RedeemContract
type RedeemContract struct {
}

// -- ClaimContract
type ClaimContract struct {
	ClaimType string `json:"claimType,omitempty"`
	ID        string `json:"id,omitempty"`
}

// -- SetAccountNameContract
type SetAccountNameContract struct {
	Name string `json:"name"`
}

// -- AssetTriggerContract
type AssetTriggerContract struct {
	TriggerType string     `json:"triggerType"`
	AssetID     string     `json:"assetId"`
	ToAddress   string     `json:"toAddress,omitempty"`
	Amount      int64      `json:"amount,omitempty"`
	MIME        string     `json:"mime,omitempty"`
	Logo        string     `json:"logo"`
	URIs        []*URI     `json:"uris"`
	Role        *RolesInfo `json:"role,omitempty"`
	Staking     *Staking   `json:"staking,omitempty"`
}

// -- ProposalContract
type ProposalContract struct {
	Parameters     map[int32]string `json:"parameters"`
	Description    string           `json:"description"`
	EpochsDuration uint32           `json:"epochsDuration"`
}

// -- VoteContract
type VoteContract struct {
	Type       string `json:"type"`
	ProposalID uint64 `json:"proposalId"`
	Amount     int64  `json:"amount"`
}

// -- ConfigITOContract
type ConfigITOContract struct {
	AssetID         string      `json:"assetId,omitempty"`
	ReceiverAddress string      `json:"receiverAddress,omitempty"`
	Status          string      `json:"status,omitempty"`
	MaxAmount       int64       `json:"maxAmount,omitempty"`
	PackInfo        []*PackInfo `json:"packInfo,omitempty"`
}

// -- SetITOPricesContract
type SetITOPricesContract struct {
	AssetID  string      `json:"assetId,omitempty"`
	PackInfo []*PackInfo `json:"packInfo,omitempty"`
}

// PackInfo holds the pack list structure for the ITO contract
type PackInfo struct {
	Key   string      `json:"key,omitempty"`
	Packs []*PackItem `json:"packs,omitempty"`
}

// PackItem hold the pack structure for the ITO contract
type PackItem struct {
	Amount int64 `json:"amount,omitempty"`
	Price  int64 `json:"price,omitempty"`
}

// -- BuyContract
type BuyContract struct {
	BuyType    string `json:"buyType,omitempty"`
	ID         string `json:"id,omitempty"`
	CurrencyID string `json:"currencyID,omitempty"`
	Amount     int64  `json:"amount,omitempty"`
}

// -- SellContract
type SellContract struct {
	MarketType    string `json:"marketType,omitempty"`
	MarketplaceID string `json:"marketplaceID,omitempty"`
	AssetID       string `json:"assetId,omitempty"`
	CurrencyID    string `json:"currencyID,omitempty"`
	Price         int64  `json:"price,omitempty"`
	ReservePrice  int64  `json:"reservePrice,omitempty"`
	EndTime       int64  `json:"endTime,omitempty"`
}

// -- CancelMarketOrderContract
type CancelMarketOrderContract struct {
	OrderID string `json:"orderID,omitempty"`
}

// -- CreateMarketplaceContract
type CreateMarketplaceContract struct {
	Name               string `json:"name,omitempty"`
	ReferralAddress    string `json:"referralAddress,omitempty"`
	ReferralPercentage uint32 `json:"referralPercentage,omitempty"`
}

// -- ConfigMarketplaceContract
type ConfigMarketplaceContract struct {
	MarketplaceID      string `json:"marketplaceID,omitempty"`
	Name               string `json:"name,omitempty"`
	ReferralAddress    string `json:"referralAddress,omitempty"`
	ReferralPercentage uint32 `json:"referralPercentage,omitempty"`
}

type AccKey struct {
	Address string `json:"address"`
	Weight  int64  `json:"weight"`
}

type AccPermission struct {
	Type           int32    `json:"type"`
	PermissionName string   `json:"permissionName"`
	Threshold      int64    `json:"threshold"`
	Operations     string   `json:"operations"`
	Signers        []AccKey `json:"signers"`
}

// -- UpdateAccountPermissionContract
type UpdateAccountPermissionContract struct {
	Permissions []AccPermission `json:"permissions"`
}
