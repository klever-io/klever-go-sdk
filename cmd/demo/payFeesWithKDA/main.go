package payfeeswithkda

import (
	"fmt"

	"github.com/klever-io/klever-go-sdk/cmd/demo"
)



func main(){

	// To pay fees with kda, you need to have a KDA with Fee Pool Active.
	accounts, wallets, kc, err := demo.InitWallets()
	if err != nil {
		panic(err)
	}

	base := accounts[0].NewBaseTX()
	// and add the current kda to kdaFee param on base struct.
	base.KdaFee = "KDA"
	tx, err := kc.Send(base, accounts[1].Address().Bech32(), 1, "")
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