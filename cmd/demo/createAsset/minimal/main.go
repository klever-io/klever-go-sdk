package main

import (
	"fmt"

	"github.com/klever-io/klever-go-sdk/cmd/demo"
	"github.com/klever-io/klever-go-sdk/models"
	"github.com/klever-io/klever-go-sdk/models/proto"
)

func main() {

	accounts, wallets, kc, err := demo.InitWallets()
	if err != nil {
		panic(err)
	}

	base := accounts[0].NewBaseTX()

	tx, err := kc.CreateKDA(
		base,
		proto.KDAData_Fungible, // Type NFT ot Fungible Token
		&models.KDAOptions{
			Name:          "KleverTest",
			Ticker:        "TST",
			Precision:     4,
			MaxSupply:     0,
			InitialSupply: 0,

			Properties: models.PropertiesInfo{
				CanMint: true, CanBurn: true,
			},
			URIs: map[string]string{"explorer": "testnet.kleverscan.org"},
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
