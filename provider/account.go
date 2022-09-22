package provider

import (
	"fmt"

	"github.com/klever-io/klever-go-sdk/models"
)

func (kc *kleverChain) GetAccount(address string) (models.Account, error) {
	result := struct {
		Data struct {
			Account models.Account `json:"account"`
		} `json:"data"`
	}{}

	err := kc.httpClient.Get(fmt.Sprintf("%s/address/%s", kc.networkConfig.GetAPIUri(), address), &result)

	return result.Data.Account, err
}
