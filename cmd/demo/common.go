package demo

import (
	"os"
	"path/filepath"
	"time"

	"github.com/klever-io/klever-go-sdk/core/account"
	"github.com/klever-io/klever-go-sdk/core/wallet"
	"github.com/klever-io/klever-go-sdk/provider"
	"github.com/klever-io/klever-go-sdk/provider/network"
	"github.com/klever-io/klever-go-sdk/provider/utils"
)

func InitWallets() (
	accounts []account.Account,
	wallets []wallet.Wallet,
	kc provider.KleverChain,
	err error,
) {

	wallets = make([]wallet.Wallet, 0)

	net := network.NewNetworkConfig(network.TestNet)
	httpClient := utils.NewHttpClient(10 * time.Second)
	kc, err = provider.NewKleverChain(net, httpClient)
	if err != nil {
		return
	}

	// load account from pemfile
	files, err := GetWallets(".", "*.pem")
	if err != nil {
		return
	}
	for _, f := range files {
		var w wallet.Wallet
		w, err = wallet.NewWalletFromPEM(f)
		if err != nil {
			return
		}
		wallets = append(wallets, w)

		var acc account.Account
		acc, err = w.GetAccount()
		if err != nil {
			return
		}
		err = acc.Sync(kc)

		accounts = append(accounts, acc)
	}

	return
}

func GetWallets(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
