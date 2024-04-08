package provider

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"

	"github.com/klever-io/klever-go-sdk/core/address"
)

const (
	U64HexLength int = 16
	U32HexLength int = 8
	U16HexLength int = 4
	U8HexLength  int = 2

	BaseHex int = 16
)

type output struct {
	Type string `json:"type"`
}

type endpoint struct {
	Name       string   `json:"name"`
	Mutability string   `json:"mutability"`
	Outputs    []output `json:"outputs"`
}

type field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type typeInfos struct {
	Type   string  `json:"type"` // struct, tuple, variadic, List
	Fields []field `json:"fields"`
}

type abiData struct {
	Endpoints []endpoint           `json:"endpoints"`
	Types     map[string]typeInfos `json:"types"`
}

type SCOutputDecoder interface {
	Decode(abiPath, endpointName, hexValue string) (interface{}, error)
}

func (a *abiData) Decode(abi io.Reader, endpoint, hex string) (interface{}, error) {
	if err := a.loadAbi(abi); err != nil {
		return nil, err
	}

	endpointIndex, err := a.findEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	a.selectDecoder(&hex, *endpointIndex)

	// TODO: To implement
	return nil, fmt.Errorf("Please implement me T-T")
}

func (a *abiData) loadAbi(r io.Reader) error {
	jsonBytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonBytes, a); err != nil {
		return err
	}

	return nil
}

func (a *abiData) findEndpoint(endpointName string) (*int, error) {
	var endpointIndex *int

	for index, endpoint := range a.Endpoints {
		if endpoint.Name == endpointName {
			endpointIndex = &index
			break
		}
	}

	if endpointIndex == nil {
		return nil, fmt.Errorf("endpoint %s not found", endpointName)
	}

	return endpointIndex, nil
}

func (a *abiData) selectDecoder(hexValue *string, endpointIndex int) (interface{}, error) {
	vType := a.Endpoints[endpointIndex].Outputs[0].Type
	typeWrapper := strings.Split(vType, "<")

	switch typeWrapper[0] {
	case "List":
		return nil, fmt.Errorf("List")
	case "Option":
		return nil, fmt.Errorf("Option")
	case "tuple":
		return nil, fmt.Errorf("tuple")
	case "variadic":
		return nil, fmt.Errorf("variadic")
	default:
		return a.decodeSingleValue(*hexValue, vType)
	}
}

func (a *abiData) decodeSingleValue(hexValue string, vType string) (interface{}, error) {
	switch vType {
	case "i8":
		return a.decodeInt8(hexValue)
	case "i16":
		return a.decodeInt16(hexValue)
	case "i32":
		return a.decodeInt32(hexValue)
	case "i64":
		return a.decodeInt64(hexValue)
	case "u8":
		return a.decodeUint8(hexValue)
	case "u16":
		return a.decodeUint16(hexValue)
	case "u32":
		return a.decodeUint32(hexValue)
	case "u64":
		return a.decodeUint64(hexValue)
	case "BigInt":
		return a.decodeBigInt(hexValue)
	case "BigUint":
		return a.decodeBigUint(hexValue)
	case "bool":
		return hexValue == "01", nil
	case
		"ManagedBuffer",
		"TokenIdentifier",
		"bytes",
		"BoxedBytes",
		"String",
		"&str",
		"Vec<u8>",
		"&[u8]":
		return a.decodeString(hexValue)
	case "Address":
		return a.decodeAddress(hexValue)
	default:
		return nil, fmt.Errorf("invalid type %s", vType)
	}
}

func (a *abiData) decodeAddress(hexValue string) (*string, error) {
	decodedAddress, err := address.NewAddressFromHex(hexValue)
	if err != nil {
		return nil, err
	}

	addressString := decodedAddress.Bech32()
	return &addressString, nil
}

func (a *abiData) decodeString(hexValue string) (string, error) {
	bytes, err := hex.DecodeString(hexValue)

	if err != nil {
		return "", err
	}

	convertedString := string(bytes)

	return convertedString, nil
}

func (a *abiData) decodeInt(hexValue string, bitSize int) (*uint, error) {
	uintValue, err := strconv.ParseUint(hexValue, BaseHex, bitSize)
	if err != nil {
		return nil, err
	}

	return a.fixSignedIntOverflow(&uintValue, len(hexValue))
}

func (a *abiData) decodeUint(hexValue string, bitSize int) (*uint64, error) {
	uintValue, err := strconv.ParseUint(hexValue, BaseHex, bitSize)
	if err != nil {
		return nil, err
	}

	return &uintValue, nil
}

func (a *abiData) decodeUint8(hexValue string) (uint8, error) {
	targetValue, err := a.decodeUint(hexValue, 8)
	uint8Decoded := uint8(*targetValue)

	return uint8Decoded, err
}

func (a *abiData) decodeUint16(hexValue string) (uint16, error) {
	targetValue, err := a.decodeUint(hexValue, 16)
	uint16Decoded := uint16(*targetValue)

	return uint16Decoded, err
}

func (a *abiData) decodeUint32(hexValue string) (uint32, error) {
	targetValue, err := a.decodeUint(hexValue, 32)
	uint32Decoded := uint32(*targetValue)

	return uint32Decoded, err
}

func (a *abiData) decodeUint64(hexValue string) (uint64, error) {
	targetValue, err := a.decodeUint(hexValue, 64)
	uint64Decoded := uint64(*targetValue)

	return uint64Decoded, err
}

func (a *abiData) decodeBigUint(hexString string) (*big.Int, error) {
	targetValue, err := a.decodeStringBigNumber(hexString)
	// if that function suceeds, then it was a string representing a decimal number
	if err == nil {
		return targetValue, nil
	}

	targetValue, ok := new(big.Int).SetString(hexString, BaseHex)
	if !ok {
		return nil, fmt.Errorf("invalid hex string to decode to uint64")
	}

	return targetValue, nil
}

func (a *abiData) decodeInt8(hexString string) (int8, error) {
	const BitSize = 8

	targetValue, err := a.decodeInt(hexString, BitSize)
	if err != nil {
		return 0, err
	}

	targetValueI8 := int8(*targetValue)
	return targetValueI8, nil
}

func (a *abiData) decodeInt16(hexString string) (int16, error) {
	const BitSize = 16

	targetValue, err := a.decodeInt(hexString, BitSize)
	if err != nil {
		return 0, err
	}

	targetValueI16 := int16(*targetValue)
	return targetValueI16, nil
}

func (a *abiData) decodeInt32(hexString string) (int32, error) {
	const BitSize = 32

	targetValue, err := a.decodeInt(hexString, BitSize)
	if err != nil {
		return 0, err
	}

	targetValueI32 := int32(*targetValue)
	return targetValueI32, nil
}

func (a *abiData) decodeInt64(hexString string) (int64, error) {
	const BitSize = 64

	targetValue, err := a.decodeInt(hexString, BitSize)
	if err != nil {
		return 0, err
	}

	targetValueI64 := int64(*targetValue)
	return targetValueI64, nil
}

func (a *abiData) decodeBigInt(hexString string) (*big.Int, error) {
	targetValue, err := a.decodeStringBigNumber(hexString)
	// if that function suceeds, then it was a string representing a decimal number
	if err == nil {
		return targetValue, nil
	}

	switch len(hexString) {
	case U8HexLength:
		decoded, err := a.decodeInt8(hexString)
		if err != nil {
			return nil, err
		}
		return big.NewInt(int64(decoded)), nil
	case U16HexLength:
		decoded, err := a.decodeInt16(hexString)
		if err != nil {
			return nil, err
		}
		return big.NewInt(int64(decoded)), nil
	case U32HexLength:
		decoded, err := a.decodeInt32(hexString)
		if err != nil {
			return nil, err
		}
		return big.NewInt(int64(decoded)), nil
	case U64HexLength:
		decoded, err := a.decodeInt64(hexString)
		if err != nil {
			return nil, err
		}
		return big.NewInt(int64(decoded)), nil
	default:
		return nil, fmt.Errorf("invalid hex string to decode to BigInt: %s", hexString)
	}
}

func (a *abiData) decodeStringBigNumber(hexString string) (*big.Int, error) {
	targetString, err := a.decodeString(hexString)
	if err != nil {
		return nil, err
	}

	targetValue, ok := new(big.Int).SetString(targetString, 10)
	if !ok {
		return nil, fmt.Errorf("invalid hex string")
	}

	return targetValue, nil
}

func (a *abiData) fixSignedIntOverflow(rawValue *uint64, hexLength int) (*uint, error) {
	var parsedValue uint

	switch hexLength {
	case U8HexLength:
		parsedValue = uint(int8(*rawValue))
		return &parsedValue, nil
	case U16HexLength:
		parsedValue = uint(int16(*rawValue))
		return &parsedValue, nil
	case U32HexLength:
		parsedValue = uint(int32(*rawValue))
		return &parsedValue, nil
	case U64HexLength:
		parsedValue = uint(int64(*rawValue))
		return &parsedValue, nil
	default:
		return nil, fmt.Errorf("invalid hex length %v", hexLength)
	}
}
