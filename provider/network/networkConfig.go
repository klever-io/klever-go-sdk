package network

type Network int

const (
	LocalNet Network = iota
	MainNet
	TestNet
	DevNet
	Custom
)

type networkConfig struct {
	network     Network
	uriAPI      string
	uriNode     string
	uriExplorer string
}

func NewNetworkConfigCustom(
	APIUri string,
	NodeUri string,
	ExplorerUri string,
) NetworkConfig {
	return &networkConfig{
		network:     Custom,
		uriAPI:      APIUri,
		uriNode:     NodeUri,
		uriExplorer: ExplorerUri,
	}
}

func NewNetworkConfig(network Network) NetworkConfig {
	APIUri := "http://localhost:8701"
	NodeUri := "http://localhost:8701"
	ExplorerUri := "http://localhost"
	switch network {
	case LocalNet:
	case MainNet:
		APIUri = "https://api.mainnet.klever.org"
		NodeUri = "https://node.mainnet.klever.org"
		ExplorerUri = "https://kleverscan.org/"
	case TestNet:
		APIUri = "https://api.testnet.klever.org"
		NodeUri = "https://node.testnet.klever.org"
		ExplorerUri = "https://testnet.kleverscan.org/"
	case DevNet:
		APIUri = "https://api.devnet.klever.org"
		NodeUri = "https://node.devnet.klever.org"
		ExplorerUri = "https://devnet.kleverscan.org/"
	default:
		panic("invalid network config")
	}

	return &networkConfig{
		network:     network,
		uriAPI:      APIUri,
		uriNode:     NodeUri,
		uriExplorer: ExplorerUri,
	}
}

func (n *networkConfig) GetNetwork() Network {
	return n.network
}

func (n *networkConfig) GetAPIUri() string {
	return n.uriAPI
}

func (n *networkConfig) GetNodeUri() string {
	return n.uriNode
}

func (n *networkConfig) GetExplorerUri() string {
	return n.uriExplorer
}
