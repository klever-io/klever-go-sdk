package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/xdg-go/pbkdf2"
)

func LoadKey(pemFile string, skIndex int, pwd string) ([]byte, string, error) {

	encodedSk, pkString, err := LoadSkPkFromPemFile(pemFile, skIndex, pwd)
	if err != nil {
		return nil, "", err
	}

	skBytes, err := hex.DecodeString(string(encodedSk))
	if err != nil {
		return nil, "", fmt.Errorf("%w for encoded secret key", err)
	}

	return skBytes, pkString, nil
}

// LoadSkPkFromPemFile loads the secret key and existing public key bytes stored in the file
func LoadSkPkFromPemFile(relativePath string, skIndex int, pwd string) ([]byte, string, error) {
	if skIndex < 0 {
		return nil, "", fmt.Errorf("invalid index")
	}

	file, err := OpenFile(relativePath)
	if err != nil {
		return nil, "", err
	}

	defer func() {
		_ = file.Close()
	}()

	buff, err := io.ReadAll(file)
	if err != nil {
		return nil, "", fmt.Errorf("%w while reading %s file", err, relativePath)
	}
	if len(buff) == 0 {
		return nil, "", fmt.Errorf("empty file provided while reading %s file", relativePath)
	}

	var blkRecovered *pem.Block

	for i := 0; i <= skIndex; i++ {
		if len(buff) == 0 {
			//less private keys present in the file than required
			return nil, "", fmt.Errorf("invalid index while reading %s file, invalid index %d", relativePath, i)
		}

		blkRecovered, buff = pem.Decode(buff)
		if blkRecovered == nil {
			return nil, "", fmt.Errorf("invalid pem file while reading %s file, error decoding", relativePath)
		}

		if IsEncryptedPEMBlock(blkRecovered) {
			if len(pwd) == 0 {
				return nil, "", errors.New("encrypted key, must provide password")
			}
			blkRecovered, err = DecryptPEMBlock(blkRecovered, pwd)
			if err != nil {
				return nil, "", fmt.Errorf("failed PEM decryption: %w", err)
			}
		}
	}

	if blkRecovered == nil {
		return nil, "", fmt.Errorf("nil pem block")
	}

	blockType := blkRecovered.Type
	header := "PRIVATE KEY for "
	if strings.Index(blockType, header) != 0 {
		return nil, "", fmt.Errorf("pem file is invalid missing '%s' in block type", header)
	}

	blockTypeString := blockType[len(header):]

	return blkRecovered.Bytes, blockTypeString, nil
}

// OpenFile method opens the file from given path - does not close the file
func OpenFile(relativePath string) (*os.File, error) {
	path, err := filepath.Abs(relativePath)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	return f, nil
}

// IsEncryptedPEMBlock returns whether the PEM block is password encrypted
// according to RFC 1423.
func IsEncryptedPEMBlock(b *pem.Block) bool {
	_, ok := b.Headers["DEK-Info"]
	return ok
}

// DecryptPEMBlock takes a PEM block encrypted according to RFC 1423 and the
// password used to encrypt it and returns a slice of decrypted DER encoded
// bytes. It inspects the DEK-Info header to determine the algorithm used for
// decryption. If no DEK-Info header is present, an error is returned. If an
// incorrect password is detected an IncorrectPasswordError is returned.
func DecryptPEMBlock(b *pem.Block, pwd string) (*pem.Block, error) {
	dek, ok := b.Headers["DEK-Info"]
	if !ok {
		return nil, errors.New("x509: no DEK-Info header in block")
	}

	mode, _, ok := strings.Cut(dek, ",")
	if !ok {
		return nil, errors.New("x509: malformed DEK-Info header")
	}

	// keep compatibility with old PEM encryption
	var password []byte
	switch mode {
	case "AES-GCM":
		password = getEncryptionKey(pwd, sha1.New)
	case "AES-256-GCM":
		password = getEncryptionKey(pwd, sha256.New)
	default:
		return nil, errors.New("invalid encryption mode")
	}

	block, err := aes.NewCipher(password)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if b.Bytes == nil || len(b.Bytes) < nonceSize {
		return nil, fmt.Errorf("invalid data size")
	}
	nonce, ciphertext := b.Bytes[:nonceSize], b.Bytes[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return &pem.Block{
		Type:  b.Type,
		Bytes: plaintext,
	}, nil

}

// EncryptPEMBlock returns a PEM block of the specified type holding the
// given DER encoded data encrypted with GCM algorithm and
// password according to RFC 1423.
func EncryptPEMBlock(blockType string, data []byte, pwd string) (*pem.Block, error) {
	// get encryption key using password and sha256
	password := getEncryptionKey(pwd, sha256.New)

	block, _ := aes.NewCipher(password)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	// append nonce to seal package
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return &pem.Block{
		Type: blockType,
		Headers: map[string]string{
			"Proc-Type": "4,ENCRYPTED",
			"DEK-Info":  "AES-256-GCM," + hex.EncodeToString(nonce),
		},
		Bytes: ciphertext,
	}, nil
}

// getEncryptionKey based for pin
func getEncryptionKey(pin string, hashType func() hash.Hash) []byte {
	idBytes := StringHash("kleverchain")
	pwdBytes := StringHash(pin)

	return PBKDFPass(pwdBytes, idBytes, hashType)
}

// StringHash computes SHA256 from string
func StringHash(text string) []byte {
	h := sha256.New()
	h.Write([]byte(text))
	return h.Sum(nil)
}

// PBKDFPass from password and salt
func PBKDFPass(password, salt []byte, hashType func() hash.Hash) []byte {
	return pbkdf2.Key(password, salt, 4096, 32, hashType)
}
