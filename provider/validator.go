package provider

import (
	"math"

	"github.com/klever-io/klever-go-sdk/models"
	"github.com/klever-io/klever-go-sdk/models/proto"
)

func (kc *kleverChain) CreateValidator(base *models.BaseTX, blsKey, ownerAddr, rewardAddr, logo string, commission, maxDelegation float64, canDelegate bool, uris map[string]string, name string) (*proto.Transaction, error) {
	parsedCommission := commission * math.Pow10(2)
	parsedMaxDelegation := maxDelegation * math.Pow10(6)

	contracts := []interface{}{models.CreateValidatorTXRequest{
		BLSPublicKey:        blsKey,
		OwnerAddress:        ownerAddr,
		RewardAddress:       rewardAddr,
		Commission:          uint32(parsedCommission),
		CanDelegate:         canDelegate,
		MaxDelegationAmount: int64(parsedMaxDelegation),
		Logo:                logo,
		Name:                name,
		URIs:                uris,
	}}

	data, err := kc.buildRequest(proto.TXContract_CreateValidatorContractType, base, contracts)
	if err != nil {
		return nil, err
	}
	return kc.PrepareTransaction(data)
}

func (kc *kleverChain) ValidatorConfig(base *models.BaseTX, blsKey, rewardAddr, logo string, commission, maxDelegation float64, canDelegate bool, uris map[string]string, name string) (*proto.Transaction, error) {
	parsedCommission := commission * math.Pow10(2)
	parsedMaxDelegation := maxDelegation * math.Pow10(6)

	contracts := []interface{}{models.ValidatorConfigTXRequest{
		BLSPublicKey:        blsKey,
		RewardAddress:       rewardAddr,
		CanDelegate:         canDelegate,
		Commission:          uint32(parsedCommission),
		MaxDelegationAmount: int64(parsedMaxDelegation),
		Logo:                logo,
		Name:                name,
		URIs:                uris,
	}}

	data, err := kc.buildRequest(proto.TXContract_ValidatorConfigContractType, base, contracts)
	if err != nil {
		return nil, err
	}
	return kc.PrepareTransaction(data)
}
