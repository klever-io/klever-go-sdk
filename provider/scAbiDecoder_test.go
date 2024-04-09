package provider_test

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/klever-io/klever-go-sdk/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Decode_Single_Value(t *testing.T) {
	jsonAbi, err := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	if err != nil {
		fmt.Println("error opening abi", err)
		return
	}
	defer jsonAbi.Close()

	abiHandler := provider.NewSCAbiHandler()

	if err := abiHandler.LoadAbi(jsonAbi); err != nil {
		fmt.Println("error loading abi reader", err)
		return
	}

	testCases := []struct {
		name     string
		endpoint string
		hex      string
		expected any
	}{
		{
			name:     "Managed_Buffer",
			endpoint: "managed_buffer",
			hex:      "74657374696e67206f757470757473207479706573",
			expected: "testing outputs types",
		},
		{
			name:     "Boolean_false",
			endpoint: "bool_false",
			hex:      "",
			expected: false,
		},
		{
			name:     "Boolean_true",
			endpoint: "bool_true",
			hex:      "01",
			expected: true,
		},
		{
			name:     "Usize",
			endpoint: "usize_number",
			hex:      "fdc20cbf",
			expected: uint32(4257352895),
		},
		{
			name:     "Isize",
			endpoint: "isize_number",
			hex:      "40cf4061",
			expected: int32(1087324257),
		},
		{
			name:     "Negative_isize",
			endpoint: "isize_minus_number",
			hex:      "bf30bf9f",
			expected: int32(-1087324257),
		},
		{
			name:     "Token_identifier",
			endpoint: "token_identifier",
			hex:      "4b4c56",
			expected: "KLV",
		},
		{
			name:     "Address",
			endpoint: "owner_address",
			hex:      "667fd274481cf5b07418b2fdc5d8baa6ae717239357f338cde99c2f612a96a9e",
			expected: "klv1velayazgrn6mqaqckt7utk9656h8zu3ex4ln8rx7n8p0vy4fd20qmwh4p5",
		},
		{
			name:     "BigInt_from_small_positive_string",
			endpoint: "big_s_10",
			hex:      "3130",
			expected: big.NewInt(10),
		},
		{
			name:     "BigInt_from_small_negative_string",
			endpoint: "big_minus_s_10",
			hex:      "2d3130",
			expected: big.NewInt(-10),
		},
		{
			name:     "BigInt_from_negative_int8",
			endpoint: "big_minus_i8",
			hex:      "ae",
			expected: big.NewInt(-82),
		},
		{
			name:     "BigInt_from_positive_int8",
			endpoint: "big_i8",
			hex:      "52",
			expected: big.NewInt(82),
		},
		{
			name:     "BigInt_from_negative_int16",
			endpoint: "big_minus_i16",
			hex:      "dede",
			expected: big.NewInt(-8482),
		},
		{
			name:     "BigInt_from_positive_int16",
			endpoint: "big_i16",
			hex:      "2122",
			expected: big.NewInt(8482),
		},
		{
			name:     "BigInt_from_negative_int32",
			endpoint: "big_minus_i32",
			hex:      "e69fa284",
			expected: big.NewInt(-425745788),
		},
		{
			name:     "BigInt_from_positive_int32",
			endpoint: "big_i32",
			hex:      "19605d7c",
			expected: big.NewInt(425745788),
		},
		{
			name:     "BigInt_from_negative_int64",
			endpoint: "big_minus_i64",
			hex:      "c91131a14fc23dac",
			expected: big.NewInt(-3958328028584329812),
		},
		{
			name:     "BigInt_from_negative_int64",
			endpoint: "big_i64",
			hex:      "36eece5eb03dc254",
			expected: big.NewInt(3958328028584329812),
		},
		{
			name:     "BigInt_from_uint128",
			endpoint: "big_u_number",
			hex:      "39bf6e49095ff7dca078957ceb928e",
			expected: func() *big.Int {
				bigInt128, _ := new(big.Int).SetString("299843598872398459348567275690758798", 10)
				return bigInt128
			}(),
		},
		{
			name:     "BigInt_negative_from_uint128",
			endpoint: "big_minus_u_number",
			hex:      "c64091b6f6a008235f876a83146d72",
			expected: func() *big.Int {
				bigInt128, _ := new(big.Int).SetString("-299843598872398459348567275690758798", 10)
				return bigInt128
			}(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := abiHandler.Decode(testCase.endpoint, testCase.hex)

			fmt.Println(result)

			require.Nil(t, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}
