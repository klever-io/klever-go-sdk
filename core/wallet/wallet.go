package wallet

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/go-bip39"
	"github.com/klever-io/klever-go-sdk/core/account"
	"github.com/klever-io/klever-go-sdk/core/address"
)

const HDPrefix = "m/44'/%d'/%d'/0'/%d'"

type wallet struct {
	privateKey ed25519.PrivateKey
	publicKey  []byte
}

func NewWallet(privateKey []byte) (Wallet, error) {
	if len(privateKey) != ed25519.SeedSize {
		return nil, fmt.Errorf("invalid private key size")
	}

	// derive pubKey
	w := &wallet{
		privateKey: ed25519.NewKeyFromSeed(privateKey),
	}

	w.publicKey = w.privateKey.Public().(ed25519.PublicKey)

	return w, nil
}

func NewWalletFromMnemonic(mnemonic string, option ...WOHDPath) (Wallet, error) {
	if len(option) == 0 {
		option = []WOHDPath{{690, 0}}
	}
	if len(option) != 1 {
		return nil, fmt.Errorf("invalid options")
	}

	path := strings.Replace(HDPrefix, "m/", "", 1)
	path = fmt.Sprintf(path, option[0].Prefix, 0, option[0].Index)

	private, err := deriveFromPath(mnemonic, path, "")
	if err != nil {
		return nil, err
	}

	return NewWallet(private[:])
}

type Path struct {
	n        uint32
	hardered bool
}

const (
	HARDENED     = uint32(0x80000000)
	ED25519_SEED = "ed25519 seed"
)

func deriveFromPath(mnemonic, path, password string) ([]byte, error) {
	seed := bip39.NewSeed(mnemonic, password)
	path = strings.Replace(path, "m/", "", 1)
	pathSplit := strings.Split(path, "/")

	pathNumbers := make([]Path, 0)
	for _, str := range pathSplit {
		n, err := strconv.ParseUint(str[:len(str)-1], 10, 32)
		if err != nil {
			n, err = strconv.ParseUint(str, 10, 32)
			if err != nil {
				return nil, err
			}
		}

		obj := Path{
			n:        uint32(n),
			hardered: str[len(str)-1:] == "'",
		}

		pathNumbers = append(pathNumbers, obj)
	}

	coinType := pathNumbers[len(pathNumbers)-4]
	account := pathNumbers[len(pathNumbers)-3]
	addressIndex := pathNumbers[len(pathNumbers)-1]
	bip := pathNumbers[0]

	pathB := []uint32{
		bip.n | HARDENED,
		coinType.n | HARDENED,
		account.n | HARDENED, // account
		HARDENED,
		addressIndex.n | HARDENED, // addressIndex
	}

	if !bip.hardered {
		pathB[0] = bip.n
	}
	if !account.hardered {
		pathB[2] = account.n
	}
	if !coinType.hardered {
		pathB[3] = coinType.n
	}
	if !addressIndex.hardered {
		pathB[4] = addressIndex.n
	}

	key := make([]byte, 0)
	digest := hmac.New(sha512.New, []byte(ED25519_SEED))
	digest.Write(seed)
	intermediary := digest.Sum(nil)
	serializedKeyLen := 32
	serializedChildIndexLen := 4
	hardenedChildPadding := byte(0x00)
	key = intermediary[:serializedKeyLen]
	chainCode := intermediary[serializedKeyLen:]
	for _, childIdx := range pathB {
		data := make([]byte, 1+serializedKeyLen+4)
		data[0] = hardenedChildPadding
		copy(data[1:1+serializedKeyLen], key)
		binary.BigEndian.PutUint32(data[1+serializedKeyLen:1+serializedKeyLen+serializedChildIndexLen], childIdx)
		digest = hmac.New(sha512.New, chainCode)
		digest.Write(data)
		intermediary = digest.Sum(nil)
		key = intermediary[:serializedKeyLen]
		chainCode = intermediary[serializedKeyLen:]
	}

	return key, nil
}

func NewWalletFroHex(privateHex string) (Wallet, error) {
	pk, err := hex.DecodeString(privateHex)
	if err != nil {
		return nil, err
	}

	return NewWallet(pk)
}

func (w *wallet) PrivateKey() []byte {
	return append([]byte{}, w.privateKey...)
}

func (w *wallet) PublicKey() []byte {
	return append([]byte{}, w.publicKey...)
}

func (w *wallet) GetAccount() (account.Account, error) {
	addr, err := address.NewAddressFromBytes(w.publicKey)
	if err != nil {
		return nil, err
	}
	return account.NewAccount(addr)
}
