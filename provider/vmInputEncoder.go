package provider

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/klever-io/klever-go-sdk/core/address"
	"github.com/klever-io/klever-go-sdk/provider/utils"
)

func EncodeInput(args []string) (string, error) {
	encodedArgs := ""
	for _, arg := range args {
		encoded, err := doEncode(arg)
		if err != nil {
			return "", fmt.Errorf("invalid item to encode `%s`", arg)
		}

		encodedArgs += "@" + encoded
	}

	return encodedArgs, nil
}

func doEncode(arg string) (string, error) {
	kv := strings.SplitN(arg, ":", 2)

	isOption := false
	// check if it is an option argument
	if strings.HasPrefix(kv[0], "option") {
		isOption = true
		kv[0] = kv[0][6:]
	}

	encoded, err := encodeSingleValue(kv[0], kv[1], isOption)
	if err != nil {
		return "", err
	}

	if isOption {
		encoded = "01" + encoded
	}

	return encoded, nil
}

func encodeSingleValue(t, v string, isNested bool) (string, error) {
	switch t {
	case
		utils.Int8, strings.ToUpper(utils.Int8), // int8
		utils.Int16, strings.ToUpper(utils.Int16), // int16
		utils.Int32, strings.ToUpper(utils.Int32), utils.Isize, // int32
		utils.Int64, strings.ToUpper(utils.Int64): // int864

		return encodeInt(v, t, isNested)
	case
		utils.Uint8, strings.ToUpper(utils.Uint8), // uint8
		utils.Uint16, strings.ToUpper(utils.Uint64), // uint16
		utils.Uint32, strings.ToUpper(utils.Uint64), utils.Usize, // uint32
		utils.Uint64, strings.ToUpper(utils.Uint64): // uint864

		return encodeUint(v, t, isNested)
	case utils.BigInt, strings.ToLower(utils.BigInt),
		utils.BigUint, strings.ToLower(utils.BigUint),
		"bi", "n", "BI", "N": // BigInt and BigUint

		return encodeBigInt(v, isNested)
	case utils.Address, strings.ToLower(utils.Address), "a", "A":
		return encodeAddress(v)
	case utils.BigFloat, strings.ToLower(utils.BigFloat),
		"bf", "BF", "f", "F": // BigFloat:

		return encodeBigFloat(v, isNested)
	case
		utils.ManagedBuffer, strings.ToLower(utils.ManagedBuffer),
		utils.TokenIdentifier, strings.ToLower(utils.TokenIdentifier),
		utils.Bytes,
		utils.BoxedBytes, strings.ToLower(utils.BoxedBytes),
		utils.String, strings.ToLower(utils.String),
		utils.StrRef,
		utils.VecU8, strings.ToLower(utils.VecU8),
		utils.SliceU8: // All types of strings

		return encodeString(v, isNested), nil
	case "bool", "boolean", "b", "B":
		return encodeBoolean(v, isNested)
	case "empty", "0", "e", "E":
		return "", nil
	default:
		return "", fmt.Errorf("invalid encode type `%s`", t)
	}
}

func encodeInt(v, t string, isNested bool) (string, error) {
	if isNested {
		return encodeNestedInt(v, t)
	}

	return encodeTopLevelInt(v)
}

func encodeTopLevelInt(v string) (string, error) {
	rawInt, err := strconv.ParseInt(v, utils.BaseDecimal, utils.Bits64)
	if err != nil {
		return "", fmt.Errorf("invalid string `%s` to convert to signed integer", v)
	}
	switch {
	// typecast to uint to correspondent bit size to use 2's complement in negative values
	case rawInt >= math.MinInt8 && rawInt <= math.MaxInt8:
		return fmt.Sprintf("%02x", uint8(rawInt)), nil
	case rawInt >= math.MinInt16 && rawInt <= math.MaxInt16:
		return fmt.Sprintf("%04x", uint16(rawInt)), nil
	case rawInt >= math.MinInt32 && rawInt <= math.MaxInt32:
		return fmt.Sprintf("%08x", uint32(rawInt)), nil
	default:
		return fmt.Sprintf("%016x", uint64(rawInt)), nil
	}
}

func encodeNestedInt(v, t string) (string, error) {
	rawInt, err := strconv.ParseInt(v, utils.BaseDecimal, utils.Bits64)
	if err != nil {
		return "", fmt.Errorf("invalid string `%s` to convert to signed integer", v)
	}
	switch t {
	// typecast to uint to correspondent bit size to use 2's complement in negative values
	case utils.Int8:
		return fmt.Sprintf("%02x", uint8(rawInt)), nil
	case utils.Int16:
		return fmt.Sprintf("%04x", uint16(rawInt)), nil
	case utils.Int32:
		return fmt.Sprintf("%08x", uint32(rawInt)), nil
	default:
		return fmt.Sprintf("%016x", uint64(rawInt)), nil
	}
}

func encodeUint(v, t string, isNested bool) (string, error) {
	if isNested {
		return encodeNestedUint(v, t)
	}

	return encodeTopLevelUint(v)
}

func encodeTopLevelUint(v string) (string, error) {
	rawInt, err := strconv.ParseUint(v, utils.BaseDecimal, utils.Bits64)
	if err != nil {
		return "", fmt.Errorf("invalid string `%s` to convert to signed integer", v)
	}
	switch {
	case rawInt <= math.MaxUint8:
		return fmt.Sprintf("%02x", uint8(rawInt)), nil
	case rawInt <= math.MaxUint16:
		return fmt.Sprintf("%04x", uint16(rawInt)), nil
	case rawInt <= math.MaxUint32:
		return fmt.Sprintf("%08x", uint32(rawInt)), nil
	default:
		return fmt.Sprintf("%016x", rawInt), nil
	}
}

func encodeNestedUint(v, t string) (string, error) {
	rawInt, err := strconv.ParseUint(v, utils.BaseDecimal, utils.Bits64)
	if err != nil {
		return "", fmt.Errorf("invalid string `%s` to convert to signed integer", v)
	}
	switch t {
	case utils.Int8:
		return fmt.Sprintf("%02x", uint8(rawInt)), nil
	case utils.Int16:
		return fmt.Sprintf("%04x", uint16(rawInt)), nil
	case utils.Int32:
		return fmt.Sprintf("%08x", uint32(rawInt)), nil
	default:
		return fmt.Sprintf("%016x", rawInt), nil
	}
}

func encodeAddress(v string) (string, error) {
	address, err := address.NewAddress(v)
	if err != nil {
		return "", fmt.Errorf("invalid address `%s` to encode to hexadecimal format", v)
	}

	return address.Hex(), nil
}

func encodeString(v string, isNested bool) string {
	hexString := ""

	for i := 0; i < len(v); i++ {
		hexString += fmt.Sprintf("%02x", v[i])
	}

	if isNested {
		hexLen := fmt.Sprintf("%08x", len(hexString)/2)
		hexString = hexLen + hexString
	}

	return hexString
}

func encodeBoolean(v string, isNested bool) (string, error) {
	switch v {
	case "true":
		return "01", nil
	case "false":
		if isNested {
			return "00", nil
		}
		return "", nil
	default:
		return "", fmt.Errorf("invalid boolean to encode, must be `true` or `false`")
	}
}

func encodeBigInt(v string, isNested bool) (string, error) {
	bi, ok := new(big.Int).SetString(v, utils.BaseDecimal)
	if !ok {
		return "", fmt.Errorf("invalid string `%s` to convert to big integer", v)
	}

	if bi.BitLen() > utils.Bits128 {
		return fmt.Sprintf("%x", v), nil
	}

	if bi.BitLen() <= utils.Bits128 && bi.Sign() == -1 {
		bigIntTwosComplement(bi)
	}

	hexBi := fmt.Sprintf("%x", bi)

	if len(hexBi)%2 != 0 {
		if strings.HasPrefix(v, "-") {
			hexBi = "f" + hexBi
		} else {
			hexBi = "0" + hexBi
		}
	}

	if isNested {
		hexLen := fmt.Sprintf("%08x", len(hexBi)/2)
		hexBi = hexLen + hexBi
	}

	return hexBi, nil
}

func bigIntTwosComplement(b *big.Int) {
	bitsLen := b.BitLen()
	// Here, even though `bitsLen` isn't a multiple of 4, by
	// adding 4 to it we are ensuring that bitMask created later is large enough
	// to cover all bits of the original value `b` in the Xor operation.
	bitsLen += utils.BitsByHexDigit

	maskHexLen := bitsLen / utils.BitsByHexDigit
	maskStringHex := strings.Repeat("f", maskHexLen)
	bitMask, _ := new(big.Int).SetString(maskStringHex, utils.BaseHex)

	b.Abs(b)

	b.Xor(b, bitMask)

	b.Add(b, big.NewInt(1))
}

func encodeBigFloat(v string, isNested bool) (string, error) {
	bf, ok := new(big.Float).SetString(v)
	if !ok {
		return "", fmt.Errorf(
			"invalid string `%s` representing a big float to encode to hexadecimal",
			v,
		)
	}

	bf.SetPrec(utils.BigFloatVMPrecision)
	bf.SetMode(big.RoundingMode(big.Exact))

	bfBytes, err := bf.GobEncode()
	if err != nil {
		return "", fmt.Errorf(
			"invalid string `%s` representing a big float to encode to hexadecimal",
			v,
		)
	}

	hexBf := hex.EncodeToString(bfBytes)

	if isNested {
		hexLen := fmt.Sprintf("%08x", len(hexBf)/2)
		hexBf = hexLen + hexBf
	}

	return hexBf, nil
}
