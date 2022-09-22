package provider

import (
	"github.com/klever-io/klever-go-sdk/models"
)

func (kc *kleverChain) SetPermission(base *models.BaseTX, permissions []models.PermissionTXRequest) (*models.Transaction, error) {

	contracts := []interface{}{models.UpdateAccountPermissionTXRequest{
		Permissions: permissions,
	}}

	data, err := kc.buildRequest(models.TXContract_UpdateAccountPermissionContractType, base, contracts)
	if err != nil {
		return nil, err
	}
	return kc.PrepareTransaction(data)
}
