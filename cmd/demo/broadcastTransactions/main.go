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
	tx, err := kc.Send(base, accounts[1].Address().Bech32(), 1, "")
	if err != nil {
		panic(err)
	}

	err = tx.Sign(wallets[0])
	if err != nil {
		panic(err)
	}

	//increase nonce
	base.Nonce++

	tx2, err := kc.Send(base, accounts[1].Address().Bech32(), 1, "")
	if err != nil {
		panic(err)
	}

	err = tx2.Sign(wallets[0])
	if err != nil {
		panic(err)
	}

	hashes, err := kc.BroadcastTransactions(tx, tx2)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n\n\nHashes: ", hashes)
}
