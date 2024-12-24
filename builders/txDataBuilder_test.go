package builders

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"testing"

	"github.com/klever-io/klever-go-sdk/core/address"
	"github.com/klever-io/klever-go/tools/check"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTxDataBuilder(t *testing.T) {
	t.Parallel()

	builder := NewTxDataBuilder()
	assert.False(t, check.IfNil(builder))
	_, err := builder.ToDataString()
	assert.Nil(t, err)
}

func TestTxDataBuilder_Function(t *testing.T) {
	t.Parallel()

	function := "sc call function"
	builder := NewTxDataBuilder()
	builder.Function(function)

	txData, err := builder.ToDataString()
	assert.Nil(t, err)
	assert.Equal(t, function, txData)

	txDataBytes, err := builder.ToDataBytes()
	assert.Nil(t, err)
	assert.Equal(t, function, string(txDataBytes))
}

func TestTxDataBuilder_AllGoodArguments(t *testing.T) {
	t.Parallel()

	t.Run("with function", func(t *testing.T) {
		t.Parallel()

		address, errBech32 := address.NewAddress("klv1mge94r8n3q44hcwu2tk9afgjcxcawmutycu0cwkap7m6jnktjlvq58355l")
		require.Nil(t, errBech32)

		builder := NewTxDataBuilder().
			Function("function").
			ArgBigInt(big.NewInt(15)).
			ArgInt64(14).
			ArgAddress(address).
			ArgHexString("eeff00").
			ArgBytes([]byte("aa")).
			ArgBigInt(big.NewInt(0)).
			ArgBytesList([][]byte{[]byte("bbb"), []byte("ccc")})

		expectedTxData := "function@" + hex.EncodeToString([]byte{15}) +
			"@" + hex.EncodeToString([]byte{14}) + "@" +
			hex.EncodeToString(address.Bytes()) + "@eeff00@" +
			hex.EncodeToString([]byte("aa")) + "@00@" +
			hex.EncodeToString([]byte("bbb")) + "@" + hex.EncodeToString([]byte("ccc"))

		txData, err := builder.ToDataString()
		require.Nil(t, err)
		require.Equal(t, expectedTxData, txData)

		txDataBytes, err := builder.ToDataBytes()
		require.Nil(t, err)

		require.Equal(t, []byte(expectedTxData), txDataBytes)
	})
	t.Run("without function name (deploy SC)", func(t *testing.T) {
		address, errBech32 := address.NewAddress("klv1mge94r8n3q44hcwu2tk9afgjcxcawmutycu0cwkap7m6jnktjlvq58355l")
		require.Nil(t, errBech32)

		builder := NewTxDataBuilder().
			ArgHexString("eeff001122").
			ArgHexString("ee").
			ArgBigInt(big.NewInt(15)).
			ArgInt64(14).
			ArgAddress(address).
			ArgBytes([]byte("aa")).
			ArgBigInt(big.NewInt(0)).
			ArgBytesList([][]byte{[]byte("bbb"), []byte("ccc")})

		expectedTxData := "eeff001122@ee@" + hex.EncodeToString([]byte{15}) +
			"@" + hex.EncodeToString([]byte{14}) + "@" +
			hex.EncodeToString(address.Bytes()) + "@" + hex.EncodeToString([]byte("aa")) + "@00@" +
			hex.EncodeToString([]byte("bbb")) + "@" + hex.EncodeToString([]byte("ccc"))

		txData, err := builder.ToDataString()
		require.Nil(t, err)
		require.Equal(t, expectedTxData, txData)

		txDataBytes, err := builder.ToDataBytes()
		require.Nil(t, err)

		require.Equal(t, []byte(expectedTxData), txDataBytes)
	})
}

func TestTxDataBuilder_InvalidArguments(t *testing.T) {
	t.Parallel()

	t.Run("invalid hex string argument", func(t *testing.T) {
		builder := NewTxDataBuilder().
			Function("function").
			ArgInt64(4)
		builder.ArgHexString("invalid hex string")

		txData, errString := builder.ToDataString()
		txDataBytes, errBytes := builder.ToDataBytes()
		assert.Equal(t, errString, errBytes)
		assert.Equal(t, "", txData)
		assert.Nil(t, txDataBytes)
		assert.NotNil(t, errString)
		assert.True(t, strings.Contains(errString.Error(), "builder.ArgHexString for string"))
	})
	t.Run("nil address argument", func(t *testing.T) {
		builder := NewTxDataBuilder().
			Function("function").
			ArgInt64(4)
		builder.ArgAddress(nil)

		txData, errString := builder.ToDataString()
		txDataBytes, errBytes := builder.ToDataBytes()
		assert.Equal(t, errString, errBytes)
		assert.Equal(t, "", txData)
		assert.Nil(t, txDataBytes)
		assert.NotNil(t, errString)
		assert.True(t, errors.Is(errString, ErrNilAddress))
	})
	t.Run("nil big int argument", func(t *testing.T) {
		builder := NewTxDataBuilder().
			Function("function").
			ArgInt64(4)
		builder.ArgBigInt(nil)

		txData, errString := builder.ToDataString()
		txDataBytes, errBytes := builder.ToDataBytes()
		assert.Equal(t, errString, errBytes)
		assert.Equal(t, "", txData)
		assert.Nil(t, txDataBytes)
		assert.NotNil(t, errString)
		assert.True(t, errors.Is(errString, ErrNilValue))
	})
	t.Run("empty bytes argument", func(t *testing.T) {
		builder := NewTxDataBuilder().
			Function("function").
			ArgInt64(4)
		builder.ArgBytes(make([]byte, 0))

		txData, errString := builder.ToDataString()
		txDataBytes, errBytes := builder.ToDataBytes()
		assert.Equal(t, errString, errBytes)
		assert.Equal(t, "", txData)
		assert.Nil(t, txDataBytes)
		assert.NotNil(t, errString)
		assert.True(t, errors.Is(errString, ErrInvalidValue))
	})
}
