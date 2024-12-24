package builders

import (
	"math/big"

	"github.com/klever-io/klever-go-sdk/core/address"
	"github.com/klever-io/klever-go-sdk/provider"
)

type vmQueryBuilder struct {
	*baseBuilder
	address    string
	callerAddr string
	function   string
}

// NewVMQueryBuilder creates a new vm query data builder
func NewVMQueryBuilder() *vmQueryBuilder {
	return &vmQueryBuilder{
		baseBuilder: &baseBuilder{},
	}
}

// Function sets the function to be called
func (builder *vmQueryBuilder) Function(function string) VMQueryBuilder {
	builder.function = function

	return builder
}

// ArgHexString adds the provided hex string to the arguments list
func (builder *vmQueryBuilder) ArgHexString(hexed string) VMQueryBuilder {
	builder.addArgHexString(hexed)

	return builder
}

// ArgAddress adds the provided address to the arguments list
func (builder *vmQueryBuilder) ArgAddress(address address.Address) VMQueryBuilder {
	builder.addArgAddress(address)

	return builder
}

// ArgBigInt adds the provided value to the arguments list
func (builder *vmQueryBuilder) ArgBigInt(value *big.Int) VMQueryBuilder {
	builder.addArgBigInt(value)

	return builder
}

// ArgInt64 adds the provided value to the arguments list
func (builder *vmQueryBuilder) ArgInt64(value int64) VMQueryBuilder {
	builder.addArgInt64(value)

	return builder
}

// ArgBytes adds the provided bytes to the arguments list. The parameter should contain at least one byte
func (builder *vmQueryBuilder) ArgBytes(bytes []byte) VMQueryBuilder {
	builder.addArgBytes(bytes)

	return builder
}

// CallerAddress sets the caller address
func (builder *vmQueryBuilder) CallerAddress(address address.Address) VMQueryBuilder {
	err := builder.checkAddress(address)
	if err != nil {
		builder.err = err
		return builder
	}

	builder.callerAddr = address.Bech32()

	return builder
}

// Address sets the destination address
func (builder *vmQueryBuilder) Address(address address.Address) VMQueryBuilder {
	err := builder.checkAddress(address)
	if err != nil {
		builder.err = err
		return builder
	}

	builder.address = address.Bech32()

	return builder
}

// ToVmValueRequest returns the VmValueRequest structure to be used in a VM call
func (builder *vmQueryBuilder) ToVmValueRequest() (*provider.VmValueRequest, error) {
	if builder.err != nil {
		return nil, builder.err
	}

	return &provider.VmValueRequest{
		Address:    builder.address,
		FuncName:   builder.function,
		CallerAddr: builder.callerAddr,
		Args:       builder.args,
	}, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (builder *vmQueryBuilder) IsInterfaceNil() bool {
	return builder == nil
}
