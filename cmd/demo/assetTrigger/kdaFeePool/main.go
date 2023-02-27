package main

import (
	"fmt"

	"github.com/klever-io/klever-go-sdk/cmd/demo"
	"github.com/klever-io/klever-go-sdk/models"
	"github.com/klever-io/klever-go-sdk/provider"
)

func main() {

	accounts, wallets, kc, err := demo.InitWallets()
	if err != nil {
		panic(err)
	}

	base := accounts[0].NewBaseTX()
	tx, err := kc.AssetTrigger(
		base,
		"KDA",
		provider.AssetTriggerType(15),
		&models.AssetTriggerOptions{
			KDAPool: &models.KDAPoolInfo{
				Active:       true,
				AdminAddress: accounts[0].Address().Bech32(),
				FRatioKLV:    1,
				FRatioKDA:    1,
			},
		})
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
