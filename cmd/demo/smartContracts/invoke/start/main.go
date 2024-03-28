package main

import (
	"fmt"
	"time"

	"github.com/klever-io/klever-go-sdk/cmd/demo"
)

func main() {
	accounts, wallets, kc, err := demo.InitWallets()

	if err != nil {
		panic(err)
	}

	base := accounts[0].NewBaseTX()

	timeStamp := time.Now().Unix()
	lotteryDuration := int64(120) // 120 seconds -> 2 minutes
	lotteryDeadline := fmt.Sprintf("optionu64:%d", timeStamp+lotteryDuration)
	tx, err := kc.InvokeSmartContract(
		base,
		"klv1qqqqqqqqqqqqqpgq4f6qkvv34kr50w5lucdn6z0238r7m96nhtxspjw22t",
		"createLotteryPool",
		map[string]int64{},
		"demo", "KLV", "N:10000000", "E:0", lotteryDeadline, "E:0", "E:0", "E:0",
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
