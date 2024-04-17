package provider

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/klever-io/klever-go-sdk/core/address"
)

const (
	// Hex length
	HexLength8Bits  int = 2
	HexLength16Bits int = 4
	HexLength32Bits int = 8
	HexLength64Bits int = 16

	// Bits count
	Bits8  int = 8
	Bits16 int = 16
	Bits32 int = 32
	Bits64 int = 64

	// Numerical bases
	BaseHex     int = 16
	BaseDecimal int = 10

	// Types wrappers
	List     string = "List"
	Option   string = "Option"
	Tuple    string = "tuple"
	Variadic string = "variadic"

	// Possible Types
	Int8            string = "i8"
	Uint8           string = "u8"
	Int16           string = "i16"
	Uint16          string = "u16"
	Int32           string = "i32"
	Uint32          string = "u32"
	Int64           string = "i64"
	Uint64          string = "u64"
	BigInt          string = "BigInt"
	BigUint         string = "BigUint"
	Address         string = "Address"
	Boolean         string = "bool"
	ManagedBuffer   string = "ManagedBuffer"
	TokenIdentifier string = "TokenIdentifier"
	Bytes           string = "bytes"
	BoxedBytes      string = "BoxedBytes"
	String          string = "String"
	StrRef          string = "&str"
	VecU8           string = "Vec<u8>"
	SliceU8         string = "&[u8]"

	LengthHexSizer int = 8

	BitsByHexDigit int = 4

	AddressHexLen int = 64
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
	Type   string  `json:"type"` // struct, tuple, variadic, List, Option
	Fields []field `json:"fields"`
}

type abiData struct {
	Endpoints []endpoint           `json:"endpoints"`
	Types     map[string]typeInfos `json:"types"`
	AbiLoaded bool
}

type AbiData interface {
	Decode(endpoint, hex string) (interface{}, error)
	LoadAbi(r io.Reader) error
}

func NewSCAbiHandler() AbiData {
	return &abiData{}
}

func (a *abiData) LoadAbi(r io.Reader) error {
	jsonBytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonBytes, a); err != nil {
		return err
	}

	a.AbiLoaded = true

	return nil
}

func (a *abiData) Decode(endpoint, hex string) (interface{}, error) {
	if !a.AbiLoaded {
		return nil, fmt.Errorf("before decode any value load your abi with `LoadAbi`")
	}

	endpointIndex, err := a.findEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	return a.doDecode(&hex, a.Endpoints[*endpointIndex].Outputs[0].Type, 0)
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

func (a *abiData) splitTypes(fullType string) (string, string) {
	index := strings.Index(fullType, "<")
	if index == -1 {
		return "", fullType
	}

	wrapperType := fullType[:index]
	valueType := fullType[index+1:]
	valueType = valueType[:len(valueType)-1]

	return wrapperType, valueType
}

func (a *abiData) doDecode(hexRef *string, fullType string, trim int) (interface{}, error) {
	wrapperType, valueType := a.splitTypes(fullType)

	return a.selectDecoder(hexRef, wrapperType, valueType, trim)
}

func (a *abiData) handleTrim(hexRef *string, trim int) string {
	hexToDecode := (*hexRef)[:trim]

	if trim == 0 {
		hexToDecode = *hexRef
	}

	*hexRef = (*hexRef)[trim:]

	return hexToDecode
}

func (a *abiData) selectDecoder(
	hexRef *string,
	typeWrapper, valueType string,
	trim int,
) (interface{}, error) {
	switch typeWrapper {
	case List:
		return a.decodeList(hexRef, valueType)
	case Option:
		return a.decodeOption(hexRef, valueType)
	case Tuple:
		return nil, fmt.Errorf("tuple")
	case Variadic:
		return nil, fmt.Errorf("variadic")
	default:
		return a.decodeSingleValue(hexRef, valueType, trim)
	}
}

func (a *abiData) decodeSingleValue(hexRef *string, valueType string, trim int) (interface{}, error) {
	switch valueType {
	case Int8:
		return a.decodeInt8(hexRef)
	case Int16:
		return a.decodeInt16(hexRef)
	case Int32:
		return a.decodeInt32(hexRef)
	case Int64:
		return a.decodeInt64(hexRef)
	case Uint8:
		return a.decodeUint8(hexRef)
	case Uint16:
		return a.decodeUint16(hexRef)
	case Uint32:
		return a.decodeUint32(hexRef)
	case Uint64:
		return a.decodeUint64(hexRef)
	case Address:
		return a.decodeAddress(hexRef)
	case BigInt:
		return a.decodeBigInt(hexRef, trim)
	case Boolean:
		return (*hexRef) == "01", nil
	case BigUint:
		return a.decodeBigUint(hexRef, trim)
	case
		ManagedBuffer,
		TokenIdentifier,
		Bytes,
		BoxedBytes,
		String,
		StrRef,
		VecU8,
		SliceU8:
		return a.decodeString(hexRef, trim)
	default:
		return nil, fmt.Errorf("invalid type %s", valueType)
	}
}

func (a *abiData) decodeString(hexRef *string, trim int) (string, error) {
	hexToDecode := a.handleTrim(hexRef, trim*2)

	bytes, err := hex.DecodeString(hexToDecode)

	return string(bytes), err
}

func (a *abiData) decodeAddress(hexRef *string) (string, error) {
	hexToDecode := (*hexRef)[:AddressHexLen]
	*hexRef = (*hexRef)[AddressHexLen:]

	decodedAddress, err := address.NewAddressFromHex(hexToDecode)

	return decodedAddress.Bech32(), err
}

func (a *abiData) decodeInt(hexRef *string, bitSize int) (*uint64, error) {
	uintValue, err := strconv.ParseUint(*hexRef, BaseHex, bitSize)

	return &uintValue, err
}

func (a *abiData) handleIntsTrim(hexRef *string, bitSize int) string {
	var hexToDecode string

	if len(*hexRef) > bitSize {
		hexToDecode = a.handleTrim(hexRef, bitSize)

		return hexToDecode
	}

	hexToDecode = *hexRef
	*hexRef = ""

	return hexToDecode
}

func (a *abiData) decodeUint8(hexRef *string) (uint8, error) {
	hexToDecode := a.handleIntsTrim(hexRef, HexLength8Bits)

	targetValue, err := a.decodeInt(&hexToDecode, Bits8)

	return uint8(*targetValue), err
}

func (a *abiData) decodeUint16(hexRef *string) (uint16, error) {
	hexToDecode := a.handleIntsTrim(hexRef, HexLength16Bits)

	targetValue, err := a.decodeInt(&hexToDecode, Bits16)

	return uint16(*targetValue), err
}

func (a *abiData) decodeUint32(hexRef *string) (uint32, error) {
	hexToDecode := a.handleIntsTrim(hexRef, HexLength32Bits)

	targetValue, err := a.decodeInt(&hexToDecode, Bits32)

	return uint32(*targetValue), err
}

func (a *abiData) decodeUint64(hexRef *string) (uint64, error) {
	hexToDecode := a.handleIntsTrim(hexRef, HexLength64Bits)

	targetValue, err := a.decodeInt(&hexToDecode, Bits64)

	return uint64(*targetValue), err
}

func (a *abiData) decodeInt8(hexRef *string) (int8, error) {
	hexToDecode := a.handleIntsTrim(hexRef, HexLength8Bits)

	targetValue, err := a.decodeInt(&hexToDecode, Bits8)

	return int8(*targetValue), err
}

func (a *abiData) decodeInt16(hexRef *string) (int16, error) {
	hexToDecode := a.handleIntsTrim(hexRef, HexLength16Bits)

	targetValue, err := a.decodeInt(&hexToDecode, Bits16)

	return int16(*targetValue), err
}

func (a *abiData) decodeInt32(hexRef *string) (int32, error) {
	hexToDecode := a.handleIntsTrim(hexRef, HexLength32Bits)

	targetValue, err := a.decodeInt(&hexToDecode, Bits32)

	return int32(*targetValue), err
}

func (a *abiData) decodeInt64(hexRef *string) (int64, error) {
	hexToDecode := a.handleIntsTrim(hexRef, HexLength64Bits)

	targetValue, err := a.decodeInt(&hexToDecode, Bits64)

	return int64(*targetValue), err
}

func (a *abiData) decodeBigUint(hexRef *string, trim int) (*big.Int, error) {
	hexToDecode := a.handleTrim(hexRef, trim*2)

	targetValue, err := a.decodeStringBigNumber(&hexToDecode)
	// if that function suceeds, then it was a string representing a decimal big number
	if err == nil {
		return targetValue, nil
	}

	targetValue, ok := new(big.Int).SetString(hexToDecode, BaseHex)
	if !ok {
		return nil, fmt.Errorf("invalid hex string %s to decode to uint64", hexToDecode)
	}

	return targetValue, nil
}

func (a *abiData) decodeBigInt(hexRef *string, trim int) (*big.Int, error) {
	hexToDecode := a.handleTrim(hexRef, trim*2)

	targetValueFromString, err := a.decodeStringBigNumber(&hexToDecode)
	// if that function suceeds, then it was a string representing a decimal number
	if err == nil {
		return targetValueFromString, nil
	}

	targetValueFromInt128, err := a.handleBigIntTill128(&hexToDecode)
	if err == nil {
		return targetValueFromInt128, nil
	}

	switch hexLen := len(hexToDecode); {
	case hexLen <= HexLength8Bits:
		decoded, err := a.decodeInt8(&hexToDecode)
		return big.NewInt(int64(decoded)), err
	case hexLen <= HexLength16Bits:
		decoded, err := a.decodeInt16(&hexToDecode)
		return big.NewInt(int64(decoded)), err
	case hexLen <= HexLength32Bits:
		decoded, err := a.decodeInt32(&hexToDecode)
		return big.NewInt(int64(decoded)), err
	case hexLen <= HexLength64Bits:
		decoded, err := a.decodeInt64(&hexToDecode)
		return big.NewInt(int64(decoded)), err
	default:
		return nil, fmt.Errorf("invalid hex string %s to decode to BigInt", *hexRef)
	}
}

func (a *abiData) decodeStringBigNumber(hexRef *string) (*big.Int, error) {
	targetString, err := a.decodeString(hexRef, 0)
	if err != nil {
		return nil, err
	}

	targetValue, ok := new(big.Int).SetString(targetString, BaseDecimal)
	if !ok {
		return nil, fmt.Errorf("invalid hex string %s to decode to BigInt", (*hexRef))
	}

	return targetValue, nil
}

func (a *abiData) handleBigIntTill128(hexRef *string) (*big.Int, error) {
	rawValue, ok := new(big.Int).SetString(*hexRef, BaseHex)
	if !ok {
		return nil, fmt.Errorf("invalid hex string to decode to BigInt: %s", *hexRef)
	}

	valueBits := len(*hexRef) * BitsByHexDigit

	two := big.NewInt(2)
	twoToTheNth := new(big.Int).Exp(two, big.NewInt(int64(valueBits-1)), nil) // 2^valueBits
	one := big.NewInt(1)

	MaxIntNBits := new(big.Int).Sub(twoToTheNth, one)

	if rawValue.Cmp(new(big.Int).SetUint64(math.MaxUint64)) == 1 && rawValue.Cmp(MaxIntNBits) == -1 {
		return rawValue, nil
	}

	if rawValue.Cmp(MaxIntNBits) == 1 {
		parsedValue := rawValue.Sub(rawValue, new(big.Int).Lsh(one, uint(valueBits)))
		return parsedValue, nil
	}

	return nil, fmt.Errorf("%v range is lower range than 128 bits", rawValue)
}

func (a *abiData) decodeOption(hexRef *string, valueType string) (interface{}, error) {
	hexToDecode := a.handleTrim(hexRef, 2)

	if hexToDecode == "00" {
		return nil, nil
	}

	return a.doDecode(hexRef, valueType, 0)
}

func stringMatch(items []string, cmpFunc func(string) bool) bool {
	for _, item := range items {
		if cmpFunc(item) {
			return true
		}
	}
	return false
}

func (a *abiData) getListTrim(hexRef *string) (int, error) {
	trimStringHex := (*hexRef)[:LengthHexSizer]

	trimRef, err := a.decodeInt(&trimStringHex, LengthHexSizer*BitsByHexDigit)

	*hexRef = (*hexRef)[LengthHexSizer:]

	return int(*trimRef), err
}

func (a *abiData) decodeList(hexRef *string, valueType string) (interface{}, error) {
	var result []interface{}
	for len((*hexRef)) > 0 {
		decoded, err := a.handleList(hexRef, valueType)
		if err != nil {
			return nil, fmt.Errorf("error decoding list item: %w", err)
		}

		result = append(result, decoded)
	}

	return result, nil
}

func (a *abiData) handleList(hexRef *string, valueType string) (interface{}, error) {
	wrapperType, innerType := a.splitTypes(valueType)

	if wrapperType == List {
		listTrim, err := a.getListTrim(hexRef)

		if err != nil {
			return nil, fmt.Errorf("error getting the list trim: %w", err)
		}

		return a.decodeNestedList(hexRef, innerType, listTrim)
	}

	var valueTrim int

	dynamicLengthTypes := []string{
		ManagedBuffer, TokenIdentifier, Bytes, BoxedBytes,
		String, StrRef, VecU8, SliceU8, BigInt, BigUint,
	}

	if stringMatch(dynamicLengthTypes, func(s string) bool { return valueType == s }) {
		calculatedTrim, err := a.getListTrim(hexRef)
		if err != nil {
			return nil, fmt.Errorf("error getting the list item trim: %w", err)
		}

		valueTrim = calculatedTrim
	}

	return a.decodeSingleValue(hexRef, innerType, valueTrim)
}

func (a *abiData) decodeNestedList(hexRef *string, valueType string, limit int) (interface{}, error) {
	var result []interface{}
	for i := 0; i < limit; i++ {
		decoded, err := a.handleList(hexRef, valueType)
		if err != nil {
			return nil, fmt.Errorf("error decoding nested list: %w", err)
		}

		result = append(result, decoded)
	}

	return result, nil
}
