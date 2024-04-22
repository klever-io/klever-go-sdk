package provider_test

import (
	"math/big"
	"os"
	"testing"

	"github.com/klever-io/klever-go-sdk/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const BaseDecimal = 10

func Test_Decode_Single_Value(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewSCAbiHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

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
				bigInt128, _ := new(big.Int).SetString("299843598872398459348567275690758798", BaseDecimal)
				return bigInt128
			}(),
		},
		{
			name:     "BigInt_negative_from_uint128",
			endpoint: "big_minus_u_number",
			hex:      "c64091b6f6a008235f876a83146d72",
			expected: func() *big.Int {
				bigInt128, _ := new(big.Int).SetString("-299843598872398459348567275690758798", BaseDecimal)
				return bigInt128
			}(),
		},
		{
			name:     "BigInt_from_random_string",
			endpoint: "big_s_number",
			hex:      "393833343735393337343536383932343739363738383930313736393831393038353637383935373639303738353132393836373938323537",
			expected: func() *big.Int {
				bigIntRandomString, _ := new(big.Int).SetString("983475937456892479678890176981908567895769078512986798257", BaseDecimal)
				return bigIntRandomString
			}(),
		},
		{
			name:     "BigInt_from_random_negative_string",
			endpoint: "big_minus_s_number",
			hex:      "2d393833343735393337343536383932343739363738383930313736393831393038353637383935373639303738353132393836373938323537",
			expected: func() *big.Int {
				bigIntRandomString, _ := new(big.Int).SetString("-983475937456892479678890176981908567895769078512986798257", BaseDecimal)
				return bigIntRandomString
			}(),
		},
		{
			name:     "BigUint_from_string_10",
			endpoint: "big_u_s_10",
			hex:      "3130",
			expected: big.NewInt(10),
		},
		{
			name:     "BigUint_from_random_string",
			endpoint: "big_u_s_number",
			hex:      "3832373432383733363433343735393733353933343736393733393637393337363938333435373836393833393035363938393739373839373839",
			expected: func() *big.Int {
				bigUintRandomString, _ := new(big.Int).SetString("82742873643475973593476973967937698345786983905698979789789", BaseDecimal)
				return bigUintRandomString
			}(),
		},
		{
			name:     "BigUint_from_u8",
			endpoint: "big_u8",
			hex:      "52",
			expected: big.NewInt(82),
		},
		{
			name:     "BigUint_from_u16",
			endpoint: "big_u16",
			hex:      "2122",
			expected: big.NewInt(8482),
		},
		{
			name:     "BigUint_from_u32",
			endpoint: "big_u32",
			hex:      "19605d7c",
			expected: big.NewInt(425745788),
		},
		{
			name:     "BigUint_from_u64",
			endpoint: "big_u64",
			hex:      "36eece5eb03dc254",
			expected: big.NewInt(3958328028584329812),
		},
		{
			name:     "BigUint_from_u128",
			endpoint: "big_u128",
			hex:      "1dc7766516260b32b52ff11612d5710e",
			expected: func() *big.Int {
				bigUintFromU128, _ := new(big.Int).SetString("39583280285843298128735477835272384782", BaseDecimal)
				return bigUintFromU128
			}(),
		},
		{
			name:     "Int8_positive",
			endpoint: "number_i8",
			hex:      "52",
			expected: int8(82),
		},
		{
			name:     "Int8_negative",
			endpoint: "number_minus_i8",
			hex:      "ae",
			expected: int8(-82),
		},
		{
			name:     "Int16_positive",
			endpoint: "number_i16",
			hex:      "2122",
			expected: int16(8482),
		},
		{
			name:     "Int16_negative",
			endpoint: "number_minus_i16",
			hex:      "dede",
			expected: int16(-8482),
		},
		{
			name:     "Int32_positive",
			endpoint: "number_i32",
			hex:      "19605d7c",
			expected: int32(425745788),
		},
		{
			name:     "Int32_negative",
			endpoint: "number_minus_i32",
			hex:      "e69fa284",
			expected: int32(-425745788),
		},
		{
			name:     "Int64_positive",
			endpoint: "number_i64",
			hex:      "36eece5eb03dc254",
			expected: int64(3958328028584329812),
		},
		{
			name:     "Int64_negative",
			endpoint: "number_minus_i64",
			hex:      "c91131a14fc23dac",
			expected: int64(-3958328028584329812),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := abiHandler.Decode(testCase.endpoint, testCase.hex)

			require.Nil(t, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func Test_Decode_List(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewSCAbiHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error loading abi", errLoad)

	testCases := []struct {
		name     string
		endpoint string
		hex      string
		expected any
	}{
		{
			name:     "Token_identifiers",
			endpoint: "list_token_identifier",
			hex:      "000000034b4c56000000034b4649000000084b49442d38473941000000084458422d483838470000000a43484950532d4e383941",
			expected: []string{"KLV", "KFI", "KID-8G9A", "DXB-H88G", "CHIPS-N89A"},
		},
		{
			name:     "int32",
			endpoint: "list_int32",
			hex:      "000000080000005700000065fffffffb",
			expected: []int32{int32(8), int32(87), int32(101), int32(-5)},
		},
		{
			name:     "int64",
			endpoint: "list_int64",
			hex:      "000000000000000200000000000000570e09174747d3c452fffffffffffffffbe51dd01a95946e83",
			expected: []int64{int64(2), int64(87), int64(1011365186236564562), int64(-5), int64(-1937163452102185341)},
		},
		{
			name:     "uint16",
			endpoint: "list_u16",
			hex:      "3fd600a7000203f3",
			expected: []uint16{uint16(16342), uint16(167), uint16(2), uint16(1011)},
		},
		{
			name:     "uint32",
			endpoint: "list_u32",
			hex:      "00000003000f880ccc61aa59003ab9b8",
			expected: []uint32{uint32(3), uint32(1017868), uint32(3428952665), uint32(3848632)},
		},
		{
			name:     "uint64",
			endpoint: "list_u64",
			hex:      "00000000053532bd0002ac5552d00e95000000000000000255868c6974ec6a9b",
			expected: []uint64{uint64(87372477), uint64(752432414985877), uint64(2), uint64(6162767524664208027)},
		},
		{
			name:     "BigInt",
			endpoint: "list_bign",
			hex:      "000000050577f695350000000109000000072d383233343732000000063533343233370000000f39bf6e49095ff7dca078957ceb928e0000000fc64091b6f6a008235f876a83146d72",
			expected: func() []*big.Int {
				bigInt128pos, _ := new(big.Int).SetString("299843598872398459348567275690758798", BaseDecimal)
				bigInt128neg, _ := new(big.Int).SetString("-299843598872398459348567275690758798", BaseDecimal)

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
			hex:      "00000001ea00000002266a000000043a9e8554000000087864b47dcf08ef8c0000004438323732333637353235343337363537363738363334373234333635383236333538363832333536383236383931323733363435373637383639383637373838373635370000000f0864a6c0c92180ec36795616644d36",
			expected: func() []*big.Int {
				bigUintString, _ := new(big.Int).SetString("82723675254376576786347243658263586823568268912736457678698677887657", BaseDecimal)
				bigUint128, _ := new(big.Int).SetString("43579827367895347689574268789869878", BaseDecimal)

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
			hex:      "667fd274481cf5b07418b2fdc5d8baa6ae717239357f338cde99c2f612a96a9e667fd274481cf5b07418b2fdc5d8baa6ae717239357f338cde99c2f612a96a9e",
			expected: []string{"klv1velayazgrn6mqaqckt7utk9656h8zu3ex4ln8rx7n8p0vy4fd20qmwh4p5", "klv1velayazgrn6mqaqckt7utk9656h8zu3ex4ln8rx7n8p0vy4fd20qmwh4p5"},
		},
		{
			name:     "Boolean",
			endpoint: "list_bool",
			hex:      "01000001",
			expected: []bool{true, false, false, true},
		},
		{
			name:     "Nested_list_boolean",
			endpoint: "list_list_bool",
			hex:      "00000004010000010000000401000100",
			expected: [][]interface{}{
				{true, false, false, true},
				{true, false, true, false},
			},
		},
		{
			name:     "Nested_list_int32",
			endpoint: "list_list_i32",
			hex:      "000000036fab02760000001cffff9ca400000003fffffffe0001e308f1b3cfdc",
			expected: [][]interface{}{
				{int32(1873478262), int32(28), int32(-25436)},
				{int32(-2), int32(123656), int32(-239874084)},
			},
		},
		{
			name:     "Nested_list_int64",
			endpoint: "list_list_i64",
			hex:      "0000000319ffee93a36dc12a000000000000001cd7e571502441e18400000003fffffffffffffffe000000000001e308fffffffff1b3cfdc",
			expected: [][]interface{}{
				{int64(1873478287878897962), int64(28), int64(-2889778996868685436)},
				{int64(-2), int64(123656), int64(-239874084)},
			},
		},
		{
			name:     "Nested_list_tokens_identifiers",
			endpoint: "list_list_tokens",
			hex:      "00000003000000034b4c56000000034b4649000000084b49442d3847394100000003000000084458422d483838470000000a43484950532d4e383941000000084646542d32424836",
			expected: [][]interface{}{
				{"KLV", "KFI", "KID-8G9A"},
				{"DXB-H88G", "CHIPS-N89A", "FFT-2BH6"},
			},
		},
		{
			name:     "Two_levels_nested_list_int64",
			endpoint: "list_list_list_i64",
			hex:      "0000000200000003000000006fab0276000000000000001cffffffffffff9ca400000003fffffffffffffffe000000000001e308fffffffff1b3cfdc0000000200000003fffffffffe4932e1000007f24bf1555700000000000000800000000300000010d9e95b690000000000000001ffffffffffffff9c",
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
			hex:      "00000002000000030000000a43484950532d4e383941000000034b4c56000000034b464900000003000000085446542d3738364a00000008534a412d4c4b394800000008514b552d3748483100000002000000030000000a43484950532d4e383941000000034b4c56000000034b464900000003000000085446542d3738364a00000008534a412d4c4b394800000008514b552d37484831",
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
			hex:      "00000002000000030000000387efdb00000002349000000008c91131a14fc23dac000000030000000f39bf6e49095ff7dca078957ceb928e0000000fc64091b6f6a008235f876a83146d72000000050103cf744100000002000000030000002839383735373638393739373839393739393837353839373332383739333532313034383438333639000000292d3938373537363839373937383939373939383735383937333238373933353231303438343833363900000002343200000003000000037810250000000154000000064d9f58c4219f",
			expected: func() [][]interface{} {
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
			result, err := abiHandler.Decode(testCase.endpoint, testCase.hex)

			require.Nil(t, err)
			assert.ElementsMatch(t, testCase.expected, result)
		})
	}
}

func Test_Decode_Option(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewSCAbiHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

	testCases := []struct {
		name     string
		endpoint string
		hex      string
		expected any
	}{
		{
			name:     "null",
			endpoint: "option_bytes_null",
			hex:      "",
			expected: nil,
		},
		{
			name:     "i8",
			endpoint: "option_i8",
			hex:      "0152",
			expected: int8(82),
		},
		{
			name:     "i16",
			endpoint: "option_i16",
			hex:      "012122",
			expected: int16(8482),
		},
		{
			name:     "i32",
			endpoint: "option_i32",
			hex:      "0119605d7c",
			expected: int32(425745788),
		},
		{
			name:     "i64",
			endpoint: "option_i64",
			hex:      "0136eece5eb03dc254",
			expected: int64(3958328028584329812),
		},
		{
			name:     "u8",
			endpoint: "option_u8",
			hex:      "0152",
			expected: uint8(82),
		},
		{
			name:     "u16",
			endpoint: "option_u16",
			hex:      "012122",
			expected: uint16(8482),
		},
		{
			name:     "u32",
			endpoint: "option_u32",
			hex:      "0119605d7c",
			expected: uint32(425745788),
		},
		{
			name:     "u64",
			endpoint: "option_u64",
			hex:      "0136eece5eb03dc254",
			expected: uint64(3958328028584329812),
		},
		{
			name:     "big_int",
			endpoint: "option_bigint",
			hex:      "01000000050577f69535",
			expected: big.NewInt(23487485237),
		},
		{
			name:     "big_int_from_buffer",
			endpoint: "option_bigint_from_buffer",
			hex:      "010000001332333438373738343735383733343835323337",
			expected: func() *big.Int {
				big, _ := new(big.Int).SetString("2348778475873485237", BaseDecimal)
				return big
			}(),
		},
		{
			name:     "bigint_from_biguint_plus",
			endpoint: "option_bigint_from_biguint_plus",
			hex:      "010000000f39bf6e49095ff7dca078957ceb928e",
			expected: func() *big.Int {
				big, _ := new(big.Int).SetString("299843598872398459348567275690758798", BaseDecimal)
				return big
			}(),
		},
		{
			name:     "bigint_from_biguint_minus",
			endpoint: "option_bigint_from_biguint_minus",
			hex:      "010000000fc64091b6f6a008235f876a83146d72",
			expected: func() *big.Int {
				big, _ := new(big.Int).SetString("-299843598872398459348567275690758798", BaseDecimal)
				return big
			}(),
		},
		{
			name:     "biguint",
			endpoint: "option_biguint",
			hex:      "010000000537952bd072",
			expected: big.NewInt(238725877874),
		},
		{
			name:     "biguint_from_buffer",
			endpoint: "option_biguint_from_buffer",
			hex:      "0100000015323338373438323734383237393235383737383734",
			expected: func() *big.Int {
				big, _ := new(big.Int).SetString("238748274827925877874", BaseDecimal)
				return big
			}(),
		},
		{
			name:     "address",
			endpoint: "option_address",
			hex:      "01667fd274481cf5b07418b2fdc5d8baa6ae717239357f338cde99c2f612a96a9e",
			expected: "klv1velayazgrn6mqaqckt7utk9656h8zu3ex4ln8rx7n8p0vy4fd20qmwh4p5",
		},
		{
			name:     "bytes",
			endpoint: "option_managed_buffer",
			hex:      "010000001574657374696e67206f757470757473207479706573",
			expected: "testing outputs types",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := abiHandler.Decode(testCase.endpoint, testCase.hex)

			require.Nil(t, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func Test_Decode_List_Option(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewSCAbiHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

	testCases := []struct {
		name     string
		endpoint string
		hex      string
		expected any
	}{
		{
			name:     "Two_levels_nested_list_token_i64",
			endpoint: "option_list_list_list_i64",
			hex:      "01000000020000000200000003000000006fab0276000000000000001cffffffffffff9ca400000003fffffffffffffffe000000000001e308fffffffff1b3cfdc0000000200000003fffffffffe4932e1000007f24bf1555700000000000000800000000300000010d9e95b690000000000000001ffffffffffffff9c",
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
			hex:      "010000000200000002000000030000000a43484950532d4e383941000000034b4c56000000034b464900000003000000085446542d3738364a00000008534a412d4c4b394800000008514b552d3748483100000002000000030000000a43484950532d4e383941000000034b4c56000000034b464900000003000000085446542d3738364a00000008534a412d4c4b394800000008514b552d37484831",
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
			hex:      "010000000200000002000000030000000387efdb00000002349000000008c91131a14fc23dac000000030000000f39bf6e49095ff7dca078957ceb928e0000000fc64091b6f6a008235f876a83146d72000000050103cf744100000002000000030000002839383735373638393739373839393739393837353839373332383739333532313034383438333639000000292d3938373537363839373937383939373939383735383937333238373933353231303438343833363900000002343200000003000000037810250000000154000000064d9f58c4219f",
			expected: func() [][]interface{} {
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
			result, err := abiHandler.Decode(testCase.endpoint, testCase.hex)

			require.Nil(t, err)
			assert.ElementsMatch(t, testCase.expected, result)
		})
	}
}

func Test_Decode_Tuple(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewSCAbiHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

	testCases := []struct {
		name     string
		endpoint string
		hex      string
		expected any
	}{
		{
			name:     "Tuple_of_big_int_address_and_i64",
			endpoint: "tuple",
			hex:      "0000000e373633343537383934333638393700000000000000000500af4a12b061e511ca9068af4c83e4918477e6c6ad6a9efffffffe84291d30",
			expected: []interface{}{
				big.NewInt(76345789436897), "klv1qqqqqqqqqqqqqpgq4a9p9vrpu5gu4yrg4axg8ey3s3m7d34dd20q5cqaav", int64(-6372647632),
			},
		},
		{
			name:     "Tuple_of_big_int_address_and_i64",
			endpoint: "tuple_nested",
			hex:      "0000000e373633343537383934333638393700000000000000000500af4a12b061e511ca9068af4c83e4918477e6c6ad6a9efffffffe84291d30490000000a68736475676668756973",
			expected: []interface{}{
				big.NewInt(76345789436897), "klv1qqqqqqqqqqqqqpgq4a9p9vrpu5gu4yrg4axg8ey3s3m7d34dd20q5cqaav", int64(-6372647632),
				[]interface{}{uint8(73), "hsdugfhuis"},
			},
		},
		{
			name:     "Tuple_with_nested_list",
			endpoint: "tuple_with_nested_list",
			hex:      "000000090106624c51203c773000000002000000036fab02760000001cffff9ca400000003fffffffe0001e308f1b3cfdc000000084458422d31593641",
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
			result, err := abiHandler.Decode(testCase.endpoint, testCase.hex)

			require.Nil(t, err)
			assert.ElementsMatch(t, testCase.expected, result)
		})
	}
}

func Test_Decode_Variadic(t *testing.T) {
	jsonAbi, errOpen := os.Open("../cmd/demo/smartContracts/decode/example.abi.json")
	require.Nil(t, errOpen, "error opening abi", errOpen)
	defer jsonAbi.Close()

	abiHandler := provider.NewSCAbiHandler()

	errLoad := abiHandler.LoadAbi(jsonAbi)
	require.Nil(t, errLoad, "error opening abi", errLoad)

	testCases := []struct {
		name     string
		endpoint string
		hex      string
		expected any
	}{
		{
			name:     "Variadic_with_only_one_i32_value",
			endpoint: "var_i32",
			hex:      "049075ea",
			expected: int32(76576234),
		},
		{
			name:     "Variadic_with_multiple_i32_value",
			endpoint: "var_i32",
			hex:      "049075ea",
			expected: int32(76576234),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := abiHandler.Decode(testCase.endpoint, testCase.hex)

			require.Nil(t, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := abiHandler.Decode(testCase.endpoint, testCase.hex)

			require.Nil(t, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}
