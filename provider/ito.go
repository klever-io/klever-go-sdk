package provider

import (
	"fmt"
	"math"
	"strings"

	"github.com/klever-io/klever-go-sdk/models"
	"github.com/klever-io/klever-go-sdk/models/proto"
)

func (kc *kleverChain) ConfigITO(base *models.BaseTX, kdaID, receiverAddress string, status int32, maxAmount float64, packs []models.ParsedPack) (*proto.Transaction, error) {
	kda, err := kc.GetAsset(kdaID)
	if err != nil {
		return nil, err
	}
	parsedmaxAmount := maxAmount * math.Pow10(int(kda.Precision))

	packInfo, err := kc.createPackInfo(kda.Precision, packs)
	if err != nil {
		return nil, err
	}

	configITO := models.ConfigITOTXRequest{
		KDA:             kdaID,
		ReceiverAddress: receiverAddress,
		Status:          status,
		MaxAmount:       int64(parsedmaxAmount),
		PackInfo:        packInfo,
	}

	data, err := kc.buildRequest(proto.TXContract_ConfigITOContractType, base, []interface{}{configITO})
	if err != nil {
		return nil, err
	}

	return kc.PrepareTransaction(data)
}

func (kc *kleverChain) createPackInfo(precision uint32, packs []models.ParsedPack) (map[string]models.PackInfoRequest, error) {
	packInfo := make(map[string]models.PackInfoRequest)

	for _, p := range packs {
		packPrecision, err := kc.getPrecision(p.Kda)
		if err != nil {
			return nil, err
		}

		packItems := make([]models.PackItemRequest, 0)
		for _, pItem := range p.Packs {
			parsedItemAmount := pItem.Amount * math.Pow10(int(precision))
			parsedItemPrice := pItem.Price * math.Pow10(int(packPrecision))

			packItems = append(packItems, models.PackItemRequest{Amount: int64(parsedItemAmount), Price: int64(parsedItemPrice)})
		}

		packInfo[p.Kda] = models.PackInfoRequest{Packs: packItems}
	}

	return packInfo, nil
}

func (kc *kleverChain) SetITOPrices(base *models.BaseTX, kdaID string, packs []models.ParsedPack) (*proto.Transaction, error) {
	kda, err := kc.GetAsset(kdaID)
	if err != nil {
		return nil, err
	}

	packInfo, err := kc.createPackInfo(kda.Precision, packs)
	if err != nil {
		return nil, err
	}

	setITOPrices := models.SetITOPricesTXRequest{
		KDA:      kdaID,
		PackInfo: packInfo,
	}

	data, err := kc.buildRequest(proto.TXContract_SetITOPricesContractType, base, []interface{}{setITOPrices})
	if err != nil {
		return nil, err
	}

	return kc.PrepareTransaction(data)
}

type ITOTriggerType uint32

const (
	SetITOPrices ITOTriggerType = iota
	UpdateStatus
	UpdateReceiverAddress
	UpdateMaxAmount
	UpdateDefaultLimitPerAddress
	UpdateTimes
	UpdateWhitelistStatus
	AddToWhitelist
	RemoveFromWhitelist
	UpdateWhitelistTimes
)

func (kc *kleverChain) ITOTrigger(
	base *models.BaseTX,
	kdaID string,
	triggerType ITOTriggerType,
	op *models.ITOTriggerOptions,
) (*proto.Transaction, error) {
	// check if kda isn't a NFT
	kda := strings.Split(kdaID, "/")
	if len(kda) > 1 {
		return nil, fmt.Errorf("invalid KDA ID")
	}

	asset, err := kc.GetAsset(kda[0])
	if err != nil {
		return nil, err
	}

	parsedMaxAmount := op.MaxAmount

	if asset.AssetType == proto.KDAData_Fungible {
		parsedMaxAmount = parsedMaxAmount * math.Pow10(int(asset.Precision))
	}

	contracts := make([]interface{}, 0)
	contracts = append(contracts, models.ITOTriggerTXRequest{
		TriggerType:            uint32(triggerType),
		KDA:                    kdaID,
		ReceiverAddress:        op.ReceiverAddress,
		MaxAmount:              int64(parsedMaxAmount),
		DefaultLimitPerAddress: int64(op.DefaultLimitPerAddress),
		Status:                 int32(op.Status),
		StartTime:              op.StartTime,
		EndTime:                op.EndTime,
		PackInfo:               op.PackInfo,
		WhitelistStatus:        int32(op.WhitelistStatus),
		WhitelistInfo:          op.WhitelistInfo,
		WhitelistStartTime:     op.WhitelistStartTime,
		WhitelistEndTime:       op.WhitelistEndTime,
	})

	data, err := kc.buildRequest(proto.TXContract_ITOTriggerContractType, base, contracts)
	if err != nil {
		return nil, err
	}
	return kc.PrepareTransaction(data)
}
