package provider

import (
	"github.com/klever-io/klever-go-sdk/provider/network"
	"github.com/klever-io/klever-go-sdk/provider/tools/hasher"
	"github.com/klever-io/klever-go-sdk/provider/tools/marshal"
	"github.com/klever-io/klever-go-sdk/provider/utils"
)

type kleverChain struct {
	networkConfig network.NetworkConfig
	httpClient    utils.HttpClient

	marshalizer marshal.Marshalizer
	hasher      hasher.Hasher
}

func NewKleverChain(network network.NetworkConfig, httpClient utils.HttpClient) (KleverChain, error) {
	hasher, err := hasher.NewHasher()
	if err != nil {
		return nil, err
	}

	marshalizer := marshal.NewProtoMarshalizer()

	return &kleverChain{
		networkConfig: network,
		httpClient:    httpClient,
		hasher:        hasher,
		marshalizer:   marshalizer,
	}, nil
}

func (kc *kleverChain) GetHasher() hasher.Hasher {
	return kc.hasher
}

func (kc *kleverChain) GetMarshalizer() marshal.Marshalizer {
	return kc.marshalizer
}
