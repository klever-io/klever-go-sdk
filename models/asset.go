package models

import "encoding/json"

type Bucket struct {
	ID            string `json:"id"`
	StakedAt      int64  `json:"stakedAt"`
	StakedEpoch   int32  `json:"stakedEpoch"`
	UnstakedEpoch int32  `json:"unstakedEpoch"`
	Balance       int64  `json:"balance"`
	Delegation    string `json:"delegation"`
	ValidatorName string `json:"validatorName"`
}

// AccountKDA is a structure that holds information about a kda account
type AccountKDA struct {
	AccountAddress  string                `json:"address"`
	AssetID         string                `json:"assetId"`
	Collection      string                `json:"collection,omitempty"`
	NFTNonce        uint64                `json:"nftNonce,omitempty"`
	AssetName       string                `json:"assetName"`
	AssetType       KDAData_EnumAssetType `json:"assetType"`
	Balance         int64                 `json:"balance"`
	Precision       uint32                `json:"precision"`
	FrozenBalance   int64                 `json:"frozenBalance"`
	UnfrozenBalance int64                 `json:"unfrozenBalance"`
	LastClaim       UserKDALastClaim      `json:"lastClaim"`
	Buckets         []UserKDABucket       `json:"buckets"`
	Metadata        string                `json:"metadata,omitempty"`
	MIME            string                `json:"mime,omitempty"`
	MarketplaceID   string                `json:"marketplaceId,omitempty"`
	OrderID         string                `json:"orderId,omitempty"`
}

// UserKDABucket is a structure that holds information about a kda user buckets
type UserKDABucket struct {
	Id            string `json:"id"`
	StakedAt      int64  `json:"stakeAt"`
	StakedEpoch   uint32 `json:"stakedEpoch"`
	UnstakedEpoch uint32 `json:"unstakedEpoch"`
	Value         int64  `json:"balance"`
	Delegation    string `json:"delegation"`
	ValidatorName string `json:"validatorName"`
}

// UserKDALastClaim is a structure that holds information about a kda user last claim
type UserKDALastClaim struct {
	Timestamp int64  `json:"timestamp"`
	Epoch     uint32 `json:"epoch"`
}

func (a AccountKDA) String() string {
	result, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		result = make([]byte, 0)
	}

	return string(result)
}
