package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

//AES ...
type AES struct {
	key []byte
	//iv       block    ，   16  ，
	iv []byte
}

//NewAES       16,24,32   ，
func NewAES(key, iv []byte) *AES {
	return &AES{key: key, iv: iv}
}

//Encrypt ...
func (a *AES) Encrypt(origData []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, a.iv[:blockSize])
	crypted := make([]byte, len(origData))
	//   CryptBlocks     ，       crypted
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//Decrypt ...
func (a *AES) Decrypt(crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, a.iv[:blockSize])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}
