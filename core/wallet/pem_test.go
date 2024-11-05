package wallet_test

import (
	"encoding/hex"
	"log"
	"os"
	"testing"

	"github.com/klever-io/klever-go-sdk/core/wallet"
	"github.com/stretchr/testify/assert"
)

func tempPemFile() string {
	file, err := os.CreateTemp("", "wallet*.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString(
		`-----BEGIN PRIVATE KEY for klv1usdnywjhrlv4tcyu6stxpl6yvhplg35nepljlt4y5r7yppe8er4qujlazy-----
ODczNDA2MmMxMTU4ZjI2YTNjYThhNGEwZGE4N2I1MjdhN2MxNjg2NTNmN2Y0Yzc3
MDQ1ZTVjZjU3MTQ5N2Q5ZA==
-----END PRIVATE KEY for klv1usdnywjhrlv4tcyu6stxpl6yvhplg35nepljlt4y5r7yppe8er4qujlazy-----`)

	return file.Name()
}

func tempPemFileEncrypted() string {
	file, err := os.CreateTemp("", "wallet*.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString(
		`-----BEGIN PRIVATE KEY for klv1usdnywjhrlv4tcyu6stxpl6yvhplg35nepljlt4y5r7yppe8er4qujlazy-----
Proc-Type: 4,ENCRYPTED
DEK-Info: AES-GCM,677bac32f3e8eeb54b5f2e40

Z3usMvPo7rVLXy5AQk/z4lKS7XaarfhI4OSA8j7pL8CeBExqaApF4Op263+qFe35
YQgq5vmh9dRHYF6YdCy7Zuv2mI0OEho8KMwtqjhBXCpiILNub9qliVG140c=
-----END PRIVATE KEY for klv1usdnywjhrlv4tcyu6stxpl6yvhplg35nepljlt4y5r7yppe8er4qujlazy-----`)

	return file.Name()
}

func tempPemFileEncryptedWrong() string {
	file, err := os.CreateTemp("", "wallet*.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString(
		`-----BEGIN PRIVATE KEY for klv1usdnywjhrlv4tcyu6stxpl6yvhplg35nepljlt4y5r7yppe8er4qujlazy-----
Proc-Type: 4,ENCRYPTED
DEK-Info: AES-256,453eb5c4c21225936c8c27b2

RT61xMISJZNsjCeyEiLgO/Nftp5Nk/l1OsGg7jlbnk/YDQ6675Nq82qg3U/IIC6b
Y3osGzZtxlO0KWQ9MOeZ1aRkIDl3Mys15RmXEqBBF+Ukqmcm1K2+oupmSgw=
-----END PRIVATE KEY for klv1usdnywjhrlv4tcyu6stxpl6yvhplg35nepljlt4y5r7yppe8er4qujlazy-----`)

	return file.Name()
}

func tempEmptyPemFile() string {
	file, err := os.CreateTemp("", "wallet*.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return file.Name()
}

func tempInvalidPemFile() string {
	file, err := os.CreateTemp("", "wallet*.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString(`-----BEGIN PRIVATE KEY for klv1usdnywjhrlv4tcyu6stxpl6yvhplg35nepljlt4y5r7yppe8er4qujlazy-----`)

	return file.Name()
}

func TestLoadKey_PemFileNotFound(t *testing.T) {

	_, _, err := wallet.LoadKey("wallet.pem", 0, "")
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestLoadKey_InvalidPemFile(t *testing.T) {
	fileName := tempInvalidPemFile()

	_, _, err := wallet.LoadKey(fileName, 0, "")
	assert.Contains(t, err.Error(), "invalid pem file while reading")
}

func TestLoadKey_EmptyPemFile(t *testing.T) {
	fileName := tempEmptyPemFile()

	_, _, err := wallet.LoadKey(fileName, 0, "")
	assert.Contains(t, err.Error(), "empty file provided while")
}

func TestLoadKey_InvalidSKIndex(t *testing.T) {
	fileName := tempPemFile()

	_, _, err := wallet.LoadKey(fileName, -1, "")
	assert.Equal(t, "invalid index", err.Error())
}

func TestLoadKey_ShouldWork(t *testing.T) {
	fileName := tempPemFile()

	pk, pub, err := wallet.LoadKey(fileName, 0, "")
	assert.Nil(t, err)
	assert.Equal(t, "klv1usdnywjhrlv4tcyu6stxpl6yvhplg35nepljlt4y5r7yppe8er4qujlazy", pub)
	assert.Equal(t, "8734062c1158f26a3ca8a4a0da87b527a7c168653f7f4c77045e5cf571497d9d", hex.EncodeToString(pk))
}

func TestLoadKey_EncryptedShouldWork(t *testing.T) {
	fileName := tempPemFileEncrypted()

	pk, pub, err := wallet.LoadKey(fileName, 0, "123")
	assert.Nil(t, err)
	assert.Equal(t, "klv1usdnywjhrlv4tcyu6stxpl6yvhplg35nepljlt4y5r7yppe8er4qujlazy", pub)
	assert.Equal(t, "8734062c1158f26a3ca8a4a0da87b527a7c168653f7f4c77045e5cf571497d9d", hex.EncodeToString(pk))
}

func TestLoadKey_EncryptedWrongPassword(t *testing.T) {
	fileName := tempPemFileEncrypted()

	_, _, err := wallet.LoadKey(fileName, 0, "")
	assert.Contains(t, err.Error(), "encrypted key, must provide password")

	_, _, err = wallet.LoadKey(fileName, 0, "1")
	assert.Contains(t, err.Error(), "failed PEM decryption")
}

func TestLoadKey_EncryptedWrong(t *testing.T) {
	fileName := tempPemFileEncryptedWrong()

	_, _, err := wallet.LoadKey(fileName, 0, "22")
	assert.Contains(t, err.Error(), "invalid encryption mode")
}

func TestLoadKey_InvalidSKIndexNotFound(t *testing.T) {
	fileName := tempPemFile()

	_, _, err := wallet.LoadKey(fileName, 1, "")
	assert.Contains(t, err.Error(), "invalid index while reading")
}

func TestEncryptPEMBlock_ShouldWork(t *testing.T) {

	pemBlock, err := wallet.EncryptPEMBlock("AES-GCM", []byte("8734062c1158f26a3ca8a4a0da87b527a7c168653f7f4c77045e5cf571497d9d"), "123")
	assert.Nil(t, err)
	assert.Len(t, pemBlock.Bytes, 92)
}
