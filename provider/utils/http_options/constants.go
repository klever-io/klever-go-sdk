package http_options

type OptionType int

const (
	Header OptionType = iota + 1
	MarshalError
)
