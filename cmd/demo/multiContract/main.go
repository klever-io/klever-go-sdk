package main

import (
	"fmt"
	"github.com/klever-io/klever-go-sdk/models"
	"github.com/klever-io/klever-go-sdk/models/proto"

	"github.com/klever-io/klever-go-sdk/cmd/demo"
)

func main() {
	accounts, wallets, kc, err := demo.InitWallets()
	if err != nil {
		panic(err)
	}

	base := accounts[0].NewBaseTX()
	sendContract := models.AnyContractRequest{
		ContractType: uint32(proto.TXContract_TransferContractType),
		Contract: models.TransferTXRequest{
			Receiver: accounts[1].Address().Bech32(),
			Amount:   1,
			KDA:      "KLV",
		},
	}

	freezeContract := models.AnyContractRequest{
		ContractType: uint32(proto.TXContract_FreezeContractType),
		Contract: models.FreezeTXRequest{
			Amount: 1000_000_000,
			KDA:    "KLV",
		},
	}

	tx, err := kc.MultiSend(base, sendContract, freezeContract)
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
