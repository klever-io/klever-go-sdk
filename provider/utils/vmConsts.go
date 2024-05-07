package utils

const LengthOfCodeMetadata = 2

// Const group for the first byte of the metadata
const (
	// MetadataUpgradeable is the bit for upgradable flag
	MetadataUpgradeable = 1
	// MetadataReadable is the bit for readable flag
	MetadataReadable = 4
)

// Const group for the second byte of the metadata
const (
	// MetadataPayable is the bit for payable flag
	MetadataPayable = 2
	// MetadataPayableBySC is the bit for payable flag
	MetadataPayableBySC = 4
)

const DefaultVMType = "0500"

type SCType int32

const (
	SmartContractInvoke SCType = iota
	SmartContractDeploy
)

// Hex length
const (
	HexLength8Bits  int = 2
	HexLength16Bits int = 4
	HexLength32Bits int = 8
	HexLength64Bits int = 16

	AddressHexLen int = 64
)

// Bits count
const (
	Bits8   int = 8
	Bits16  int = 16
	Bits32  int = 32
	Bits64  int = 64
	Bits128 int = 128
)

// Numerical bases
const (
	BaseHex     int = 16
	BaseDecimal int = 10
)

// Types wrappers
const (
	List     string = "List"
	Option   string = "Option"
	Tuple    string = "tuple"
	Variadic string = "variadic"
)

// Possible Types
const (
	Int8            string = "i8"
	Uint8           string = "u8"
	Int16           string = "i16"
	Uint16          string = "u16"
	Int32           string = "i32"
	Isize           string = "isize"
	Uint32          string = "u32"
	Usize           string = "usize"
	Int64           string = "i64"
	Uint64          string = "u64"
	BigInt          string = "BigInt"
	BigUint         string = "BigUint"
	BigFloat        string = "BigFloat"
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
)

// Others
const (
	LengthHexSizer int = 8
	BitsByHexDigit int = 4
)

const BigFloatVMPrecision uint = 53
