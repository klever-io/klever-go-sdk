package provider

import (
	"fmt"

	"github.com/klever-io/klever-go-sdk/models"
)

func (kc *kleverChain) GetAsset(assetID string) (models.KDAData, error) {
	result := struct {
		Data struct {
			Asset models.KDAData `json:"asset"`
		} `json:"data"`
	}{}

	err := kc.httpClient.Get(fmt.Sprintf("%s/assets/%s", kc.networkConfig.GetAPIUri(), assetID), &result)

	return result.Data.Asset, err
}
