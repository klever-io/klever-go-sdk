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
	Types     map[string]typeInfos `json:"types"`
	Endpoints []endpoint           `json:"endpoints"`
	AbiLoaded bool
}

type VMOutputData interface {
	DecodeHex(endpoint string, hex []string) (interface{}, error)
	DecodeQuery(endpoint string, base64 []string) (interface{}, error)
	LoadAbi(r io.Reader) error
}

func NewVMOutputHandler() VMOutputData {
	return &vmOutputData{}
}

func (kc *kleverChain) NewScOutputDecoder() VMOutputData {
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
func (a *vmOutputData) DecodeHex(endpoint string, hexData []string) (interface{}, error) {
	if !a.AbiLoaded {
		return nil, fmt.Errorf("before decode any value load your abi with `LoadAbi`")
	}

	endpointIndex, err := a.findEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	if len(hexData) == 1 {
		return a.doDecode(&hexData[0], a.Endpoints[*endpointIndex].Outputs[0].Type, 0)
	}

	var decodedValues []interface{}

	var outputIndex int
	isMultipleOutputs := len(a.Endpoints[*endpointIndex].Outputs) > 1

	for i, h := range hexData {
		if isMultipleOutputs {
			outputIndex = i
		}
		decoded, err := a.doDecode(&h, a.Endpoints[*endpointIndex].Outputs[outputIndex].Type, 0)
		if err != nil {
			return nil, fmt.Errorf("error decoding hex value: %w", err)
		}
		decodedValues = append(decodedValues, decoded)
	}

	return decodedValues, nil
}

// mainly for multi result outputs
func (a *vmOutputData) DecodeQuery(endpoint string, base64Data []string) (interface{}, error) {
	var hexData []string

	for _, data := range base64Data {
		dataBytes, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return nil, fmt.Errorf("invalid base64 string: %w", err)
		}
		hexData = append(hexData, hex.EncodeToString(dataBytes))
	}

	return a.DecodeHex(endpoint, hexData)
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

func (a *vmOutputData) doDecode(hexRef *string, fullType string, trim int) (interface{}, error) {
	wrapperType, valueType := utils.SplitTypes(fullType)

	return a.selectDecoder(hexRef, wrapperType, valueType, trim)
}

func (a *vmOutputData) handleTrim(hexRef *string, trim int) string {
	hexToDecode := (*hexRef)[:trim]

	if trim == 0 {
		hexToDecode = *hexRef
	}

	*hexRef = (*hexRef)[trim:]

	return hexToDecode
}

func (a *vmOutputData) selectDecoder(
	hexRef *string,
	typeWrapper, valueType string,
	trim int,
) (interface{}, error) {
	switch typeWrapper {
	case utils.List:
		return a.decodeList(hexRef, valueType)
	case utils.Option:
		return a.decodeOption(hexRef, valueType)
	case utils.Tuple:
		return a.decodeTuple(hexRef, valueType)
	case utils.Variadic:
		return a.decodeVariadic(hexRef, valueType)
	default:
		return a.decodeSingleValue(hexRef, valueType, trim)
	}
}

func (a *vmOutputData) decodeSingleValue(
	hexRef *string,
	valueType string,
	trim int,
) (interface{}, error) {
	switch valueType {
	case utils.Int8:
		decodedValue, err := a.decodeInt(hexRef, utils.HexLength8Bits)
		return int8(decodedValue), err
	case utils.Int16:
		decodedValue, err := a.decodeInt(hexRef, utils.HexLength16Bits)
		return int16(decodedValue), err
	case utils.Int32, utils.Isize:
		decodedValue, err := a.decodeInt(hexRef, utils.HexLength32Bits)
		return int32(decodedValue), err
	case utils.Int64:
		decodedValue, err := a.decodeInt(hexRef, utils.HexLength64Bits)
		return int64(decodedValue), err
	case utils.Uint8:
		decodedValue, err := a.decodeUint(hexRef, utils.HexLength8Bits)
		return uint8(decodedValue), err
	case utils.Uint16:
		decodedValue, err := a.decodeUint(hexRef, utils.HexLength16Bits)
		return uint16(decodedValue), err
	case utils.Uint32, utils.Usize:
		decodedValue, err := a.decodeUint(hexRef, utils.HexLength32Bits)
		return uint32(decodedValue), err
	case utils.Uint64:
		decodedValue, err := a.decodeUint(hexRef, utils.HexLength64Bits)
		return uint64(decodedValue), err
	case utils.Address:
		return a.decodeAddress(hexRef)
	case utils.BigInt:
		return a.decodeBigInt(hexRef, trim)
	case utils.BigUint:
		return a.decodeBigUint(hexRef, trim)
	case utils.BigFloat:
		return a.decodeBigFloat(hexRef, trim)
	case utils.Boolean:
		return a.decodeBoolean(hexRef), nil
	case
		utils.ManagedBuffer,
		utils.TokenIdentifier,
		utils.Bytes,
		utils.BoxedBytes,
		utils.String,
		utils.StrRef,
		utils.VecU8,
		utils.SliceU8:
		return a.decodeString(hexRef, trim)
	default:
		return a.decodeStruct(hexRef, valueType)
	}
}

func (a *vmOutputData) decodeString(hexRef *string, trim int) (string, error) {
	hexToDecode := a.handleTrim(hexRef, trim*2)

	bytes, err := hex.DecodeString(hexToDecode)

	return string(bytes), err
}

func (a *vmOutputData) decodeAddress(hexRef *string) (string, error) {
	hexToDecode := (*hexRef)[:utils.AddressHexLen]
	*hexRef = (*hexRef)[utils.AddressHexLen:]

	decodedAddress, err := address.NewAddressFromHex(hexToDecode)

	return decodedAddress.Bech32(), err
}

func (a *vmOutputData) decodeBaseUint(hexRef *string, bitSize int) (*uint64, error) {
	uintValue, err := strconv.ParseUint(*hexRef, utils.BaseHex, bitSize)

	return &uintValue, err
}

func (a *vmOutputData) getDynamicTrim(hexRef *string, trimSize int) string {
	var hexToDecode string

	if len(*hexRef) > trimSize {
		hexToDecode = a.handleTrim(hexRef, trimSize)

		return hexToDecode
	}

	hexToDecode = *hexRef
	*hexRef = ""

	return hexToDecode
}

func (a *vmOutputData) decodeUint(hexRef *string, bitsHexLen int) (uint, error) {
	hexToDecode := a.getDynamicTrim(hexRef, bitsHexLen)

	switch rawHexLen := len(hexToDecode); {
	case rawHexLen <= utils.HexLength8Bits:
		targetValue, err := a.decodeBaseUint(&hexToDecode, utils.Bits8)
		return uint(*targetValue), err
	case rawHexLen <= utils.HexLength16Bits:
		targetValue, err := a.decodeBaseUint(&hexToDecode, utils.Bits16)
		return uint(*targetValue), err
	case rawHexLen <= utils.HexLength32Bits:
		targetValue, err := a.decodeBaseUint(&hexToDecode, utils.Bits32)
		return uint(*targetValue), err
	case rawHexLen <= utils.HexLength64Bits:
		targetValue, err := a.decodeBaseUint(&hexToDecode, utils.Bits64)
		return uint(*targetValue), err
	default:
		return 0, fmt.Errorf("invalid hex string %s to dedode to uint", hexToDecode)
	}
}

func (a *vmOutputData) decodeInt(hexRef *string, bitsHexLen int) (int, error) {
	hexToDecode := a.getDynamicTrim(hexRef, bitsHexLen)

	switch rawHexLen := len(hexToDecode); {
	case rawHexLen <= utils.HexLength8Bits:
		targetValue, err := a.decodeBaseUint(&hexToDecode, utils.Bits8)
		return int(int8(*targetValue)), err
	case rawHexLen <= utils.HexLength16Bits:
		targetValue, err := a.decodeBaseUint(&hexToDecode, utils.Bits16)
		return int(int16(*targetValue)), err
	case rawHexLen <= utils.HexLength32Bits:
		targetValue, err := a.decodeBaseUint(&hexToDecode, utils.Bits32)
		return int(int32(*targetValue)), err
	case rawHexLen <= utils.HexLength64Bits:
		targetValue, err := a.decodeBaseUint(&hexToDecode, utils.Bits64)
		return int(int64(*targetValue)), err
	default:
		return 0, fmt.Errorf("invalid hex string %s to decode to int", hexToDecode)
	}
}

func (a *vmOutputData) decodeBigUint(hexRef *string, trim int) (*big.Int, error) {
	hexToDecode := a.handleTrim(hexRef, trim*2)

	targetValue, err := a.decodeStringBigNumber(&hexToDecode)
	// if that function suceeds, then it was a string representing a decimal big number
	if err == nil {
		return targetValue, nil
	}

	targetValue, ok := new(big.Int).SetString(hexToDecode, utils.BaseHex)
	if !ok {
		return nil, fmt.Errorf("invalid hex string %s to decode to uint64", hexToDecode)
	}

	return targetValue, nil
}

func (a *vmOutputData) decodeBigInt(hexRef *string, trim int) (*big.Int, error) {
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

	targetValue, err := a.decodeInt(&hexToDecode, len(hexToDecode))
	if err != nil {
		return nil, fmt.Errorf("invelid hex string %s to decode to big int", hexToDecode)
	}

	return big.NewInt(int64(targetValue)), nil
}

func (a *vmOutputData) decodeStringBigNumber(hexRef *string) (*big.Int, error) {
	targetString, err := a.decodeString(hexRef, 0)
	if err != nil {
		return nil, err
	}

	targetValue, ok := new(big.Int).SetString(targetString, utils.BaseDecimal)
	if !ok {
		return nil, fmt.Errorf("invalid hex string %s to decode to BigInt", (*hexRef))
	}

	return targetValue, nil
}

func (a *vmOutputData) handleBigIntTill128(hexRef *string) (*big.Int, error) {
	rawValue, ok := new(big.Int).SetString(*hexRef, utils.BaseHex)
	if !ok {
		return nil, fmt.Errorf("invalid hex string to decode to BigInt: %s", *hexRef)
	}

	valueBits := len(*hexRef) * utils.BitsByHexDigit

	two := big.NewInt(2)
	twoToTheNth := new(big.Int).Exp(two, big.NewInt(int64(valueBits-1)), nil) // 2^valueBits
	one := big.NewInt(1)

	MaxIntNBits := new(big.Int).Sub(twoToTheNth, one)

	if rawValue.Cmp(new(big.Int).SetUint64(math.MaxUint64)) == 1 &&
		rawValue.Cmp(MaxIntNBits) == -1 {
		return rawValue, nil
	}

	if rawValue.Cmp(MaxIntNBits) == 1 {
		decodedValue := rawValue.Sub(rawValue, new(big.Int).Lsh(one, uint(valueBits)))
		return decodedValue, nil
	}

	return nil, fmt.Errorf("%v range is lower range than 128 bits", rawValue)
}

func (a *vmOutputData) decodeBoolean(hexRef *string) bool {
	const BooleanLength int = 2

	hexToDecode := a.getDynamicTrim(hexRef, BooleanLength)
	return hexToDecode == "01"
}

func (a *vmOutputData) decodeBigFloat(hexRef *string, trim int) (interface{}, error) {
	hexToDecode := a.handleTrim(hexRef, trim*2)

	hexBytes, err := hex.DecodeString(hexToDecode)
	if err != nil {
		return nil, fmt.Errorf("error decoding hex string to bytes: %w", err)
	}

	decodedValue := new(big.Float)

	if err := decodedValue.GobDecode(hexBytes); err != nil {
		return nil, fmt.Errorf("error decoding big float: %w", err)
	}

	return decodedValue, nil
}

func (a *vmOutputData) decodeOption(hexRef *string, valueType string) (interface{}, error) {
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

	return a.doDecode(hexRef, valueType, trim)
}

func (a *vmOutputData) getFixedTrim(hexRef *string) (int, error) {
	trimStringHex := (*hexRef)[:utils.LengthHexSizer]

	trimRef, err := a.decodeBaseUint(&trimStringHex, utils.LengthHexSizer*utils.BitsByHexDigit)

	*hexRef = (*hexRef)[utils.LengthHexSizer:]

	return int(*trimRef), err
}

func (a *vmOutputData) decodeList(hexRef *string, valueType string) (interface{}, error) {
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

func (a *vmOutputData) handleList(hexRef *string, valueType string) (interface{}, error) {
	wrapperType, coreType := utils.SplitTypes(valueType)

	if wrapperType == utils.List {
		listTrim, err := a.getFixedTrim(hexRef)
		if err != nil {
			return nil, fmt.Errorf("error getting the list trim: %w", err)
		}

		return a.decodeNestedList(hexRef, coreType, listTrim)
	}

	var valueTrim int

	if utils.IsDynamicLengthType(valueType) {
		calculatedTrim, err := a.getFixedTrim(hexRef)
		if err != nil {
			return nil, fmt.Errorf("error getting the list item trim: %w", err)
		}

		valueTrim = calculatedTrim
	}

	return a.decodeSingleValue(hexRef, coreType, valueTrim)
}

func (a *vmOutputData) decodeNestedList(
	hexRef *string,
	valueType string,
	limit int,
) (interface{}, error) {
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

func (a *vmOutputData) decodeTuple(hexRef *string, valueType string) ([]interface{}, error) {
	var result []interface{}

	types := utils.SplitTupleTypes(valueType)

	for _, t := range types {
		if strings.HasPrefix(t, utils.List) {
			decodedList, err := a.handleList(hexRef, t)
			if err != nil {
				return nil, err
			}

			result = append(result, decodedList)
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

		decodedValue, err := a.doDecode(hexRef, t, valueTrim)
		if err != nil {
			return nil, fmt.Errorf("error decoding tuple value %w", err)
		}

		result = append(result, decodedValue)
	}

	return result, nil
}

func (a *vmOutputData) decodeVariadic(hexRef *string, valueType string) (interface{}, error) {
	decodedValue, err := a.doDecode(hexRef, valueType, 0)
	if err != nil {
		return nil, fmt.Errorf("error trying to decode variadic value: %w", err)
	}

	return decodedValue, nil
}

func (a *vmOutputData) decodeStruct(
	hexRef *string,
	valueType string,
) (map[string]interface{}, error) {
	typeDef, exists := a.Types[valueType]
	if !exists {
		return nil, fmt.Errorf("type %s not found in provided abi", valueType)
	}

	result := make(map[string]interface{})

	for _, field := range typeDef.Fields {
		if strings.HasPrefix(field.Type, utils.List) {
			decodedList, err := a.handleList(hexRef, field.Type)
			if err != nil {
				return nil, fmt.Errorf(
					"error %w decoding list value of key %s of custom type %s",
					err,
					field.Type,
					valueType,
				)
			}

			result[field.Name] = decodedList
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

		decodedValue, err := a.doDecode(hexRef, field.Type, trim)
		if err != nil {
			return nil, fmt.Errorf(
				"error %w decoding value of key %s of custom type %s",
				err,
				field.Type,
				valueType,
			)
		}

		result[field.Name] = decodedValue
	}

	return result, nil
}
