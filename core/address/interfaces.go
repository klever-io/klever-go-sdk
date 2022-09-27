package address

type Address interface {
	Bech32() string
	Hex() string
	Bytes() []byte
	IsInterfaceNil() bool
}
