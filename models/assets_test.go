package models_test

import (
	"testing"

	"github.com/klever-io/klever-go-sdk/models"
	"github.com/stretchr/testify/assert"
)

func TestAssets_ToString(t *testing.T) {
	kda := &models.AccountKDA{
		AssetID: "KLV",
	}
	assert.Equal(t,
		"{\n\t\"address\": \"\",\n\t\"assetId\": \"KLV\",\n\t\"assetName\": \"\",\n\t\"assetType\": 0,\n\t\"balance\": 0,\n\t\"precision\": 0,\n\t\"frozenBalance\": 0,\n\t\"unfrozenBalance\": 0,\n\t\"lastClaim\": {\n\t\t\"timestamp\": 0,\n\t\t\"epoch\": 0\n\t},\n\t\"buckets\": null\n}",
		kda.String())
}
