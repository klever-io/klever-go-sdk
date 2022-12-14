package main

import (
	"fmt"
	"time"

	"github.com/klever-io/klever-go-sdk/provider"
	"github.com/klever-io/klever-go-sdk/provider/network"
	"github.com/klever-io/klever-go-sdk/provider/utils"
)

func main() {

	net := network.NewNetworkConfig(network.TestNet)
	httpClient := utils.NewHttpClient(10 * time.Second)
	kc, err := provider.NewKleverChain(net, httpClient)
	if err != nil {
		panic(err)
	}

	acc, err := kc.GetAccount("klv1mt8yw657z6nk9002pccmwql8w90k0ac6340cjqkvm9e7lu0z2wjqudt69s")
	if err != nil {
		panic(err)
	}

	fmt.Println(acc.String())
}
