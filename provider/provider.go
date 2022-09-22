package provider

import (
	"github.com/klever-io/klever-go-sdk/provider/network"
	"github.com/klever-io/klever-go-sdk/provider/utils"
)

type kleverChain struct {
	networkConfig network.NetworkConfig
	httpClient    utils.HttpClient
}

func NewKleverChain(network network.NetworkConfig, httpClient utils.HttpClient) (KleverChain, error) {

	return &kleverChain{
		networkConfig: network,
		httpClient:    httpClient,
	}, nil
}
