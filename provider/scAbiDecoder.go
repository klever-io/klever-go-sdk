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
)

const (
	U64HexLength int = 16
	U32HexLength int = 8
	U16HexLength int = 4
	U8HexLength  int = 2
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

	// return nil, nil
}

func (a *abiData) decodeSingleValue(hexValue string, vType string) (interface{}, error) {
	switch vType {
	case "u64":
		return a.decodeUint64(hexValue)
	case "u32":
		return a.decodeUint32(hexValue)
	case "u16":
		return a.decodeUint16(hexValue)
	case "u8":
		return a.decodeUint8(hexValue)
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
	}
	return nil, fmt.Errorf("Please implement me T-T")
}

func (a *abiData) decodeUint(hexValue string) (*uint64, error) {
	// decode List/Option/tuple/variadic will "cut" the original string
	// value := (*hex)[:length]
	// *hex = (*hex)[length:]

	uintValue, err := strconv.ParseUint(hexValue, 16, 64)
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

func (a *abiData) decodeString(hexValue string) (*string, error) {
	bytes, err := hex.DecodeString(hexValue)

	if err != nil {
		return nil, err
	}

	convertedString := string(bytes)

	return &convertedString, nil
}

func (a *abiData) decodeBigInt(hexString string) (*big.Int, error) {
	if len(hexString) >= U64HexLength { // it will be a string representing a decimal number
		targetString, err := a.decodeString(hexString)
		if err != nil {
			return nil, err
		}

		targetBig, ok := new(big.Int).SetString(*targetString, 10)
		if !ok {
			return nil, fmt.Errorf("invalid hex string")
		}

		return targetBig, nil
	}

	targetBig, ok := new(big.Int).SetString(hexString, 16)
	if !ok {
		return nil, fmt.Errorf("invalid hex string")
	}

	maxInt64big := big.NewInt(math.MaxInt64)
	if targetBig.Cmp(maxInt64big) > 0 {
		targetBig.Sub(targetBig, new(big.Int).Lsh(big.NewInt(1), 64))
	}

	return targetBig, nil
}
