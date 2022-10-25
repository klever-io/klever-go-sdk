package models_test

import (
	"testing"

	"github.com/klever-io/klever-go-sdk/models"
	"github.com/stretchr/testify/assert"
)

func TestAccount_ToString(t *testing.T) {
	acc := &models.Account{
		AccountInfo: &models.AccountInfo{
			Address: "1234",
		},
	}
	assert.Equal(t,
		"{\n\t\"address\": \"1234\",\n\t\"nonce\": 0,\n\t\"balance\": 0,\n\t\"frozenBalance\": 0,\n\t\"allowance\": 0,\n\t\"permissions\": null,\n\t\"timestamp\": 0,\n\t\"assets\": null\n}",
		acc.String(),
	)
}

func TestAccountAlloance_ToString(t *testing.T) {
	acc := &models.AccountAllowance{
		Allowance:      123,
		StakingRewards: 456,
	}
	assert.Equal(t,
		"{\n\t\"allowance\": 123,\n\t\"stakingRewards\": 456\n}",
		acc.String(),
	)
}
