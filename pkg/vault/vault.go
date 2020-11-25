package vault

import (
	"crypto/aes"
	"crypto/cipher"

	"golang.org/x/crypto/scrypt"
)

const (
	blockSize int = 2048
	lenSize   int = 2
	macSize   int = 32

	keyCost = 1048576
	keyLen  = 32
	keyP    = 1
	keyR    = 8
)

func createStream(key, iv []byte) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewOFB(block, iv)
	return stream, nil
}

func generateKey(pass, salt []byte) []byte {
	key, _ := scrypt.Key(pass, salt, keyCost, keyR, keyP, keyLen)
	return key
}
