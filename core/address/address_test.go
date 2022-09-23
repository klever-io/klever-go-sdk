package address_test

import (
	"encoding/hex"
	"testing"

	"github.com/klever-io/klever-go-sdk/core/address"
	"github.com/stretchr/testify/assert"
)

func TestAddress_ZeroAddress(t *testing.T) {
	addr := address.ZeroAddress()
	assert.Equal(t, addr.Bech32(), "klv1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqpgm89z")
}

func TestAddress_Address_ShouldFail(t *testing.T) {
	// invalid size
	_, err := address.NewAddressFromHex("000000000000000000000000000000000000000000000000000000000000000")
	assert.NotNil(t, err)

	// invalid hex
	_, err = address.NewAddressFromHex("X000000000000000000000000000000000000000000000000000000000000000")
	assert.NotNil(t, err)

	// invalid bech32
	_, err = address.NewAddress("klv1qy352eufzqg3yyc5z5v3wxqeyqsjygeyy5nzw2pfxqcnyve5x5mqfrkqfh")
	assert.NotNil(t, err)

	// invalid bech32 prefix
	_, err = address.NewAddress("kfi1qy352eufzqg3yyc5z5v3wxqeyqsjygeyy5nzw2pfxqcnyve5x5mq7ze5xk")
	assert.Nil(t, err)
}

func TestAddress_Address_ShouldWork(t *testing.T) {
	hexString := "0123456789101112131415191718192021222324252627282930313233343536"
	bytes, _ := hex.DecodeString(hexString)
	bech32Addr := "klv1qy352eufzqg3yyc5z5v3wxqeyqsjygeyy5nzw2pfxqcnyve5x5mqfrkqfg"

	// invalid size
	addr, err := address.NewAddressFromHex(hexString)
	assert.Nil(t, err)
	assert.Equal(t, addr.Bech32(), bech32Addr)
	assert.Equal(t, addr.Hex(), hexString)

	addr, err = address.NewAddressFromBytes(bytes)
	assert.Nil(t, err)
	assert.Equal(t, addr.Bech32(), bech32Addr)
	assert.Equal(t, addr.Hex(), hexString)

	addr, err = address.NewAddress(bech32Addr)
	assert.Nil(t, err)
	assert.Equal(t, addr.Bech32(), bech32Addr)
	assert.Equal(t, addr.Hex(), hexString)
}
