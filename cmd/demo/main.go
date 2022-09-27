package main

import (
	"fmt"
	"time"

	"github.com/klever-io/klever-go-sdk/provider"
	"github.com/klever-io/klever-go-sdk/provider/network"
	"github.com/klever-io/klever-go-sdk/provider/utils"
)

func main() {

	net := network.NewNetworkConfig(network.MainNet)
	httpClient := utils.NewHttpClient(10 * time.Second)
	kc, err := provider.NewKleverChain(net, httpClient)
	if err != nil {
		panic(err)
	}

	acc, err := kc.GetAccount("klv1usdnywjhrlv4tcyu6stxpl6yvhplg35nepljlt4y5r7yppe8er4qujlazy")
	if err != nil {
		panic(err)
	}

	fmt.Println(acc.String())
}
