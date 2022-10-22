package provider

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/klever-io/klever-go-sdk/models"
)

func (kc *kleverChain) GetAsset(assetID string) (*models.KDAData, error) {
	result := struct {
		Data struct {
			Asset *models.KDAData `json:"asset"`
		} `json:"data"`
	}{}

	err := kc.httpClient.Get(fmt.Sprintf("%s/assets/%s", kc.networkConfig.GetAPIUri(), assetID), &result)

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

type AssetTriggerOprions struct {
	Amount               float64
	AddRolesMint         []string
	AddRolesSetITOPrices []string
	Staking              map[string]string
	Receiver             string
	Mime                 string
	Logo                 string
	URIs                 map[string]string
}

func (kc *kleverChain) AssetTrigger(
	base *models.BaseTX,
	kdaID string,
	triggerType AssetTriggerType,
	op *AssetTriggerOprions,
) (*models.Transaction, error) {
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

	if asset.AssetType == models.KDAData_Fungible {
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

	stakingInfo := &models.StakingInfo{}
	if len(op.Staking) > 0 {
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

	data, err := kc.buildRequest(models.TXContract_AssetTriggerContractType, base, contracts)
	if err != nil {
		return nil, err
	}
	return kc.PrepareTransaction(data)
}
