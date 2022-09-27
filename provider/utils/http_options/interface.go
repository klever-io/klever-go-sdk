package http_options

type IOptions interface {
	GetType() OptionType
	GetHeaders() [][2]string
	GetMarshalError() bool
}
