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

	tx, err := kc.DeploySmartContract(
		base,
		"./cmd/demo/smartContracts/scFiles/lottery-kda.wasm", // `.` is the root of that project
		true, true, true, true, "",
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
