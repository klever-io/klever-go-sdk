package provider

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/klever-io/klever-go-sdk/core/address"
	"github.com/klever-io/klever-go-sdk/provider/utils"
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

type vmOutputData struct {
	Endpoints []endpoint           `json:"endpoints"`
	Types     map[string]typeInfos `json:"types"`
	AbiLoaded bool
}

type VMOutputData interface {
	ParseHex(endpoint string, hex []string) (interface{}, error)
	ParseQuery(endpoint string, base64 []string) (interface{}, error)
	LoadAbi(r io.Reader) error
}

func NewVMOutputHandler() VMOutputData {
	return &vmOutputData{}
}

func (kc *kleverChain) NewScOutputParser() VMOutputData {
	return NewVMOutputHandler()
}

func (a *vmOutputData) LoadAbi(r io.Reader) error {
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

// only for single result outputs
func (a *vmOutputData) ParseHex(endpoint string, hexData []string) (interface{}, error) {
	if !a.AbiLoaded {
		return nil, fmt.Errorf("before parse any value load your abi with `LoadAbi`")
	}

	endpointIndex, err := a.findEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	if len(hexData) == 1 {
		return a.doParse(&hexData[0], a.Endpoints[*endpointIndex].Outputs[0].Type, 0)
	}

	var parsedValues []interface{}

	var outputIndex int
	isMultipleOutputs := len(a.Endpoints[*endpointIndex].Outputs) > 1

	for i, h := range hexData {
		if isMultipleOutputs {
			outputIndex = i
		}
		parsed, err := a.doParse(&h, a.Endpoints[*endpointIndex].Outputs[outputIndex].Type, 0)
		if err != nil {
			return nil, fmt.Errorf("error decoding hex value: %w", err)
		}
		parsedValues = append(parsedValues, parsed)
	}

	return parsedValues, nil
}

// mainly for multi result outputs
func (a *vmOutputData) ParseQuery(endpoint string, base64Data []string) (interface{}, error) {
	var hexData []string

	for _, data := range base64Data {
		dataBytes, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return nil, fmt.Errorf("invalid base64 string: %w", err)
		}
		hexData = append(hexData, hex.EncodeToString(dataBytes))
	}

	return a.ParseHex(endpoint, hexData)
}

func (a *vmOutputData) findEndpoint(endpointName string) (*int, error) {
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

func (a *vmOutputData) doParse(hexRef *string, fullType string, trim int) (interface{}, error) {
	wrapperType, valueType := utils.SplitTypes(fullType)

	return a.selectParser(hexRef, wrapperType, valueType, trim)
}

func (a *vmOutputData) handleTrim(hexRef *string, trim int) string {
	hexToParse := (*hexRef)[:trim]

	if trim == 0 {
		hexToParse = *hexRef
	}

	*hexRef = (*hexRef)[trim:]

	return hexToParse
}

func (a *vmOutputData) selectParser(
	hexRef *string,
	typeWrapper, valueType string,
	trim int,
) (interface{}, error) {
	switch typeWrapper {
	case utils.List:
		return a.parseList(hexRef, valueType)
	case utils.Option:
		return a.parseOption(hexRef, valueType)
	case utils.Tuple:
		return a.parseTuple(hexRef, valueType)
	case utils.Variadic:
		return a.parseVariadic(hexRef, valueType)
	default:
		return a.parseSingleValue(hexRef, valueType, trim)
	}
}

func (a *vmOutputData) parseSingleValue(hexRef *string, valueType string, trim int) (interface{}, error) {
	switch valueType {
	case utils.Int8:
		parsedValue, err := a.parseInt(hexRef, utils.HexLength8Bits)
		return int8(parsedValue), err
	case utils.Int16:
		parsedValue, err := a.parseInt(hexRef, utils.HexLength16Bits)
		return int16(parsedValue), err
	case utils.Int32, utils.Isize:
		parsedValue, err := a.parseInt(hexRef, utils.HexLength32Bits)
		return int32(parsedValue), err
	case utils.Int64:
		parsedValue, err := a.parseInt(hexRef, utils.HexLength64Bits)
		return int64(parsedValue), err
	case utils.Uint8:
		parsedValue, err := a.parseUint(hexRef, utils.HexLength8Bits)
		return uint8(parsedValue), err
	case utils.Uint16:
		parsedValue, err := a.parseUint(hexRef, utils.HexLength16Bits)
		return uint16(parsedValue), err
	case utils.Uint32, utils.Usize:
		parsedValue, err := a.parseUint(hexRef, utils.HexLength32Bits)
		return uint32(parsedValue), err
	case utils.Uint64:
		parsedValue, err := a.parseUint(hexRef, utils.HexLength64Bits)
		return uint64(parsedValue), err
	case utils.Address:
		return a.parseAddress(hexRef)
	case utils.BigInt:
		return a.parseBigInt(hexRef, trim)
	case utils.BigUint:
		return a.parseBigUint(hexRef, trim)
	case utils.BigFloat:
		return a.parseBigFloat(hexRef, trim)
	case utils.Boolean:
		return a.parseBoolean(hexRef), nil
	case
		utils.ManagedBuffer,
		utils.TokenIdentifier,
		utils.Bytes,
		utils.BoxedBytes,
		utils.String,
		utils.StrRef,
		utils.VecU8,
		utils.SliceU8:
		return a.parseString(hexRef, trim)
	default:
		return a.parseStruct(hexRef, valueType)
	}
}

func (a *vmOutputData) parseString(hexRef *string, trim int) (string, error) {
	hexToParse := a.handleTrim(hexRef, trim*2)

	bytes, err := hex.DecodeString(hexToParse)

	return string(bytes), err
}

func (a *vmOutputData) parseAddress(hexRef *string) (string, error) {
	hexToParse := (*hexRef)[:utils.AddressHexLen]
	*hexRef = (*hexRef)[utils.AddressHexLen:]

	parsedAddress, err := address.NewAddressFromHex(hexToParse)

	return parsedAddress.Bech32(), err
}

func (a *vmOutputData) parseBaseUint(hexRef *string, bitSize int) (*uint64, error) {
	uintValue, err := strconv.ParseUint(*hexRef, utils.BaseHex, bitSize)

	return &uintValue, err
}

func (a *vmOutputData) getDynamicTrim(hexRef *string, trimSize int) string {
	var hexToParse string

	if len(*hexRef) > trimSize {
		hexToParse = a.handleTrim(hexRef, trimSize)

		return hexToParse
	}

	hexToParse = *hexRef
	*hexRef = ""

	return hexToParse
}

func (a *vmOutputData) parseUint(hexRef *string, bitsHexLen int) (uint, error) {
	hexToParse := a.getDynamicTrim(hexRef, bitsHexLen)

	switch rawHexLen := len(hexToParse); {
	case rawHexLen <= utils.HexLength8Bits:
		targetValue, err := a.parseBaseUint(&hexToParse, utils.Bits8)
		return uint(*targetValue), err
	case rawHexLen <= utils.HexLength16Bits:
		targetValue, err := a.parseBaseUint(&hexToParse, utils.Bits16)
		return uint(*targetValue), err
	case rawHexLen <= utils.HexLength32Bits:
		targetValue, err := a.parseBaseUint(&hexToParse, utils.Bits32)
		return uint(*targetValue), err
	case rawHexLen <= utils.HexLength64Bits:
		targetValue, err := a.parseBaseUint(&hexToParse, utils.Bits64)
		return uint(*targetValue), err
	default:
		return 0, fmt.Errorf("invalid hex string %s to parse to uint", hexToParse)
	}
}

func (a *vmOutputData) parseInt(hexRef *string, bitsHexLen int) (int, error) {
	hexToParse := a.getDynamicTrim(hexRef, bitsHexLen)

	switch rawHexLen := len(hexToParse); {
	case rawHexLen <= utils.HexLength8Bits:
		targetValue, err := a.parseBaseUint(&hexToParse, utils.Bits8)
		return int(int8(*targetValue)), err
	case rawHexLen <= utils.HexLength16Bits:
		targetValue, err := a.parseBaseUint(&hexToParse, utils.Bits16)
		return int(int16(*targetValue)), err
	case rawHexLen <= utils.HexLength32Bits:
		targetValue, err := a.parseBaseUint(&hexToParse, utils.Bits32)
		return int(int32(*targetValue)), err
	case rawHexLen <= utils.HexLength64Bits:
		targetValue, err := a.parseBaseUint(&hexToParse, utils.Bits64)
		return int(int64(*targetValue)), err
	default:
		return 0, fmt.Errorf("invalid hex string %s to parse to int", hexToParse)
	}
}

func (a *vmOutputData) parseBigUint(hexRef *string, trim int) (*big.Int, error) {
	hexToParse := a.handleTrim(hexRef, trim*2)

	targetValue, err := a.parseStringBigNumber(&hexToParse)
	// if that function suceeds, then it was a string representing a decimal big number
	if err == nil {
		return targetValue, nil
	}

	targetValue, ok := new(big.Int).SetString(hexToParse, utils.BaseHex)
	if !ok {
		return nil, fmt.Errorf("invalid hex string %s to parse to uint64", hexToParse)
	}

	return targetValue, nil
}

func (a *vmOutputData) parseBigInt(hexRef *string, trim int) (*big.Int, error) {
	hexToParse := a.handleTrim(hexRef, trim*2)

	targetValueFromString, err := a.parseStringBigNumber(&hexToParse)
	// if that function suceeds, then it was a string representing a decimal number
	if err == nil {
		return targetValueFromString, nil
	}

	targetValueFromInt128, err := a.handleBigIntTill128(&hexToParse)
	if err == nil {
		return targetValueFromInt128, nil
	}

	targetValue, err := a.parseInt(&hexToParse, len(hexToParse))
	if err != nil {
		return nil, fmt.Errorf("invelid hex string %s to parse to big int", hexToParse)
	}

	return big.NewInt(int64(targetValue)), nil
}

func (a *vmOutputData) parseStringBigNumber(hexRef *string) (*big.Int, error) {
	targetString, err := a.parseString(hexRef, 0)
	if err != nil {
		return nil, err
	}

	targetValue, ok := new(big.Int).SetString(targetString, utils.BaseDecimal)
	if !ok {
		return nil, fmt.Errorf("invalid hex string %s to parse to BigInt", (*hexRef))
	}

	return targetValue, nil
}

func (a *vmOutputData) handleBigIntTill128(hexRef *string) (*big.Int, error) {
	rawValue, ok := new(big.Int).SetString(*hexRef, utils.BaseHex)
	if !ok {
		return nil, fmt.Errorf("invalid hex string to parse to BigInt: %s", *hexRef)
	}

	valueBits := len(*hexRef) * utils.BitsByHexDigit

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

func (a *vmOutputData) parseBoolean(hexRef *string) bool {
	const BooleanLength int = 2

	hexToParse := a.getDynamicTrim(hexRef, BooleanLength)
	return hexToParse == "01"
}

func (a *vmOutputData) parseBigFloat(hexRef *string, trim int) (interface{}, error) {
	hexToParse := a.handleTrim(hexRef, trim*2)

	hexBytes, err := hex.DecodeString(hexToParse)
	if err != nil {
		return nil, fmt.Errorf("error decoding hex string to bytes: %w", err)
	}

	parsedValue := new(big.Float)

	if err := parsedValue.GobDecode(hexBytes); err != nil {
		return nil, fmt.Errorf("error decoding big float: %w", err)
	}

	return parsedValue, nil
}

func (a *vmOutputData) parseOption(hexRef *string, valueType string) (interface{}, error) {
	if *hexRef == "" {
		return nil, nil
	}

	const OptionStringLenght int = 2
	a.getDynamicTrim(hexRef, OptionStringLenght)

	hasListPrefix := strings.HasPrefix(valueType, utils.List)

	var trim int

	if utils.IsDynamicLengthType(valueType) || hasListPrefix {
		calculatedTrim, err := a.getFixedTrim(hexRef)
		if err != nil {
			return nil, fmt.Errorf("error while triming option hex string %w", err)
		}

		trim = calculatedTrim
	}

	return a.doParse(hexRef, valueType, trim)
}

func (a *vmOutputData) getFixedTrim(hexRef *string) (int, error) {
	trimStringHex := (*hexRef)[:utils.LengthHexSizer]

	trimRef, err := a.parseBaseUint(&trimStringHex, utils.LengthHexSizer*utils.BitsByHexDigit)

	*hexRef = (*hexRef)[utils.LengthHexSizer:]

	return int(*trimRef), err
}

func (a *vmOutputData) parseList(hexRef *string, valueType string) (interface{}, error) {
	var result []interface{}
	for len((*hexRef)) > 0 {
		parsed, err := a.handleList(hexRef, valueType)
		if err != nil {
			return nil, fmt.Errorf("error decoding list item: %w", err)
		}

		result = append(result, parsed)
	}

	return result, nil
}

func (a *vmOutputData) handleList(hexRef *string, valueType string) (interface{}, error) {
	wrapperType, coreType := utils.SplitTypes(valueType)

	if wrapperType == utils.List {
		listTrim, err := a.getFixedTrim(hexRef)

		if err != nil {
			return nil, fmt.Errorf("error getting the list trim: %w", err)
		}

		return a.parseNestedList(hexRef, coreType, listTrim)
	}

	var valueTrim int

	if utils.IsDynamicLengthType(valueType) {
		calculatedTrim, err := a.getFixedTrim(hexRef)
		if err != nil {
			return nil, fmt.Errorf("error getting the list item trim: %w", err)
		}

		valueTrim = calculatedTrim
	}

	return a.parseSingleValue(hexRef, coreType, valueTrim)
}

func (a *vmOutputData) parseNestedList(hexRef *string, valueType string, limit int) (interface{}, error) {
	var result []interface{}
	for i := 0; i < limit; i++ {
		parsed, err := a.handleList(hexRef, valueType)
		if err != nil {
			return nil, fmt.Errorf("error decoding nested list: %w", err)
		}

		result = append(result, parsed)
	}

	return result, nil
}

func (a *vmOutputData) parseTuple(hexRef *string, valueType string) ([]interface{}, error) {
	var result []interface{}

	types := utils.SplitTupleTypes(valueType)

	for _, t := range types {
		if strings.HasPrefix(t, utils.List) {
			parsedList, err := a.handleList(hexRef, t)
			if err != nil {
				return nil, err
			}

			result = append(result, parsedList)
			continue
		}

		var valueTrim int

		if utils.IsDynamicLengthType(t) {
			calculatedTrim, err := a.getFixedTrim(hexRef)
			if err != nil {
				return nil, fmt.Errorf("error getting the list item trim: %w", err)
			}

			valueTrim = calculatedTrim
		}

		parsedValue, err := a.doParse(hexRef, t, valueTrim)
		if err != nil {
			return nil, fmt.Errorf("error decoding tuple value %w", err)
		}

		result = append(result, parsedValue)
	}

	return result, nil
}

func (a *vmOutputData) parseVariadic(hexRef *string, valueType string) (interface{}, error) {
	parsedValue, err := a.doParse(hexRef, valueType, 0)
	if err != nil {
		return nil, fmt.Errorf("error trying to parse variadic value: %w", err)
	}

	return parsedValue, nil
}

func (a *vmOutputData) parseStruct(hexRef *string, valueType string) (map[string]interface{}, error) {
	typeDef, exists := a.Types[valueType]
	if !exists {
		return nil, fmt.Errorf("type %s not found in provided abi", valueType)
	}

	result := make(map[string]interface{})

	for _, field := range typeDef.Fields {
		if strings.HasPrefix(field.Type, utils.List) {
			parsedList, err := a.handleList(hexRef, field.Type)
			if err != nil {
				return nil, fmt.Errorf("error %w decoding list value of key %s of custom type %s", err, field.Type, valueType)
			}

			result[field.Name] = parsedList
			continue
		}

		var trim int

		if utils.IsDynamicLengthType(field.Type) {
			calculatedTrim, err := a.getFixedTrim(hexRef)
			if err != nil {
				return nil, fmt.Errorf("error while triming option hex string %w", err)
			}

			trim = calculatedTrim
		}

		parsedValue, err := a.doParse(hexRef, field.Type, trim)
		if err != nil {
			return nil, fmt.Errorf("error %w decoding value of key %s of custom type %s", err, field.Type, valueType)
		}

		result[field.Name] = parsedValue
	}

	return result, nil
}
