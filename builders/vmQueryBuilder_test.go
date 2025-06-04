package builders

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/klever-io/klever-go-sdk/core/address"
	"github.com/klever-io/klever-go/tools/check"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVMQueryBuilder(t *testing.T) {
	t.Parallel()

	builder := NewVMQueryBuilder()
	assert.False(t, check.IfNil(builder))
	_, err := builder.ToVmValueRequest()
	assert.Nil(t, err)
}

func TestTxDataBuilder_Address(t *testing.T) {
	t.Parallel()

	addr, errBech32 := address.NewAddress("klv1mge94r8n3q44hcwu2tk9afgjcxcawmutycu0cwkap7m6jnktjlvq58355l")
	require.Nil(t, errBech32)

	t.Run("nil address should contain error", func(t *testing.T) {
		builder := NewVMQueryBuilder()
		builder.Address(nil)
		valueRequest, err := builder.ToVmValueRequest()
		assert.True(t, errors.Is(err, ErrNilAddress))
		assert.Nil(t, valueRequest)
	})
	// TODO: check if this test is valid, since we need to ignore the error to proceed
	// t.Run("invalid address should contain error", func(t *testing.T) {
	// 	builder := NewVMQueryBuilder()
	// 	invalidAddr, _ := address.NewAddressFromBytes(make([]byte, 0))
	// 	builder.Address(invalidAddr)
	// 	valueRequest, err := builder.ToVmValueRequest()
	// 	assert.True(t, errors.Is(err, ErrInvalidAddress))
	// 	assert.Nil(t, valueRequest)
	// })
	t.Run("should work", func(t *testing.T) {
		builder := NewVMQueryBuilder()
		builder.Address(addr)
		valueRequest, err := builder.ToVmValueRequest()
		assert.Nil(t, err)

		addressAsBech32String := addr.Bech32()
		assert.Equal(t, addressAsBech32String, valueRequest.Address)
	})
}

func TestTxDataBuilder_CallerAddress(t *testing.T) {
	t.Parallel()

	addr, errBech32 := address.NewAddress("klv1mge94r8n3q44hcwu2tk9afgjcxcawmutycu0cwkap7m6jnktjlvq58355l")
	require.Nil(t, errBech32)

	t.Run("nil address should contain error", func(t *testing.T) {
		builder := NewVMQueryBuilder()
		builder.CallerAddress(nil)
		valueRequest, err := builder.ToVmValueRequest()
		assert.True(t, errors.Is(err, ErrNilAddress))
		assert.Nil(t, valueRequest)
	})
	// TODO: check if this test is valid, since we need to ignore the error to proceed
	// t.Run("invalid address should contain error", func(t *testing.T) {
	// 	builder := NewVMQueryBuilder()
	// 	invalidAddr, _ := address.NewAddressFromBytes(make([]byte, 0))
	// 	builder.CallerAddress(invalidAddr)
	// 	valueRequest, err := builder.ToVmValueRequest()
	// 	assert.True(t, errors.Is(err, ErrInvalidAddress))
	// 	assert.Nil(t, valueRequest)
	// })
	t.Run("should work", func(t *testing.T) {
		builder := NewVMQueryBuilder()
		builder.CallerAddress(addr)
		valueRequest, err := builder.ToVmValueRequest()
		assert.Nil(t, err)
		addressAsBech32String := addr.Bech32()
		assert.Equal(t, addressAsBech32String, valueRequest.CallerAddr)
	})
}

func TestVmQueryBuilder_AllGoodArguments(t *testing.T) {
	t.Parallel()

	address, errBech32 := address.NewAddress("klv1mge94r8n3q44hcwu2tk9afgjcxcawmutycu0cwkap7m6jnktjlvq58355l")
	require.Nil(t, errBech32)

	builder := NewVMQueryBuilder().
		Function("function").
		ArgBigInt(big.NewInt(15)).
		ArgInt64(14).
		ArgAddress(address).
		ArgHexString("eeff00").
		ArgBytes([]byte("aa")).
		ArgBigInt(big.NewInt(0))

	valueRequest, err := builder.ToVmValueRequest()
	assert.Nil(t, err)
	assert.Equal(t, "function", valueRequest.FuncName)

	expectedArgs := []string{
		hex.EncodeToString([]byte{15}),
		hex.EncodeToString([]byte{14}),
		hex.EncodeToString(address.Bytes()),
		"eeff00",
		hex.EncodeToString([]byte("aa")),
		"00",
	}

	require.Equal(t, expectedArgs, valueRequest.Args)
}
