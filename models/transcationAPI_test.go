package models_test

import (
	"testing"

	"github.com/klever-io/klever-go-sdk/models"
	"github.com/stretchr/testify/assert"
)

func TestTranscationAPI_ToString(t *testing.T) {
	tx := &models.TransactionAPI{
		Hash: "1234",
	}
	assert.Equal(t,
		"{\n\t\"hash\": \"1234\",\n\t\"sender\": \"\",\n\t\"nonce\": 0,\n\t\"kAppFee\": 0,\n\t\"bandwidthFee\": 0,\n\t\"status\": \"\",\n\t\"searchOrder\": 0,\n\t\"receipts\": null,\n\t\"contract\": null\n}",
		tx.String(),
	)
}
