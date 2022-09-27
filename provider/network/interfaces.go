package network

type NetworkConfig interface {
	GetNetwork() Network
	GetAPIUri() string
	GetNodeUri() string
	GetExplorerUri() string
}
