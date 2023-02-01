package main

import (
	"fmt"

	"github.com/klever-io/klever-go-sdk/cmd/demo"
	"github.com/klever-io/klever-go-sdk/models"
)

func main() {

	accounts, wallets, kc, err := demo.InitWallets()
	if err != nil {
		panic(err)
	}

	base := accounts[0].NewBaseTX()
	tx, err := kc.Deposit(
		base,
		&models.DepositOptions{
			Amount:      100,
			DepositType: 1,
			KDAID:       "KDA",
			CurrencyID:  "KLV",
		},
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
