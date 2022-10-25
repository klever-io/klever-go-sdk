package provider

import (
	"math"

	"github.com/klever-io/klever-go-sdk/models"
	"github.com/klever-io/klever-go-sdk/models/proto"
)

func (kc *kleverChain) CreateMarketplace(base *models.BaseTX, name, referralAddr string, referralPercent float64) (*proto.Transaction, error) {
	parsedReferralPercent := referralPercent * math.Pow10(2)

	createMarketplace := models.CreateMarketplaceTXRequest{
		Name:               name,
		ReferralAddress:    referralAddr,
		ReferralPercentage: uint32(parsedReferralPercent),
	}

	data, err := kc.buildRequest(proto.TXContract_CreateMarketplaceContractType, base, []interface{}{createMarketplace})
	if err != nil {
		return nil, err
	}

	return kc.PrepareTransaction(data)
}

func (kc *kleverChain) ConfigMarketplace(base *models.BaseTX, id, name, referralAddr string, referralPercent float64) (*proto.Transaction, error) {
	parsedReferralPercent := referralPercent * math.Pow10(2)

	configMarketplace := models.ConfigMarketplaceTXRequest{
		MarketplaceID:      id,
		Name:               name,
		ReferralAddress:    referralAddr,
		ReferralPercentage: uint32(parsedReferralPercent),
	}

	data, err := kc.buildRequest(proto.TXContract_ConfigMarketplaceContractType, base, []interface{}{configMarketplace})
	if err != nil {
		return nil, err
	}

	return kc.PrepareTransaction(data)
}

func (kc *kleverChain) BuyOrder(base *models.BaseTX, id, currency string, amount float64, buyType int32) (*proto.Transaction, error) {
	parsedAmount := amount

	precision, err := kc.getPrecision(id)
	if err != nil {
		return nil, err
	}

	parsedAmount = amount * math.Pow10(int(precision))

	buyOrder := models.BuyTXRequest{
		BuyType:    buyType,
		ID:         id,
		CurrencyID: currency,
		Amount:     int64(parsedAmount),
	}

	data, err := kc.buildRequest(proto.TXContract_BuyContractType, base, []interface{}{buyOrder})
	if err != nil {
		return nil, err
	}

	return kc.PrepareTransaction(data)
}

func (kc *kleverChain) SellOrder(base *models.BaseTX, kdaID, currency, mktID string, price, reservePrice float64, endTime int64, mktType int32, message string) (*proto.Transaction, error) {
	precision, err := kc.getPrecision(currency)
	if err != nil {
		return nil, err
	}

	parsedPrice := price * math.Pow10(int(precision))
	parsedReservePrice := reservePrice * math.Pow10(int(precision))

	sellOrder := models.SellTXRequest{
		MarketType:    mktType,
		MarketplaceID: mktID,
		AssetID:       kdaID,
		CurrencyID:    currency,
		Price:         int64(parsedPrice),
		ReservePrice:  int64(parsedReservePrice),
		EndTime:       endTime,
	}

	data, err := kc.buildRequest(proto.TXContract_SellContractType, base, []interface{}{sellOrder})
	if err != nil {
		return nil, err
	}

	return kc.PrepareTransaction(data)
}

func (kc *kleverChain) CancelMarketOrder(base *models.BaseTX, orderID string) (*proto.Transaction, error) {
	cancelMarketOrder := models.CancelMarketOrderTXRequest{
		OrderID: orderID,
	}

	data, err := kc.buildRequest(proto.TXContract_CancelMarketOrderContractType, base, []interface{}{cancelMarketOrder})
	if err != nil {
		return nil, err
	}

	return kc.PrepareTransaction(data)
}
