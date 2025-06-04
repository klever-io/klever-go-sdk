package builders

import (
	"math/big"
	"strings"

	"github.com/klever-io/klever-go-sdk/core/address"
)

const dataSeparator = "@"

// txDataBuilder can be used to easy construct a transaction's data field for a smart contract call
// can also be used to construct a VmValueRequest instance ready to be used on a VM query
type txDataBuilder struct {
	*baseBuilder
	function string
}

// NewTxDataBuilder creates a new transaction data builder
func NewTxDataBuilder() *txDataBuilder {
	return &txDataBuilder{
		baseBuilder: &baseBuilder{},
	}
}

// Function sets the function to be called
func (builder *txDataBuilder) Function(function string) TxDataBuilder {
	builder.function = function

	return builder
}

// ArgHexString adds the provided hex string to the arguments list
func (builder *txDataBuilder) ArgHexString(hexed string) TxDataBuilder {
	builder.addArgHexString(hexed)

	return builder
}

// ArgAddress adds the provided address to the arguments list
func (builder *txDataBuilder) ArgAddress(address address.Address) TxDataBuilder {
	builder.addArgAddress(address)

	return builder
}

// ArgBigInt adds the provided value to the arguments list
func (builder *txDataBuilder) ArgBigInt(value *big.Int) TxDataBuilder {
	builder.addArgBigInt(value)

	return builder
}

// ArgInt64 adds the provided value to the arguments list
func (builder *txDataBuilder) ArgInt64(value int64) TxDataBuilder {
	builder.addArgInt64(value)

	return builder
}

// ArgBytes adds the provided bytes to the arguments list. The parameter should contain at least one byte
func (builder *txDataBuilder) ArgBytes(bytes []byte) TxDataBuilder {
	builder.addArgBytes(bytes)

	return builder
}

// ArgBytesList adds the provided list of bytes. Each argument from the list should contain at least one byte
func (builder *txDataBuilder) ArgBytesList(list [][]byte) TxDataBuilder {
	for _, arg := range list {
		builder.addArgBytes(arg)
	}

	return builder
}

// ToDataString returns the formatted data string ready to be used in a transaction call
func (builder *txDataBuilder) ToDataString() (string, error) {
	if builder.err != nil {
		return "", builder.err
	}

	var parts []string
	if len(builder.function) > 0 {
		parts = append(parts, builder.function)
	}

	parts = append(parts, builder.args...)

	return strings.Join(parts, dataSeparator), nil
}

// ToDataBytes returns the formatted data string ready to be used in a transaction call as bytes
func (builder *txDataBuilder) ToDataBytes() ([]byte, error) {
	dataField, err := builder.ToDataString()
	if err != nil {
		return nil, err
	}

	return []byte(dataField), err
}

// IsInterfaceNil returns true if there is no value under the interface
func (builder *txDataBuilder) IsInterfaceNil() bool {
	return builder == nil
}
