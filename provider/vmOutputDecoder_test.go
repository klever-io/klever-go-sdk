package provider_test

import (
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/klever-io/klever-go-sdk/provider"
)

const BaseDecimal = 10

func Test_Parse_Single_Value(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewVMOutputHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

	testCases := []struct {
		name     string
		endpoint string
		hex      []string
		expected any
	}{
		{
			name:     "Managed_Buffer",
			endpoint: "managed_buffer",
			hex:      []string{"74657374696e67206f757470757473207479706573"},
			expected: "testing outputs types",
		},
		{
			name:     "Boolean_false",
			endpoint: "bool_false",
			hex:      []string{""},
			expected: false,
		},
		{
			name:     "Boolean_true",
			endpoint: "bool_true",
			hex:      []string{"01"},
			expected: true,
		},
		{
			name:     "Usize",
			endpoint: "usize_number",
			hex:      []string{"fdc20cbf"},
			expected: uint32(4257352895),
		},
		{
			name:     "Isize",
			endpoint: "isize_number",
			hex:      []string{"40cf4061"},
			expected: int32(1087324257),
		},
		{
			name:     "Negative_isize",
			endpoint: "isize_minus_number",
			hex:      []string{"bf30bf9f"},
			expected: int32(-1087324257),
		},
		{
			name:     "Token_identifier",
			endpoint: "token_identifier",
			hex:      []string{"4b4c56"},
			expected: "KLV",
		},
		{
			name:     "Address",
			endpoint: "owner_address",
			hex:      []string{"667fd274481cf5b07418b2fdc5d8baa6ae717239357f338cde99c2f612a96a9e"},
			expected: "klv1velayazgrn6mqaqckt7utk9656h8zu3ex4ln8rx7n8p0vy4fd20qmwh4p5",
		},
		{
			name:     "BigInt_from_negative_int8",
			endpoint: "big_minus_i8",
			hex:      []string{"ae"},
			expected: big.NewInt(-82),
		},
		{
			name:     "BigInt_from_positive_int8",
			endpoint: "big_i8",
			hex:      []string{"52"},
			expected: big.NewInt(82),
		},
		{
			name:     "BigInt_from_negative_int16",
			endpoint: "big_minus_i16",
			hex:      []string{"dede"},
			expected: big.NewInt(-8482),
		},
		{
			name:     "BigInt_from_positive_int16",
			endpoint: "big_i16",
			hex:      []string{"2122"},
			expected: big.NewInt(8482),
		},
		{
			name:     "BigInt_from_negative_int32",
			endpoint: "big_minus_i32",
			hex:      []string{"e69fa284"},
			expected: big.NewInt(-425745788),
		},
		{
			name:     "BigInt_from_positive_int32",
			endpoint: "big_i32",
			hex:      []string{"19605d7c"},
			expected: big.NewInt(425745788),
		},
		{
			name:     "BigInt_from_negative_int64",
			endpoint: "big_minus_i64",
			hex:      []string{"c91131a14fc23dac"},
			expected: big.NewInt(-3958328028584329812),
		},
		{
			name:     "BigInt_from_negative_int64",
			endpoint: "big_i64",
			hex:      []string{"36eece5eb03dc254"},
			expected: big.NewInt(3958328028584329812),
		},
		{
			name:     "BigInt_from_uint128",
			endpoint: "big_u_number",
			hex:      []string{"39bf6e49095ff7dca078957ceb928e"},
			expected: func() *big.Int {
				bigInt128, _ := new(
					big.Int,
				).SetString("299843598872398459348567275690758798", BaseDecimal)
				return bigInt128
			}(),
		},
		{
			name:     "BigInt_negative_from_uint128",
			endpoint: "big_minus_u_number",
			hex:      []string{"c64091b6f6a008235f876a83146d72"},
			expected: func() *big.Int {
				bigInt128, _ := new(
					big.Int,
				).SetString("-299843598872398459348567275690758798", BaseDecimal)
				return bigInt128
			}(),
		},
		{
			name:     "BigUint_from_u8",
			endpoint: "big_u8",
			hex:      []string{"52"},
			expected: big.NewInt(82),
		},
		{
			name:     "BigUint_from_u16",
			endpoint: "big_u16",
			hex:      []string{"2122"},
			expected: big.NewInt(8482),
		},
		{
			name:     "BigUint_from_u32",
			endpoint: "big_u32",
			hex:      []string{"19605d7c"},
			expected: big.NewInt(425745788),
		},
		{
			name:     "BigUint_from_u64",
			endpoint: "big_u64",
			hex:      []string{"36eece5eb03dc254"},
			expected: big.NewInt(3958328028584329812),
		},
		{
			name:     "BigUint_from_u128",
			endpoint: "big_u128",
			hex:      []string{"1dc7766516260b32b52ff11612d5710e"},
			expected: func() *big.Int {
				bigUintFromU128, _ := new(
					big.Int,
				).SetString("39583280285843298128735477835272384782", BaseDecimal)
				return bigUintFromU128
			}(),
		},
		{
			name:     "Int8_positive",
			endpoint: "number_i8",
			hex:      []string{"52"},
			expected: int8(82),
		},
		{
			name:     "Int8_negative",
			endpoint: "number_minus_i8",
			hex:      []string{"ae"},
			expected: int8(-82),
		},
		{
			name:     "Int16_positive",
			endpoint: "number_i16",
			hex:      []string{"2122"},
			expected: int16(8482),
		},
		{
			name:     "Int16_negative",
			endpoint: "number_minus_i16",
			hex:      []string{"dede"},
			expected: int16(-8482),
		},
		{
			name:     "Int32_positive",
			endpoint: "number_i32",
			hex:      []string{"19605d7c"},
			expected: int32(425745788),
		},
		{
			name:     "Int32_negative",
			endpoint: "number_minus_i32",
			hex:      []string{"e69fa284"},
			expected: int32(-425745788),
		},
		{
			name:     "Int64_positive",
			endpoint: "number_i64",
			hex:      []string{"36eece5eb03dc254"},
			expected: int64(3958328028584329812),
		},
		{
			name:     "Int64_negative",
			endpoint: "number_minus_i64",
			hex:      []string{"c91131a14fc23dac"},
			expected: int64(-3958328028584329812),
		},
		{
			name:     "BigFloat_from_fraction",
			endpoint: "test_from_fraction",
			hex:      []string{"010a0000003500000000c000000000000000"},
			expected: big.NewFloat(0.75),
		},
		{
			name:     "BigFloat_from_fraction_negative",
			endpoint: "test_float_neg",
			hex:      []string{"010b00000035000000008000000000000000"},
			expected: big.NewFloat(-0.5),
		},
		{
			name:     "BigFloat_from_parts",
			endpoint: "test_float_parts",
			hex:      []string{"010a000000350000000f93eac189374bc800"},
			expected: big.NewFloat(18933.378),
		},
		{
			name:     "BigFloat_from_scientific_notation",
			endpoint: "test_float_sci",
			hex:      []string{"010a00000035fffffffcf5c28f5c28f5c000"},
			expected: big.NewFloat(0.06),
		},
		{
			name:     "BigFloat_from_big_int",
			endpoint: "test_float_bi",
			hex:      []string{"010b0000003500000053d7da3fa6c03de800"},
			expected: func() *big.Float {
				bf, _ := new(big.Float).SetString("-8154678164717479819989764")
				bf.SetPrec(53)
				bf.SetMode(big.RoundingMode(big.Exact))
				return bf
			}(),
		},
		{
			name:     "BigFloat_from_big_uint",
			endpoint: "test_float_bu",
			hex:      []string{"010a000000350000005bb35b36c567c96000"},
			expected: func() *big.Float {
				bf, _ := new(big.Float).SetString("1734627739277592794878918977")
				bf.SetPrec(53)
				bf.SetMode(big.RoundingMode(big.Exact))
				return bf
			}(),
		},
		{
			name:     "BigFloat_from_generics",
			endpoint: "test_float_from_generics",
			hex:      []string{"010a00000035000000028000000000000000"},
			expected: big.NewFloat(2),
		},
		{
			name:     "BigFloat_from_generics",
			endpoint: "test_float_from_generics",
			hex:      []string{"010a00000035000000028000000000000000"},
			expected: big.NewFloat(2),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := abiHandler.DecodeHex(testCase.endpoint, testCase.hex)

			require.Nil(t, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func Test_Parse_Single_Enum_Value(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewVMOutputHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

	testCases := []struct {
		name     string
		endpoint string
		hex      []string
		expected string
	}{
		{
			name:     "Zero_discriminant",
			endpoint: "enum_test",
			hex:      []string{""},
			expected: "Neutral",
		},
		{
			name:     "One_discriminant",
			endpoint: "enum_test",
			hex:      []string{"01"},
			expected: "Down",
		},
		{
			name:     "Two_discriminant",
			endpoint: "enum_test",
			hex:      []string{"02"},
			expected: "Up",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := abiHandler.DecodeHex(testCase.endpoint, testCase.hex)

			require.Nil(t, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func Test_Parse_List(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewVMOutputHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error loading abi", errLoad)

	testCases := []struct {
		name     string
		endpoint string
		hex      []string
		expected any
	}{
		{
			name:     "Token_identifiers",
			endpoint: "list_token_identifier",
			hex: []string{
				"000000034b4c56000000034b4649000000084b49442d38473941000000084458422d483838470000000a43484950532d4e383941",
			},
			expected: []string{"KLV", "KFI", "KID-8G9A", "DXB-H88G", "CHIPS-N89A"},
		},
		{
			name:     "int32",
			endpoint: "list_int32",
			hex:      []string{"000000080000005700000065fffffffb"},
			expected: []int32{int32(8), int32(87), int32(101), int32(-5)},
		},
		{
			name:     "int64",
			endpoint: "list_int64",
			hex: []string{
				"000000000000000200000000000000570e09174747d3c452fffffffffffffffbe51dd01a95946e83",
			},
			expected: []int64{
				int64(2),
				int64(87),
				int64(1011365186236564562),
				int64(-5),
				int64(-1937163452102185341),
			},
		},
		{
			name:     "uint16",
			endpoint: "list_u16",
			hex:      []string{"3fd600a7000203f3"},
			expected: []uint16{uint16(16342), uint16(167), uint16(2), uint16(1011)},
		},
		{
			name:     "uint32",
			endpoint: "list_u32",
			hex:      []string{"00000003000f880ccc61aa59003ab9b8"},
			expected: []uint32{uint32(3), uint32(1017868), uint32(3428952665), uint32(3848632)},
		},
		{
			name:     "uint64",
			endpoint: "list_u64",
			hex:      []string{"00000000053532bd0002ac5552d00e95000000000000000255868c6974ec6a9b"},
			expected: []uint64{
				uint64(87372477),
				uint64(752432414985877),
				uint64(2),
				uint64(6162767524664208027),
			},
		},
		{
			name:     "BigInt",
			endpoint: "list_bign",
			hex: []string{
				"000000050577f69535000000010900000004fff36f50000000030826dd0000000f39bf6e49095ff7dca078957ceb928e0000000fc64091b6f6a008235f876a83146d72",
			},
			expected: func() []*big.Int {
				bigInt128pos, _ := new(
					big.Int,
				).SetString("299843598872398459348567275690758798", BaseDecimal)
				bigInt128neg, _ := new(
					big.Int,
				).SetString("-299843598872398459348567275690758798", BaseDecimal)

				return []*big.Int{
					big.NewInt(23487485237),
					big.NewInt(9),
					big.NewInt(-823472),
					big.NewInt(534237),
					bigInt128pos,
					bigInt128neg,
				}
			}(),
		},
		{
			name:     "BigUint",
			endpoint: "list_bigun",
			hex: []string{
				"00000001ea00000002266a000000043a9e8554000000087864b47dcf08ef8c0000001d0311821D7EB083EE2D255B7DFA8F665639FED96C2998FDB113657562A90000000f0864a6c0c92180ec36795616644d36",
			},
			expected: func() []*big.Int {
				bigUintString, _ := new(
					big.Int,
				).SetString("82723675254376576786347243658263586823568268912736457678698677887657", BaseDecimal)
				bigUint128, _ := new(
					big.Int,
				).SetString("43579827367895347689574268789869878", BaseDecimal)

				return []*big.Int{
					big.NewInt(234),
					big.NewInt(9834),
					big.NewInt(983467348),
					big.NewInt(8675257234659798924),
					bigUintString,
					bigUint128,
				}
			}(),
		},
		{
			name:     "Address",
			endpoint: "list_address",
			hex: []string{
				"667fd274481cf5b07418b2fdc5d8baa6ae717239357f338cde99c2f612a96a9e667fd274481cf5b07418b2fdc5d8baa6ae717239357f338cde99c2f612a96a9e",
			},
			expected: []string{
				"klv1velayazgrn6mqaqckt7utk9656h8zu3ex4ln8rx7n8p0vy4fd20qmwh4p5",
				"klv1velayazgrn6mqaqckt7utk9656h8zu3ex4ln8rx7n8p0vy4fd20qmwh4p5",
			},
		},
		{
			name:     "Boolean",
			endpoint: "list_bool",
			hex:      []string{"01000001"},
			expected: []bool{true, false, false, true},
		},
		{
			name:     "Nested_list_boolean",
			endpoint: "list_list_bool",
			hex:      []string{"00000004010000010000000401000100"},
			expected: [][]interface{}{
				{true, false, false, true},
				{true, false, true, false},
			},
		},
		{
			name:     "Nested_list_int32",
			endpoint: "list_list_i32",
			hex:      []string{"000000036fab02760000001cffff9ca400000003fffffffe0001e308f1b3cfdc"},
			expected: [][]interface{}{
				{int32(1873478262), int32(28), int32(-25436)},
				{int32(-2), int32(123656), int32(-239874084)},
			},
		},
		{
			name:     "Nested_list_int64",
			endpoint: "list_list_i64",
			hex: []string{
				"0000000319ffee93a36dc12a000000000000001cd7e571502441e18400000003fffffffffffffffe000000000001e308fffffffff1b3cfdc",
			},
			expected: [][]interface{}{
				{int64(1873478287878897962), int64(28), int64(-2889778996868685436)},
				{int64(-2), int64(123656), int64(-239874084)},
			},
		},
		{
			name:     "Nested_list_tokens_identifiers",
			endpoint: "list_list_tokens",
			hex: []string{
				"00000003000000034b4c56000000034b4649000000084b49442d3847394100000003000000084458422d483838470000000a43484950532d4e383941000000084646542d32424836",
			},
			expected: [][]interface{}{
				{"KLV", "KFI", "KID-8G9A"},
				{"DXB-H88G", "CHIPS-N89A", "FFT-2BH6"},
			},
		},
		{
			name:     "Two_levels_nested_list_int64",
			endpoint: "list_list_list_i64",
			hex: []string{
				"0000000200000003000000006fab0276000000000000001cffffffffffff9ca400000003fffffffffffffffe000000000001e308fffffffff1b3cfdc0000000200000003fffffffffe4932e1000007f24bf1555700000000000000800000000300000010d9e95b690000000000000001ffffffffffffff9c",
			},
			expected: [][]interface{}{
				{
					[]interface{}{int64(1873478262), int64(28), int64(-25436)},
					[]interface{}{int64(-2), int64(123656), int64(-239874084)},
				},
				{
					[]interface{}{int64(-28757279), int64(8737237587287), int64(128)},
					[]interface{}{int64(72375425897), int64(1), int64(-100)},
				},
			},
		},
		{
			name:     "Two_levels_nested_list_token-identifiers",
			endpoint: "list_list_list_token",
			hex: []string{
				"00000002000000030000000a43484950532d4e383941000000034b4c56000000034b464900000003000000085446542d3738364a00000008534a412d4c4b394800000008514b552d3748483100000002000000030000000a43484950532d4e383941000000034b4c56000000034b464900000003000000085446542d3738364a00000008534a412d4c4b394800000008514b552d37484831",
			},
			expected: [][]interface{}{
				{
					[]interface{}{string("CHIPS-N89A"), string("KLV"), string("KFI")},
					[]interface{}{string("TFT-786J"), string("SJA-LK9H"), string("QKU-7HH1")},
				},
				{
					[]interface{}{string("CHIPS-N89A"), string("KLV"), string("KFI")},
					[]interface{}{string("TFT-786J"), string("SJA-LK9H"), string("QKU-7HH1")},
				},
			},
		},
		{
			name:     "Two_levels_nested_list_big_int",
			endpoint: "list_list_list_big_int",
			hex: []string{
				"00000002000000030000000387efdb00000002349000000008c91131a14fc23dac000000030000000f39bf6e49095ff7dca078957ceb928e0000000fc64091b6f6a008235f876a83146d72000000050103cf74410000000200000003000000111d05b3eb926950d5852179a29534618bf100000011e2fa4c146d96af2a7ade865d6acb9e740f000000012a00000003000000037810250000000154000000064d9f58c4219f",
			},
			expected: func() [][]interface{} {
				big1, _ := new(big.Int).SetString("-7868453", BaseDecimal)
				big2, _ := new(big.Int).SetString("13456", BaseDecimal)
				big3, _ := new(big.Int).SetString("-3958328028584329812", BaseDecimal)
				big4, _ := new(
					big.Int,
				).SetString("299843598872398459348567275690758798", BaseDecimal)
				big5, _ := new(
					big.Int,
				).SetString("-299843598872398459348567275690758798", BaseDecimal)
				big6, _ := new(big.Int).SetString("4358894657", BaseDecimal)
				big7, _ := new(
					big.Int,
				).SetString("9875768979789979987589732879352104848369", BaseDecimal)
				big8, _ := new(
					big.Int,
				).SetString("-9875768979789979987589732879352104848369", BaseDecimal)
				big9, _ := new(big.Int).SetString("42", BaseDecimal)
				big10, _ := new(big.Int).SetString("7868453", BaseDecimal)
				big11, _ := new(big.Int).SetString("84", BaseDecimal)
				big12, _ := new(big.Int).SetString("85346784387487", BaseDecimal)

				return [][]interface{}{
					{
						[]interface{}{big1, big2, big3},
						[]interface{}{big4, big5, big6},
					},
					{
						[]interface{}{big7, big8, big9},
						[]interface{}{big10, big11, big12},
					},
				}
			}(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := abiHandler.DecodeHex(testCase.endpoint, testCase.hex)

			require.Nil(t, err)
			assert.ElementsMatch(t, testCase.expected, result)
		})
	}
}

func Test_Parse_Option(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewVMOutputHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

	testCases := []struct {
		name     string
		endpoint string
		hex      []string
		expected any
	}{
		{
			name:     "null",
			endpoint: "option_bytes_null",
			hex:      []string{""},
			expected: nil,
		},
		{
			name:     "i8",
			endpoint: "option_i8",
			hex:      []string{"0152"},
			expected: int8(82),
		},
		{
			name:     "i16",
			endpoint: "option_i16",
			hex:      []string{"012122"},
			expected: int16(8482),
		},
		{
			name:     "i32",
			endpoint: "option_i32",
			hex:      []string{"0119605d7c"},
			expected: int32(425745788),
		},
		{
			name:     "i64",
			endpoint: "option_i64",
			hex:      []string{"0136eece5eb03dc254"},
			expected: int64(3958328028584329812),
		},
		{
			name:     "u8",
			endpoint: "option_u8",
			hex:      []string{"0152"},
			expected: uint8(82),
		},
		{
			name:     "u16",
			endpoint: "option_u16",
			hex:      []string{"012122"},
			expected: uint16(8482),
		},
		{
			name:     "u32",
			endpoint: "option_u32",
			hex:      []string{"0119605d7c"},
			expected: uint32(425745788),
		},
		{
			name:     "u64",
			endpoint: "option_u64",
			hex:      []string{"0136eece5eb03dc254"},
			expected: uint64(3958328028584329812),
		},
		{
			name:     "big_int",
			endpoint: "option_bigint",
			hex:      []string{"01000000050577f69535"},
			expected: big.NewInt(23487485237),
		},
		{
			name:     "big_int_from_buffer",
			endpoint: "option_bigint_from_buffer",
			hex:      []string{"0100000008209889945685c1b5"},
			expected: func() *big.Int {
				big, _ := new(big.Int).SetString("2348778475873485237", BaseDecimal)
				return big
			}(),
		},
		{
			name:     "bigint_from_biguint_plus",
			endpoint: "option_bigint_from_biguint_plus",
			hex:      []string{"010000000f39bf6e49095ff7dca078957ceb928e"},
			expected: func() *big.Int {
				big, _ := new(
					big.Int,
				).SetString("299843598872398459348567275690758798", BaseDecimal)
				return big
			}(),
		},
		{
			name:     "bigint_from_biguint_minus",
			endpoint: "option_bigint_from_biguint_minus",
			hex:      []string{"010000000fc64091b6f6a008235f876a83146d72"},
			expected: func() *big.Int {
				big, _ := new(
					big.Int,
				).SetString("-299843598872398459348567275690758798", BaseDecimal)
				return big
			}(),
		},
		{
			name:     "biguint",
			endpoint: "option_biguint",
			hex:      []string{"010000000537952bd072"},
			expected: big.NewInt(238725877874),
		},
		{
			name:     "biguint_from_buffer",
			endpoint: "option_biguint_from_buffer",
			hex:      []string{"01000000090cf14c43036fdb7472"},
			expected: func() *big.Int {
				big, _ := new(big.Int).SetString("238748274827925877874", BaseDecimal)
				return big
			}(),
		},
		{
			name:     "address",
			endpoint: "option_address",
			hex: []string{
				"01667fd274481cf5b07418b2fdc5d8baa6ae717239357f338cde99c2f612a96a9e",
			},
			expected: "klv1velayazgrn6mqaqckt7utk9656h8zu3ex4ln8rx7n8p0vy4fd20qmwh4p5",
		},
		{
			name:     "bytes",
			endpoint: "option_managed_buffer",
			hex:      []string{"010000001574657374696e67206f757470757473207479706573"},
			expected: "testing outputs types",
		},
		{
			name:     "bigfloat",
			endpoint: "test_float_option",
			hex:      []string{"0100000012010a00000035000000028000000000000000"},
			expected: big.NewFloat(2),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := abiHandler.DecodeHex(testCase.endpoint, testCase.hex)

			require.Nil(t, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func Test_Parse_List_Option(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewVMOutputHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

	testCases := []struct {
		name     string
		endpoint string
		hex      []string
		expected any
	}{
		{
			name:     "Two_levels_nested_list_token_i64",
			endpoint: "option_list_list_list_i64",
			hex: []string{
				"01000000020000000200000003000000006fab0276000000000000001cffffffffffff9ca400000003fffffffffffffffe000000000001e308fffffffff1b3cfdc0000000200000003fffffffffe4932e1000007f24bf1555700000000000000800000000300000010d9e95b690000000000000001ffffffffffffff9c",
			},
			expected: [][]interface{}{
				{
					[]interface{}{int64(1873478262), int64(28), int64(-25436)},
					[]interface{}{int64(-2), int64(123656), int64(-239874084)},
				},
				{
					[]interface{}{int64(-28757279), int64(8737237587287), int64(128)},
					[]interface{}{int64(72375425897), int64(1), int64(-100)},
				},
			},
		},
		{
			name:     "Two_levels_nested_list_token_identifiers",
			endpoint: "option_list_list_list_token",
			hex: []string{
				"010000000200000002000000030000000a43484950532d4e383941000000034b4c56000000034b464900000003000000085446542d3738364a00000008534a412d4c4b394800000008514b552d3748483100000002000000030000000a43484950532d4e383941000000034b4c56000000034b464900000003000000085446542d3738364a00000008534a412d4c4b394800000008514b552d37484831",
			},
			expected: [][]interface{}{
				{
					[]interface{}{string("CHIPS-N89A"), string("KLV"), string("KFI")},
					[]interface{}{string("TFT-786J"), string("SJA-LK9H"), string("QKU-7HH1")},
				},
				{
					[]interface{}{string("CHIPS-N89A"), string("KLV"), string("KFI")},
					[]interface{}{string("TFT-786J"), string("SJA-LK9H"), string("QKU-7HH1")},
				},
			},
		},
		{
			name:     "Two_levels_nested_list_big_int",
			endpoint: "option_list_list_list_big_int",
			hex: []string{
				"0100000002000000020000000300000003" + "87efdb" + "00000002" + "3490" + "00000008" + "c91131a14fc23dac" + "00000003" + "0000000f" + "39bf6e49095ff7dca078957ceb928e" + "0000000f" + "c64091b6f6a008235f876a83146d72" + "00000005" + "0103cf7441" + "00000002" + "00000003" + "00000011" + "1d05b3eb926950d5852179a29534618bf1" + "00000011" + "e2fa4c146d96af2a7ade865d6acb9e740f" + "00000001" + "2a" + "00000003" + "00000003" + "781025" + "00000001" + "54" + "00000006" + "4d9f58c4219f",
			},
			expected: func() [][]interface{} {
				big1, _ := new(big.Int).SetString("-7868453", BaseDecimal)
				big2, _ := new(big.Int).SetString("13456", BaseDecimal)
				big3, _ := new(big.Int).SetString("-3958328028584329812", BaseDecimal)
				big4, _ := new(
					big.Int,
				).SetString("299843598872398459348567275690758798", BaseDecimal)
				big5, _ := new(
					big.Int,
				).SetString("-299843598872398459348567275690758798", BaseDecimal)
				big6, _ := new(big.Int).SetString("4358894657", BaseDecimal)
				big7, _ := new(
					big.Int,
				).SetString("9875768979789979987589732879352104848369", BaseDecimal)
				big8, _ := new(
					big.Int,
				).SetString("-9875768979789979987589732879352104848369", BaseDecimal)
				big9, _ := new(big.Int).SetString("42", BaseDecimal)
				big10, _ := new(big.Int).SetString("7868453", BaseDecimal)
				big11, _ := new(big.Int).SetString("84", BaseDecimal)
				big12, _ := new(big.Int).SetString("85346784387487", BaseDecimal)

				return [][]interface{}{
					{
						[]interface{}{big1, big2, big3},
						[]interface{}{big4, big5, big6},
					},
					{
						[]interface{}{big7, big8, big9},
						[]interface{}{big10, big11, big12},
					},
				}
			}(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := abiHandler.DecodeHex(testCase.endpoint, testCase.hex)

			require.Nil(t, err)
			assert.ElementsMatch(t, testCase.expected, result)
		})
	}
}

func Test_Parse_Tuple(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewVMOutputHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

	testCases := []struct {
		name     string
		endpoint string
		hex      []string
		expected any
	}{
		{
			name:     "Tuple_of_big_int_address_and_i64",
			endpoint: "tuple",
			hex: []string{
				"00000006456fa3a8d3e100000000000000000500af4a12b061e511ca9068af4c83e4918477e6c6ad6a9efffffffe84291d30",
			},
			expected: []interface{}{
				big.NewInt(
					76345789436897,
				), "klv1qqqqqqqqqqqqqpgq4a9p9vrpu5gu4yrg4axg8ey3s3m7d34dd20q5cqaav", int64(-6372647632),
			},
		},
		{
			name:     "Tuple_of_big_int_address_and_i64",
			endpoint: "tuple_nested",
			hex: []string{
				"00000006456fa3a8d3e100000000000000000500af4a12b061e511ca9068af4c83e4918477e6c6ad6a9efffffffe84291d30490000000a68736475676668756973",
			},
			expected: []interface{}{
				big.NewInt(
					76345789436897,
				), "klv1qqqqqqqqqqqqqpgq4a9p9vrpu5gu4yrg4axg8ey3s3m7d34dd20q5cqaav", int64(-6372647632),
				[]interface{}{uint8(73), "hsdugfhuis"},
			},
		},
		{
			name:     "Tuple_with_nested_list",
			endpoint: "tuple_with_nested_list",
			hex: []string{
				"000000090106624c51203c773000000002000000036fab02760000001cffff9ca400000003fffffffe0001e308f1b3cfdc000000084458422d31593641",
			},
			expected: func() []interface{} {
				bigUint, _ := new(big.Int).SetString("18906758096971659056", BaseDecimal)
				return []interface{}{
					bigUint,
					[]interface{}{
						[]interface{}{int32(1873478262), int32(28), int32(-25436)},
						[]interface{}{int32(-2), int32(123656), int32(-239874084)},
					},
					"DXB-1Y6A",
				}
			}(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := abiHandler.DecodeHex(testCase.endpoint, testCase.hex)

			require.Nil(t, err)
			assert.ElementsMatch(t, testCase.expected, result)
		})
	}
}

func Test_Parse_Query_Variadic(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewVMOutputHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

	testCases := []struct {
		name     string
		endpoint string
		hex      []string
		expected any
	}{
		{
			name:     "i32",
			endpoint: "variadic_i32",
			hex:      []string{"CA==", "Vw==", "PEg1Qg==", "+w==", "9HMSAw=="},
			expected: []interface{}{
				int32(8),
				int32(87),
				int32(1011365186),
				int32(-5),
				int32(-193785341),
			},
		},
		{
			name:     "i64",
			endpoint: "variadic_i64",
			hex:      []string{"Ag==", "Vw==", "DgkXR0fTxFI=", "+w==", "5R3QGpWUboM="},
			expected: []interface{}{
				int64(2),
				int64(87),
				int64(1011365186236564562),
				int64(-5),
				int64(-1937163452102185341),
			},
		},
		{
			name:     "u32",
			endpoint: "variadic_u32",
			hex:      []string{"Aw==", "D4gM", "zGGqWQ==", "Orm4"},
			expected: []interface{}{
				uint32(3),
				uint32(1017868),
				uint32(3428952665),
				uint32(3848632),
			},
		},
		{
			name:     "u64",
			endpoint: "variadic_u64",
			hex:      []string{"BTUyvQ==", "AqxVUtAOlQ==", "Ag==", "VYaMaXTsaps="},
			expected: []interface{}{
				uint64(87372477),
				uint64(752432414985877),
				uint64(2),
				uint64(6162767524664208027),
			},
		},
		{
			name:     "big_int",
			endpoint: "variadic_bign",
			hex: []string{
				"BXf2lTU=",
				"CQ==",
				"829Q",
				"CCbd",
				"Ob9uSQlf99ygeJV865KO",
				"xkCRtvagCCNfh2qDFG1y",
			},
			expected: func() []*big.Int {
				bigInt128pos, _ := new(
					big.Int,
				).SetString("299843598872398459348567275690758798", BaseDecimal)
				bigInt128neg, _ := new(
					big.Int,
				).SetString("-299843598872398459348567275690758798", BaseDecimal)

				return []*big.Int{
					big.NewInt(23487485237),
					big.NewInt(9),
					big.NewInt(-823472),
					big.NewInt(534237),
					bigInt128pos,
					bigInt128neg,
				}
			}(),
		},
		{
			name:     "big_uint",
			endpoint: "variadic_bigu",
			hex: []string{
				"6g==",
				"Jmo=",
				"Op6FVA==",
				"eGS0fc8I74w=",
				"AxGCHX6wg+4tJVt9+o9mVjn+2WwpmP2xE2V1Yqk=",
				"CGSmwMkhgOw2eVYWZE02",
			},
			expected: func() []*big.Int {
				bigUintString, _ := new(
					big.Int,
				).SetString("82723675254376576786347243658263586823568268912736457678698677887657", BaseDecimal)
				bigUint128, _ := new(
					big.Int,
				).SetString("43579827367895347689574268789869878", BaseDecimal)

				return []*big.Int{
					big.NewInt(234),
					big.NewInt(9834),
					big.NewInt(983467348),
					big.NewInt(8675257234659798924),
					bigUintString,
					bigUint128,
				}
			}(),
		},
		{
			name:     "token_identifiers",
			endpoint: "variadic_token",
			hex:      []string{"S0xW", "S0ZJ", "S0lELThHOUE=", "RFhCLUg4OEc=", "Q0hJUFMtTjg5QQ=="},
			expected: []interface{}{"KLV", "KFI", "KID-8G9A", "DXB-H88G", "CHIPS-N89A"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := abiHandler.DecodeQuery(testCase.endpoint, testCase.hex)

			require.Nil(t, err)
			assert.ElementsMatch(t, testCase.expected, result)
		})
	}
}

func Test_Parse_Struct(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewVMOutputHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

	bigInt, _ := new(big.Int).SetString("7372987927357988018876673641787893775877874", BaseDecimal)
	bigUint, _ := new(big.Int).SetString("23723672699978725877874", BaseDecimal)

	expected := map[string]interface{}{
		"address_field":  "klv1velayazgrn6mqaqckt7utk9656h8zu3ex4ln8rx7n8p0vy4fd20qmwh4p5",
		"bigint_field":   bigInt,
		"biguint_field":  bigUint,
		"bool_field":     false,
		"i16_field":      int16(8482),
		"i32_field":      int32(425745788),
		"i64_field":      int64(3958328028584329812),
		"i8_field":       int8(82),
		"list_token":     []interface{}{"KLV", "KFI", "KID-8G9A", "DXB-H88G", "CHIPS-N89A"},
		"mngd_buf_field": "testing outputs types",
		"token_field":    "KLV",
		"u16_field":      uint16(8482),
		"u32_field":      uint32(425745788),
		"u64_field":      uint64(3958328028584329812),
		"u8_field":       uint8(82),
		"bigfloat_field": func() *big.Float {
			numerator := big.NewFloat(8)
			denominator := big.NewFloat(3)

			result := new(big.Float).Quo(numerator, denominator)
			result.SetMode(big.RoundingMode(big.Exact))
			return result
		}(),
	}

	hex := "0052212219605d7c36eece5eb03dc25452212219605d7c36eece5eb03dc2540000001574657374696e67206f757470757473207479706573000000034b4c56667fd274481cf5b07418b2fdc5d8baa6ae717239357f338cde99c2f612a96a9e0000000a050610188339c82a68720000001254a3439ee3f36ccadd6573fdfb67310c46f200000012010a0000003500000002aaaaaaaaaaaaa80000000005000000034b4c56000000034b4649000000084b49442d38473941000000084458422d483838470000000a43484950532d4e383941"

	result, err := abiHandler.DecodeHex("struct_test", []string{hex})

	require.Nil(t, err)
	assert.Equal(t, expected, result)
}

func Test_Parse_Struct_With_Simple_And_Tuple_Variant_Enum(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewVMOutputHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

	expected := map[string]interface{}{
		"id":                 uint64(1),
		"category":           "KDA",
		"creator":            "klv1velayazgrn6mqaqckt7utk9656h8zu3ex4ln8rx7n8p0vy4fd20qmwh4p5",
		"creation_timestamp": uint64(1723721192),
		"deadline":           uint64(1723728352),
		"status":             true,
		"subject":            []any{"URB-3P2A"},
	}

	hex := "000000000000000100667fd274481cf5b07418b2fdc5d8baa6ae717239357f338cde99c2f612a96a9e0000000066bde5e80000000066be01e00100000000085552422d33503241"

	result, err := abiHandler.DecodeHex("struct_with_enums", []string{hex})

	require.Nil(t, err)
	assert.Equal(t, expected, result)
}

func Test_Parse_List_of_Struct(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/scFiles/lottery-kda.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewVMOutputHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

	hexOutput := "0000000b226474c28c44db45a9b52945192f533afefe9c4dbf24fccee16dc0553972ac1f000000040ee6b28000000010eaf4866822d00fcfdd9247e953de7f937619c1720336a5c03d6f2939ce63bacd000000040a21fe8000000002eaf4866822d00fcfdd9247e953de7f937619c1720336a5c03d6f2939ce63bacd0000000401312d0000000001eaf4866822d00fcfdd9247e953de7f937619c1720336a5c03d6f2939ce63bacd00000003989680000000012fdc794513bd5d6a96bd1c4369244285a9e9ed840cf6d2480296fc2716ed2f710000000398968000000003226474c28c44db45a9b52945192f533afefe9c4dbf24fccee16dc0553972ac1f0000000401c9c3800000000243016b136874f81ed9c4b1b12c2b1db466127b74d7e276296a6cf7f95992c3bd0000000401c9c380000000019f1354706d75aeb684f26d7dea1fbda17e264c7595cc1eddec0b8968c1be85240000000398968000000001667fd274481cf5b07418b2fdc5d8baa6ae717239357f338cde99c2f612a96a9e00000003989680"

	result, err := abiHandler.DecodeHex("getWinnersInfo", []string{hexOutput})

	expectedResult := []interface{}([]interface{}{
		map[string]interface{}{
			"drawn_ticket_number": uint32(11),
			"prize":               big.NewInt(250000000),
			"winner_address":      "klv1yfj8fs5vgnd5t2d499z3jt6n8tl0a8zdhuj0enhpdhq92wtj4s0snj96jg",
		},
		map[string]interface{}{
			"drawn_ticket_number": uint32(16),
			"prize":               big.NewInt(170000000),
			"winner_address":      "klv1at6gv6pz6q8ulhvjgl548hnljdmpnstjqvm2tspadu5nnnnrhtxsj29zmr",
		},
		map[string]interface{}{
			"drawn_ticket_number": uint32(2),
			"prize":               big.NewInt(20000000),
			"winner_address":      "klv1at6gv6pz6q8ulhvjgl548hnljdmpnstjqvm2tspadu5nnnnrhtxsj29zmr",
		},
		map[string]interface{}{
			"drawn_ticket_number": uint32(1),
			"prize":               big.NewInt(10000000),
			"winner_address":      "klv1at6gv6pz6q8ulhvjgl548hnljdmpnstjqvm2tspadu5nnnnrhtxsj29zmr",
		},
		map[string]interface{}{
			"drawn_ticket_number": uint32(1),
			"prize":               big.NewInt(10000000),
			"winner_address":      "klv19lw8j3gnh4wk494ar3pkjfzzsk57nmvypnmdyjqzjm7zw9hd9acs6qdz6w",
		},
		map[string]interface{}{
			"drawn_ticket_number": uint32(3),
			"prize":               big.NewInt(30000000),
			"winner_address":      "klv1yfj8fs5vgnd5t2d499z3jt6n8tl0a8zdhuj0enhpdhq92wtj4s0snj96jg",
		},
		map[string]interface{}{
			"drawn_ticket_number": uint32(2),
			"prize":               big.NewInt(30000000),
			"winner_address":      "klv1gvqkkymgwnupakwykxcjc2cak3npy7m56l38v2t2dnmljkvjcw7sraj2ny",
		},
		map[string]interface{}{
			"drawn_ticket_number": uint32(1),
			"prize":               big.NewInt(10000000),
			"winner_address":      "klv1nuf4gurdwkhtdp8jd47758aa59lzvnr4jhxpah0vpwyk3sd7s5jqy6mut7",
		},
		map[string]interface{}{
			"drawn_ticket_number": uint32(1),
			"prize":               big.NewInt(10000000),
			"winner_address":      "klv1velayazgrn6mqaqckt7utk9656h8zu3ex4ln8rx7n8p0vy4fd20qmwh4p5",
		},
	})

	require.Nil(t, err)
	assert.ElementsMatch(t, expectedResult, result)
}

func Test_ParseQuery_output_struct(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewVMOutputHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

	bigInt, _ := new(big.Int).SetString("7372987927357988018876673641787893775877874", BaseDecimal)
	bigUint, _ := new(big.Int).SetString("23723672699978725877874", BaseDecimal)

	expected := map[string]interface{}{
		"address_field":  "klv1velayazgrn6mqaqckt7utk9656h8zu3ex4ln8rx7n8p0vy4fd20qmwh4p5",
		"bigint_field":   bigInt,
		"biguint_field":  bigUint,
		"bool_field":     false,
		"i16_field":      int16(8482),
		"i32_field":      int32(425745788),
		"i64_field":      int64(3958328028584329812),
		"i8_field":       int8(82),
		"list_token":     []interface{}{"KLV", "KFI", "KID-8G9A", "DXB-H88G", "CHIPS-N89A"},
		"mngd_buf_field": "testing outputs types",
		"token_field":    "KLV",
		"u16_field":      uint16(8482),
		"u32_field":      uint32(425745788),
		"u64_field":      uint64(3958328028584329812),
		"u8_field":       uint8(82),
		"bigfloat_field": func() *big.Float {
			numerator := big.NewFloat(8)
			denominator := big.NewFloat(3)

			result := new(big.Float).Quo(numerator, denominator)
			result.SetMode(big.RoundingMode(big.Exact))
			return result
		}(),
	}

	queryBase64 := "AFIhIhlgXXw27s5esD3CVFIhIhlgXXw27s5esD3CVAAAABV0ZXN0aW5nIG91dHB1dHMgdHlwZXMAAAADS0xWZn/SdEgc9bB0GLL9xdi6pq5xcjk1fzOM3pnC9hKpap4AAAAKBQYQGIM5yCpocgAAABJUo0Oe4/Nsyt1lc/37ZzEMRvIAAAASAQoAAAA1AAAAAqqqqqqqqqgAAAAABQAAAANLTFYAAAADS0ZJAAAACEtJRC04RzlBAAAACERYQi1IODhHAAAACkNISVBTLU44OUE="

	result, err := abiHandler.DecodeQuery("struct_test", []string{queryBase64})

	require.Nil(t, err)
	assert.Equal(t, expected, result)
}

func Test_ParseQuery_List_of_Struct(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/scFiles/lottery-kda.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewVMOutputHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

	queryOutput := "AAAACyJkdMKMRNtFqbUpRRkvUzr+/pxNvyT8zuFtwFU5cqwfAAAABA7msoAAAAAQ6vSGaCLQD8/dkkfpU95/k3YZwXIDNqXAPW8pOc5jus0AAAAECiH+gAAAAALq9IZoItAPz92SR+lT3n+TdhnBcgM2pcA9byk5zmO6zQAAAAQBMS0AAAAAAer0hmgi0A/P3ZJH6VPef5N2GcFyAzalwD1vKTnOY7rNAAAAA5iWgAAAAAEv3HlFE71dapa9HENpJEKFqenthAz20kgClvwnFu0vcQAAAAOYloAAAAADImR0woxE20WptSlFGS9TOv7+nE2/JPzO4W3AVTlyrB8AAAAEAcnDgAAAAAJDAWsTaHT4HtnEsbEsKx20ZhJ7dNfidilqbPf5WZLDvQAAAAQBycOAAAAAAZ8TVHBtda62hPJtfeofvaF+Jkx1lcwe3ewLiWjBvoUkAAAAA5iWgAAAAAFmf9J0SBz1sHQYsv3F2LqmrnFyOTV/M4zemcL2EqlqngAAAAOYloA="

	result, err := abiHandler.DecodeQuery("getWinnersInfo", []string{queryOutput})

	expectedResult := []interface{}([]interface{}{
		map[string]interface{}{
			"drawn_ticket_number": uint32(11),
			"prize":               big.NewInt(250000000),
			"winner_address":      "klv1yfj8fs5vgnd5t2d499z3jt6n8tl0a8zdhuj0enhpdhq92wtj4s0snj96jg",
		},
		map[string]interface{}{
			"drawn_ticket_number": uint32(16),
			"prize":               big.NewInt(170000000),
			"winner_address":      "klv1at6gv6pz6q8ulhvjgl548hnljdmpnstjqvm2tspadu5nnnnrhtxsj29zmr",
		},
		map[string]interface{}{
			"drawn_ticket_number": uint32(2),
			"prize":               big.NewInt(20000000),
			"winner_address":      "klv1at6gv6pz6q8ulhvjgl548hnljdmpnstjqvm2tspadu5nnnnrhtxsj29zmr",
		},
		map[string]interface{}{
			"drawn_ticket_number": uint32(1),
			"prize":               big.NewInt(10000000),
			"winner_address":      "klv1at6gv6pz6q8ulhvjgl548hnljdmpnstjqvm2tspadu5nnnnrhtxsj29zmr",
		},
		map[string]interface{}{
			"drawn_ticket_number": uint32(1),
			"prize":               big.NewInt(10000000),
			"winner_address":      "klv19lw8j3gnh4wk494ar3pkjfzzsk57nmvypnmdyjqzjm7zw9hd9acs6qdz6w",
		},
		map[string]interface{}{
			"drawn_ticket_number": uint32(3),
			"prize":               big.NewInt(30000000),
			"winner_address":      "klv1yfj8fs5vgnd5t2d499z3jt6n8tl0a8zdhuj0enhpdhq92wtj4s0snj96jg",
		},
		map[string]interface{}{
			"drawn_ticket_number": uint32(2),
			"prize":               big.NewInt(30000000),
			"winner_address":      "klv1gvqkkymgwnupakwykxcjc2cak3npy7m56l38v2t2dnmljkvjcw7sraj2ny",
		},
		map[string]interface{}{
			"drawn_ticket_number": uint32(1),
			"prize":               big.NewInt(10000000),
			"winner_address":      "klv1nuf4gurdwkhtdp8jd47758aa59lzvnr4jhxpah0vpwyk3sd7s5jqy6mut7",
		},
		map[string]interface{}{
			"drawn_ticket_number": uint32(1),
			"prize":               big.NewInt(10000000),
			"winner_address":      "klv1velayazgrn6mqaqckt7utk9656h8zu3ex4ln8rx7n8p0vy4fd20qmwh4p5",
		},
	})

	require.Nil(t, err)
	assert.ElementsMatch(t, expectedResult, result)
}

func Test_ParseQuery_MultiValue_of_struct_and_nested_list(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewVMOutputHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

	bigInt, _ := new(big.Int).SetString("7372987927357988018876673641787893775877874", BaseDecimal)
	bigUint, _ := new(big.Int).SetString("23723672699978725877874", BaseDecimal)

	expectedMap := map[string]interface{}{
		"address_field":  "klv1velayazgrn6mqaqckt7utk9656h8zu3ex4ln8rx7n8p0vy4fd20qmwh4p5",
		"bigint_field":   bigInt,
		"biguint_field":  bigUint,
		"bool_field":     false,
		"i16_field":      int16(8482),
		"i32_field":      int32(425745788),
		"i64_field":      int64(3958328028584329812),
		"i8_field":       int8(82),
		"list_token":     []interface{}{"KLV", "KFI", "KID-8G9A", "DXB-H88G", "CHIPS-N89A"},
		"mngd_buf_field": "testing outputs types",
		"token_field":    "KLV",
		"u16_field":      uint16(8482),
		"u32_field":      uint32(425745788),
		"u64_field":      uint64(3958328028584329812),
		"u8_field":       uint8(82),
		"bigfloat_field": func() *big.Float {
			numerator := big.NewFloat(8)
			denominator := big.NewFloat(3)

			result := new(big.Float).Quo(numerator, denominator)
			result.SetMode(big.RoundingMode(big.Exact))
			return result
		}(),
	}

	expectedListI32 := []interface{}{
		[]interface{}{int32(1873478262), int32(28), int32(-25436)},
		[]interface{}{int32(-2), int32(123656), int32(-239874084)},
	}

	expectedListBigInt := func() []interface{} {
		big1, _ := new(big.Int).SetString("-7868453", BaseDecimal)
		big2, _ := new(big.Int).SetString("13456", BaseDecimal)
		big3, _ := new(big.Int).SetString("-3958328028584329812", BaseDecimal)
		big4, _ := new(big.Int).SetString("299843598872398459348567275690758798", BaseDecimal)
		big5, _ := new(big.Int).SetString("-299843598872398459348567275690758798", BaseDecimal)
		big6, _ := new(big.Int).SetString("4358894657", BaseDecimal)
		big7, _ := new(big.Int).SetString("9875768979789979987589732879352104848369", BaseDecimal)
		big8, _ := new(big.Int).SetString("-9875768979789979987589732879352104848369", BaseDecimal)
		big9, _ := new(big.Int).SetString("42", BaseDecimal)
		big10, _ := new(big.Int).SetString("7868453", BaseDecimal)
		big11, _ := new(big.Int).SetString("84", BaseDecimal)
		big12, _ := new(big.Int).SetString("85346784387487", BaseDecimal)

		return []interface{}{
			[]interface{}{
				[]interface{}{big1, big2, big3},
				[]interface{}{big4, big5, big6},
			},
			[]interface{}{
				[]interface{}{big7, big8, big9},
				[]interface{}{big10, big11, big12},
			},
		}
	}()

	expectedOutput := []interface{}{expectedMap, expectedListI32, expectedListBigInt}

	hexInputs := []string{
		"AFIhIhlgXXw27s5esD3CVFIhIhlgXXw27s5esD3CVAAAABV0ZXN0aW5nIG91dHB1dHMgdHlwZXMAAAADS0xWZn/SdEgc9bB0GLL9xdi6pq5xcjk1fzOM3pnC9hKpap4AAAAKBQYQGIM5yCpocgAAABJUo0Oe4/Nsyt1lc/37ZzEMRvIAAAASAQoAAAA1AAAAAqqqqqqqqqgAAAAABQAAAANLTFYAAAADS0ZJAAAACEtJRC04RzlBAAAACERYQi1IODhHAAAACkNISVBTLU44OUE=",
		"AAAAA2+rAnYAAAAc//+cpAAAAAP////+AAHjCPGzz9w=",
		"AAAAAgAAAAMAAAADh+/bAAAAAjSQAAAACMkRMaFPwj2sAAAAAwAAAA85v25JCV/33KB4lXzrko4AAAAPxkCRtvagCCNfh2qDFG1yAAAABQEDz3RBAAAAAgAAAAMAAAARHQWz65JpUNWFIXmilTRhi/EAAAAR4vpMFG2Wryp63oZdasuedA8AAAABKgAAAAMAAAADeBAlAAAAAVQAAAAGTZ9YxCGf",
	}

	result, err := abiHandler.DecodeQuery("multi_value_nested_list_struct", hexInputs)

	require.Nil(t, err)
	assert.ElementsMatch(t, expectedOutput, result)
}

func Test_DecodeStruct_AddLiquidityEvent(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/scFiles/pair.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewVMOutputHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

	hex := "50630e5d63c70e1b7985344af7e33699736771d0a5db02929653c89b9ee444690000000844564b2d33345a480000000301e240000000034b4c560000000230390000000b44564b4b4c562d5047594c00000002987e000000061cd96aee9ede000000065b3abdcdbdea00000006091f79615fca0000000001b25b9200000000000014a100000000698b92c8"

	decodedValue, err := abiHandler.DecodeStruct("AddLiquidityEvent", hex)
	require.Nil(t, err)

	expected := map[string]interface{}{
		"caller":                "klv12p3suhtrcu8pk7v9x3900cekn9ekwuws5hds9y5k20yfh8hyg35ss3s8ex",
		"first_token_id":        "DVK-34ZH",
		"first_token_amount":    big.NewInt(123456),
		"second_token_id":       "KLV",
		"second_token_amount":   big.NewInt(12345),
		"lp_token_id":           "DVKKLV-PGYL",
		"lp_token_amount":       big.NewInt(39038),
		"lp_supply":             big.NewInt(31720127504094),
		"first_token_reserves":  big.NewInt(100307850608106),
		"second_token_reserves": big.NewInt(10030785060810),
		"block":                 uint64(28466066),
		"epoch":                 uint64(5281),
		"timestamp":             uint64(1770754760),
	}

	assert.Equal(t, expected, decodedValue)
}

func Test_DecodeStruct_KdaTokenPayment(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/scFiles/pair.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewVMOutputHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error loading abi", errLoad)

	testCases := []struct {
		name     string
		hex      string
		expected map[string]interface{}
	}{
		{
			name: "BABYDGKO_large_amount",
			hex:  "0000000d4241425944474b4f2d3353363700000000000000000000000602a6db03cc64",
			expected: map[string]interface{}{
				"token_identifier": "BABYDGKO-3S67",
				"token_nonce":      uint64(0),
				"amount":           big.NewInt(2915662285924),
			},
		},
		{
			name: "KFI_nonzero_amount",
			hex:  "000000034b464900000000000000000000000313054f",
			expected: map[string]interface{}{
				"token_identifier": "KFI",
				"token_nonce":      uint64(0),
				"amount":           big.NewInt(1246543),
			},
		},
		{
			name: "KLV_nonzero_amount",
			hex:  "000000034b4c5600000000000000000000000364e728",
			expected: map[string]interface{}{
				"token_identifier": "KLV",
				"token_nonce":      uint64(0),
				"amount":           big.NewInt(6612776),
			},
		},
		{
			name: "KFI_zero_amount",
			hex:  "000000034b4649000000000000000000000000",
			expected: map[string]interface{}{
				"token_identifier": "KFI",
				"token_nonce":      uint64(0),
				"amount":           big.NewInt(0),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := abiHandler.DecodeStruct("KdaTokenPayment", tc.hex)
			require.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}
