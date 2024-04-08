package provider

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"os"
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

func (a *abiData) Decode(abiPath, endpoint, hex string) (interface{}, error) {
	if err := a.loadAbi(abiPath); err != nil {
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

func (a *abiData) loadAbi(abiPath string) error {
	jsonBytes, err := os.ReadFile(abiPath)
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

func (a *abiData) decodeString(hexValue string) (*string, error) {
	bytes, err := hex.DecodeString(hexValue)

	if err != nil {
		return nil, err
	}

	convertedString := string(bytes)

	return &convertedString, nil
}

func (a *abiData) decodeUint(hexValue string) (*uint64, error) {
	// decode List/Option/tuple/variadic will "cut" the original string
	// value := (*hex)[:length]
	// *hex = (*hex)[length:]

	uintValue, err := strconv.ParseUint(hexValue, BaseHex, 64)
	if err != nil {
		return nil, err
	}

	return &uintValue, nil
}

func (a *abiData) decodeUint8(hexValue string) (*uint8, error) {
	targetValue, err := a.decodeUint(hexValue)
	uint8Decoded := uint8(*targetValue)

	return &uint8Decoded, err
}

func (a *abiData) decodeUint16(hexValue string) (*uint16, error) {
	targetValue, err := a.decodeUint(hexValue)
	uint16Decoded := uint16(*targetValue)

	return &uint16Decoded, err
}

func (a *abiData) decodeUint32(hexValue string) (*uint32, error) {
	targetValue, err := a.decodeUint(hexValue)
	uint32Decoded := uint32(*targetValue)

	return &uint32Decoded, err
}

func (a *abiData) decodeUint64(hexValue string) (*uint64, error) {
	return a.decodeUint(hexValue)
}

func (a *abiData) decodeBigUint(hexValue string) (*big.Int, error) {
	if len(hexValue) > BaseHex {
		return a.decodeStringBigNumber(hexValue)
	}

	targetValue, ok := new(big.Int).SetString(hexValue, BaseHex)
	if !ok {
		return nil, fmt.Errorf("invalid hex string to decode to uint64")
	}

	return targetValue, nil
}

func (a *abiData) decodeInt8(hexString string) (*int8, error) {
	const BitSize = 8

	parsedValue, err := strconv.ParseInt(hexString, BaseHex, BitSize)
	if err != nil {
		return nil, err
	}

	if uint8(parsedValue) > math.MaxInt8 {
		ptrtTrgetValue, err := a.fixIntOverflow(hexString, BitSize)
		if err != nil {
			return nil, err
		}

		targetValue := int8(*ptrtTrgetValue)
		return &targetValue, nil
	}

	targetValue := int8(parsedValue)
	return &targetValue, nil
}

func (a *abiData) decodeInt16(hexString string) (*int16, error) {
	const BitSize = 16

	parsedValue, err := strconv.ParseInt(hexString, BaseHex, BitSize)
	if err != nil {
		return nil, err
	}

	if uint16(parsedValue) > math.MaxInt16 {
		ptrtTrgetValue, err := a.fixIntOverflow(hexString, BitSize)
		if err != nil {
			return nil, err
		}

		targetValue := int16(*ptrtTrgetValue)
		return &targetValue, nil
	}

	targetValue := int16(parsedValue)
	return &targetValue, nil
}

func (a *abiData) decodeInt32(hexString string) (*int32, error) {
	const BitSize = 32

	parsedValue, err := strconv.ParseInt(hexString, BaseHex, BitSize)
	if err != nil {
		return nil, err
	}

	if uint32(parsedValue) > math.MaxInt32 {
		ptrtTrgetValue, err := a.fixIntOverflow(hexString, BitSize)
		if err != nil {
			return nil, err
		}

		targetValue := int32(*ptrtTrgetValue)
		return &targetValue, nil
	}

	targetValue := int32(parsedValue)
	return &targetValue, nil
}

func (a *abiData) decodeInt64(hexString string) (*int64, error) {
	const BitSize = 64

	targetValue, err := strconv.ParseUint(hexString, BaseHex, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string to decode to int64 %s: %v", hexString, err)
	}

	if targetValue > math.MaxInt64 {
		ptrtTargetValue, err := a.fixIntOverflow(hexString, BitSize)
		if err != nil {
			return nil, err
		}

		targetValue := int64(*ptrtTargetValue)
		return &targetValue, nil
	}

	targetValueI64 := int64(targetValue)
	return &targetValueI64, nil
}

func (a *abiData) decodeBigInt(hexString string) (*big.Int, error) {
	targetValue, err := a.decodeStringBigNumber(hexString)
	// if that function suceeds, then it was a string representin a decimal number
	if err == nil {
		return targetValue, nil
	}

	switch len(hexString) {
	case U8HexLength:
		decoded, err := a.decodeInt8(hexString)
		if err != nil {
			return nil, err
		}
		return big.NewInt(int64(*decoded)), nil
	case U16HexLength:
		decoded, err := a.decodeInt16(hexString)
		if err != nil {
			return nil, err
		}
		return big.NewInt(int64(*decoded)), nil
	case U32HexLength:
		decoded, err := a.decodeInt32(hexString)
		if err != nil {
			return nil, err
		}
		return big.NewInt(int64(*decoded)), nil
	case U64HexLength:
		decoded, err := a.decodeInt64(hexString)
		if err != nil {
			return nil, err
		}
		return big.NewInt(int64(*decoded)), nil
	default:
		return nil, fmt.Errorf("invalid hex string to decode to BigInt: %s", hexString)
	}
}

func (a *abiData) decodeStringBigNumber(hexString string) (*big.Int, error) {
	targetString, err := a.decodeString(hexString)
	if err != nil {
		return nil, err
	}

	targetValue, ok := new(big.Int).SetString(*targetString, 10)
	if !ok {
		return nil, fmt.Errorf("invalid hex string")
	}

	return targetValue, nil
}

func (a *abiData) fixIntOverflow(hexString string, bitSize int) (*int64, error) {
	if bitSize > 64 {
		return nil, fmt.Errorf("bitSize too large for uint64")
	}

	targetValue, err := strconv.ParseUint(hexString, BaseHex, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string %s: %v", hexString, err)
	}

	targetValue -= 1 << bitSize

	targetValueInt64 := int64(targetValue)

	return &targetValueInt64, nil
}
