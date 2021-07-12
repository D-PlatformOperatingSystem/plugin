package crypto

import (
	"crypto/cipher"
	"crypto/des"
)

//DES ...
type DES struct {
	key []byte
	//iv       block
	iv []byte
}

//NewDES ...
func NewDES(key, iv []byte) *DES {
	return &DES{key: key, iv: iv}
}

//Encrypt ...
func (d *DES) Encrypt(origData []byte) ([]byte, error) {
	block, err := des.NewCipher(d.key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, d.iv[:block.BlockSize()])
	crypted := make([]byte, len(origData))
	//   CryptBlocks     ，       crypted
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//Decrypt   key    8
func (d *DES) Decrypt(crypted []byte) ([]byte, error) {
	block, err := des.NewCipher(d.key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, d.iv[:block.BlockSize()])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

//TripleDES ...
type TripleDES struct {
	key []byte
	//iv       block
	iv []byte
}

//NewTripleDES ...
func NewTripleDES(key, iv []byte) *TripleDES {
	return &TripleDES{key: key, iv: iv}
}

//Encrypt 3DES   24
func (d *TripleDES) Encrypt(origData []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(d.key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, d.iv[:block.BlockSize()])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//Decrypt 3DES
func (d *TripleDES) Decrypt(crypted []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(d.key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, d.iv[:block.BlockSize()])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}
