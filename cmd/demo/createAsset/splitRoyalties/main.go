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

	splitRoyalties := map[string]*models.RoyaltySplitInfo{
		accounts[0].Address().Bech32(): {
			PercentITOPercentage: 50,
			PercentITOFixed:      50,
		},
		accounts[1].Address().Bech32(): {
			PercentITOPercentage: 50,
			PercentITOFixed:      50,
		},
	}

	tx, err := kc.CreateKDA(
		base,
		proto.KDAData_Fungible, // Type NFT ot Fungible Token
		&models.KDAOptions{
			Name:          "KleverTest",
			Ticker:        "TST",
			Precision:     4,
			MaxSupply:     1000,
			InitialSupply: 10,
			AddRolesMint:  []string{accounts[0].Address().Bech32(), accounts[1].Address().Bech32()},
			Properties: models.PropertiesInfo{
				CanMint: true, CanBurn: true,
			},
			Royalties: models.RoyaltiesInfo{
				ITOPercentage:  10,
				ITOFixed:       10,
				SplitRoyalties: splitRoyalties,
			},
			URIs:         map[string]string{"explorer": "testnet.kleverscan.org"},
			AdminAddress: base.FromAddress,
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
