package vault

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/scrypt"
)

func (s *HeaderTestSuite) TestVaultInternalGenerateKey() {
	pass := []byte("secret pass")
	salt := make([]byte, 32)

	result, err := scrypt.Key(pass, salt, keyCost, keyR, keyP, keyLen)
	if assert.Nil(s.T(), err) {
		assert.NotNil(s.T(), result)
	}
}
