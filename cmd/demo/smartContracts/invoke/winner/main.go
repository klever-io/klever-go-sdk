package main

import (
	"fmt"

	"github.com/klever-io/klever-go-sdk/cmd/demo"
)

func main() {
	accounts, wallets, kc, err := demo.InitWallets()

	if err != nil {
		panic(err)
	}

	base := accounts[0].NewBaseTX()

	tx, err := kc.InvokeSmartContract(
		base,
		"klv1qqqqqqqqqqqqqpgq4f6qkvv34kr50w5lucdn6z0238r7m96nhtxspjw22t",
		"determine_winner",
		map[string]int64{},
		"demo",
	)

	if err != nil {
		panic(err)
	}

	err = tx.Sign(wallets[0])
	if err != nil {
		panic(err)
	}

	hash, err := tx.Broadcast(kc)
	if err != nil {
		panic(err)
	}

	fmt.Println("TxHash: ", hash)
}
