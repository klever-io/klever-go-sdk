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
	tx, err := kc.MultiTransfer(
		base,
		[]models.ToAmount{{
			ToAddress: accounts[1].Address().Bech32(),
			Amount:    1000,
			KDA:       "KLV",
		}})
	if err != nil {
		panic(err)
	}

	err = tx.Sign(wallets[0])
	if err != nil {
		panic(err)
	}

	decodedTx, err := kc.Decode(tx)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n\n\nEncodedTX: ", tx)
	fmt.Println("\n\n\nDecodedTX: ", decodedTx)
}
