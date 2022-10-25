package address

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/klever-io/klever-go-sdk/core"
)

type address struct {
	bytes  []byte
	bech32 string
}

const addrLen = 32

func ZeroAddress() Address {
	addr, _ := NewAddressFromHex("0000000000000000000000000000000000000000000000000000000000000000")
	return addr
}

func NewAddressFromBytes(pkBytes []byte) (Address, error) {
	// if len(pkBytes) != addrLen {
	// 	return nil, fmt.Errorf("decoding address, expected length %d, received %d",
	// 		addrLen, len(pkBytes))
	// }

	//since the errors generated here are usually because of a bad config, they will be treated here
	conv, err := bech32.ConvertBits(pkBytes, 8, 5, true)
	if err != nil {
		return nil, err
	}

	converted, err := bech32.Encode(core.HRP, conv)
	if err != nil {
		return nil, err
	}

	return newAddress(pkBytes, converted), nil
}

func NewAddressFromHex(hexAddr string) (Address, error) {
	decodedBytes, err := hex.DecodeString(hexAddr)
	if err != nil {
		return nil, err
	}

	if len(decodedBytes) != addrLen {
		return nil, fmt.Errorf("decoding address, expected length %d, received %d",
			addrLen, len(decodedBytes))
	}

	return NewAddressFromBytes(decodedBytes)
}

func NewAddress(bech32Addr string) (Address, error) {
	decodedPrefix, buff, err := bech32.Decode(bech32Addr)
	if err != nil {
		return nil, err
	}
	if decodedPrefix != core.HRP {
		return nil, fmt.Errorf("invalid address prefix")
	}

	// warning: mind the order of the parameters, those should be inverted
	decodedBytes, err := bech32.ConvertBits(buff, 5, 8, false)
	if err != nil {
		return nil, err
	}

	if len(decodedBytes) != addrLen {
		return nil, fmt.Errorf("decoding address, expected length %d, received %d",
			addrLen, len(decodedBytes))
	}

	return newAddress(decodedBytes, bech32Addr), nil
}

func newAddress(bytes []byte, bech32 string) Address {
	addr := &address{
		bech32: bech32,
		bytes:  make([]byte, len(bytes)),
	}

	copy(addr.bytes, bytes)

	return addr
}

func (a address) Hex() string {
	return hex.EncodeToString(a.bytes)
}

func (a address) Bech32() string {
	return a.bech32
}

func (a address) Bytes() []byte {
	return append([]byte{}, a.bytes...)
}

func (a *address) IsInterfaceNil() bool {
	return a == nil
}
