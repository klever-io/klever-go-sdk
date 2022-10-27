package provider

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/klever-io/klever-go-sdk/models"
	"github.com/klever-io/klever-go-sdk/models/proto"
)

func (kc *kleverChain) GetAsset(assetID string) (*proto.KDAData, error) {
	result := struct {
		Data struct {
			Asset *proto.KDAData `json:"asset"`
		} `json:"data"`
	}{}

	err := kc.httpClient.Get(fmt.Sprintf("%s/asset/%s", kc.networkConfig.GetNodeUri(), assetID), &result)

	return result.Data.Asset, err
}

type AssetTriggerType uint32

const (
	Mint AssetTriggerType = iota
	Burn
	Wipe
	Pause
	Resume
	ChangeOwner
	AddRole
	RemoveRole
	UpdateMetadata
	StopNFTMint
	UpdateLogo
	UpdateURIs
	ChangeRoyaltiesReceiver
	UpdateStaking
	UpdateRoyalties
	UpdateKDAFeePool
)

func (kc *kleverChain) AssetTrigger(
	base *models.BaseTX,
	kdaID string,
	triggerType AssetTriggerType,
	op *models.AssetTriggerOptions,
) (*proto.Transaction, error) {
	// check if is NFT
	kda := strings.Split(kdaID, "/")
	if len(kda) > 2 {
		return nil, fmt.Errorf("invalid KDA ID")
	}

	asset, err := kc.GetAsset(kda[0])
	if err != nil {
		return nil, err
	}

	parsedAmount := op.Amount

	if asset.AssetType == proto.KDAData_Fungible {
		parsedAmount = parsedAmount * math.Pow10(int(asset.Precision))
	}

	if len(op.AddRolesMint) == 1 &&
		len(op.AddRolesSetITOPrices) == 1 &&
		op.AddRolesMint[0] != op.AddRolesSetITOPrices[0] {
		return nil, fmt.Errorf("can only set one address roler per trigger")
	}

	role := &models.RolesInfo{}
	switch len(op.AddRolesMint) {
	case 0:
	case 1:
		role.Address = op.AddRolesMint[0]
		role.HasRoleMint = true
	default:
		return nil, fmt.Errorf("can only add one roler per trigger")
	}

	switch len(op.AddRolesSetITOPrices) {
	case 0:
	case 1:
		role.Address = op.AddRolesSetITOPrices[0]
		role.HasRoleSetITOPrices = true
	default:
		return nil, fmt.Errorf("can only add one roler per trigger")
	}

	var stakingInfo *models.StakingInfo
	if len(op.Staking) > 0 {
		stakingInfo = &models.StakingInfo{}
		apr, err := strconv.ParseFloat(op.Staking["apr"], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid apr %s: %w", op.Staking["apr"], err)
		}

		claim, err := strconv.ParseUint(op.Staking["claim"], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid claim min epochs %s: %w", op.Staking["claim"], err)
		}

		unstake, err := strconv.ParseUint(op.Staking["unstake"], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid unstake min epochs %s: %w", op.Staking["unstake"], err)
		}

		withdraw, err := strconv.ParseUint(op.Staking["withdraw"], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid withdraw min epochs %s: %w", op.Staking["withdraw"], err)
		}

		stakingInfo.APR = uint32(apr * math.Pow10(2))
		stakingInfo.MinEpochsToClaim = uint32(claim)
		stakingInfo.MinEpochsToUnstake = uint32(unstake)
		stakingInfo.MinEpochsToWithdraw = uint32(withdraw)
	}

	contracts := make([]interface{}, 0)
	contracts = append(contracts, models.AssetTriggerTXRequest{
		TriggerType: uint32(triggerType),
		AssetID:     kdaID,
		Amount:      int64(parsedAmount),
		Receiver:    op.Receiver,
		MIME:        op.Mime,
		Logo:        op.Logo,
		URIs:        op.URIs,
		Role:        role,
		Staking:     stakingInfo,
	})

	data, err := kc.buildRequest(proto.TXContract_AssetTriggerContractType, base, contracts)
	if err != nil {
		return nil, err
	}
	return kc.PrepareTransaction(data)
}

func (kc *kleverChain) CreateKDA(
	base *models.BaseTX,
	kdaType proto.KDAData_EnumAssetType,
	op *models.KDAOptions,
) (*proto.Transaction, error) {
	if !IsNameValid(op.Name) {
		return nil, fmt.Errorf("invalid KDA name")
	}

	if !IsTickerValid(op.Ticker) {
		return nil, fmt.Errorf("invalid KDA ticker")
	}

	if !IsPrecisionValid(op.Precision) {
		return nil, fmt.Errorf("invalid KDA precision")
	}

	if len(op.Roles) == 0 {
		op.Roles = []*models.RolesInfo{
			{
				Address:             base.FromAddress,
				HasRoleMint:         true,
				HasRoleSetITOPrices: true,
			},
		}
	}
	if len(op.Royalties.Address) == 0 {
		op.Royalties.Address = base.FromAddress
	}

	contracts := make([]interface{}, 0)
	contracts = append(contracts, models.CreateAssetTXRequest{
		Type:          uint32(kdaType),
		OwnerAddress:  base.FromAddress,
		Name:          op.Name,
		Ticker:        op.Ticker,
		Precision:     uint32(op.Precision),
		InitialSupply: int64(op.InitialSupply * math.Pow10(int(op.Precision))),
		MaxSupply:     int64(op.MaxSupply * math.Pow10(int(op.Precision))),
		Logo:          op.Logo,
		URIs:          op.URIs,
		Royalties:     &op.Royalties,
		Attributes:    &op.Attributes,
		Properties:    &op.Properties,
		Staking:       &op.Staking,
		Roles:         op.Roles,
	})

	data, err := kc.buildRequest(proto.TXContract_CreateAssetContractType, base, contracts)
	if err != nil {
		return nil, err
	}
	return kc.PrepareTransaction(data)
}

func IsNameValid(name string) bool {
	if len(name) < 1 ||
		len(name) > 32 ||
		name == "KLV" ||
		name == "KFI" {
		return false
	}

	for _, ch := range []byte(name) {
		isSmallCharacter := ch >= 'a' && ch <= 'z'
		isBigCharacter := ch >= 'A' && ch <= 'Z'
		isNumber := ch >= '0' && ch <= '9'
		isSpace := ch == ' '
		isReadable := isSmallCharacter || isBigCharacter || isNumber || isSpace
		if !isReadable {
			return false
		}
	}
	return true
}

func IsTickerValid(tickerName string) bool {
	if len(tickerName) < 3 ||
		len(tickerName) > 10 ||
		tickerName == "KLV" ||
		tickerName == "KFI" {
		return false
	}

	for _, ch := range []byte(tickerName) {
		isBigCharacter := ch >= 'A' && ch <= 'Z'
		isNumber := ch >= '0' && ch <= '9'
		isReadable := isBigCharacter || isNumber
		if !isReadable {
			return false
		}
	}

	return true
}

func IsPrecisionValid(precision int) bool {
	if precision < 0 || precision > 8 {
		return false
	}

	return true
}
