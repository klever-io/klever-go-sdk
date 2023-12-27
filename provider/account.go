package provider

import (
	"context"
	"fmt"

	"github.com/klever-io/klever-go-sdk/models"
)

func (kc *kleverChain) GetAccount(ctx context.Context, address string) (*models.Account, error) {
	result := struct {
		Data struct {
			Account *models.Account `json:"account"`
		} `json:"data"`
	}{}

	err := kc.httpClient.Get(ctx, fmt.Sprintf("%s/address/%s", kc.networkConfig.GetAPIUri(), address), &result)

	return result.Data.Account, err
}

func (kc *kleverChain) GetAccountAllowance(address string, kda string) (*models.AccountAllowance, error) {
	result := struct {
		Data struct {
			Result *models.AccountAllowance `json:"result"`
		} `json:"data"`
	}{}

	err := kc.httpClient.Get(context.Background(), fmt.Sprintf("%s/address/%s/allowance?assetID=%s", kc.networkConfig.GetAPIUri(), address, kda), &result)

	return result.Data.Result, err
}
