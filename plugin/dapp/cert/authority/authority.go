// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authority

import (
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"sync"

	"bytes"

	"github.com/D-PlatformOperatingSystem/dpos/common/crypto"
	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/cert/authority/core"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/cert/authority/utils"
	ty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/cert/types"
)

var (
	alog   = log.New("module", "authority")
	cpuNum = runtime.NumCPU()

	// OrgName
	OrgName = "DplatformOS"

	// Author
	Author = &Authority{}

	// IsAuthEnable
	IsAuthEnable = false
)

// Authority
type Authority struct {
	//
	cryptoPath string
	// certByte
	authConfig *core.AuthConfig
	//
	validator core.Validator
	//
	signType int
	//
	validCertCache [][]byte
	//
	HistoryCertCache *HistoryCertData
}

// HistoryCertData
type HistoryCertData struct {
	CryptoCfg *core.AuthConfig
	CurHeight int64
	NxtHeight int64
}

// Init    auth
func (auth *Authority) Init(conf *ty.Authority) error {
	if conf == nil || !conf.Enable {
		return nil
	}

	if len(conf.CryptoPath) == 0 {
		alog.Error("Crypto config path can not be null")
		return types.ErrInvalidParam
	}
	auth.cryptoPath = conf.CryptoPath

	sign := types.GetSignType("cert", conf.SignType)
	if sign == types.Invalid {
		alog.Error(fmt.Sprintf("Invalid sign type:%s", conf.SignType))
		return types.ErrInvalidParam
	}
	auth.signType = sign

	authConfig, err := core.GetAuthConfig(conf.CryptoPath)
	if err != nil {
		alog.Error("Get authority crypto config failed")
		return err
	}
	auth.authConfig = authConfig

	vldt, err := core.GetLocalValidator(authConfig, auth.signType)
	if err != nil {
		alog.Error(fmt.Sprintf("Get loacal validator failed. err:%s", err.Error()))
		return err
	}
	auth.validator = vldt

	auth.validCertCache = make([][]byte, 0)
	auth.HistoryCertCache = &HistoryCertData{authConfig, -1, -1}

	IsAuthEnable = true
	return nil
}

// newAuthConfig store    authConfig
func newAuthConfig(store *types.HistoryCertStore) *core.AuthConfig {
	ret := &core.AuthConfig{}
	ret.RootCerts = make([][]byte, len(store.Rootcerts))
	for i, v := range store.Rootcerts {
		ret.RootCerts[i] = append(ret.RootCerts[i], v...)
	}

	ret.IntermediateCerts = make([][]byte, len(store.IntermediateCerts))
	for i, v := range store.IntermediateCerts {
		ret.IntermediateCerts[i] = append(ret.IntermediateCerts[i], v...)
	}

	ret.RevocationList = make([][]byte, len(store.RevocationList))
	for i, v := range store.RevocationList {
		ret.RevocationList[i] = append(ret.RevocationList[i], v...)
	}

	return ret
}

// ReloadCert               ，
func (auth *Authority) ReloadCert(store *types.HistoryCertStore) error {
	if !IsAuthEnable {
		return nil
	}

	//
	if len(store.Rootcerts) == 0 {
		auth.authConfig = nil
		auth.validator, _ = core.NewNoneValidator()
	} else {
		auth.authConfig = newAuthConfig(store)
		//
		vldt, err := core.GetLocalValidator(auth.authConfig, auth.signType)
		if err != nil {
			return err
		}
		auth.validator = vldt
	}

	//
	auth.validCertCache = auth.validCertCache[:0]

	//
	auth.HistoryCertCache = &HistoryCertData{auth.authConfig, store.CurHeigth, store.NxtHeight}

	return nil
}

// ReloadCertByHeght    authdir        ，
func (auth *Authority) ReloadCertByHeght(currentHeight int64) error {
	if !IsAuthEnable {
		return nil
	}

	authConfig, err := core.GetAuthConfig(auth.cryptoPath)
	if err != nil {
		alog.Error("Get authority crypto config failed")
		return err
	}
	auth.authConfig = authConfig

	//
	vldt, err := core.GetLocalValidator(auth.authConfig, auth.signType)
	if err != nil {
		return err
	}
	auth.validator = vldt

	//
	auth.validCertCache = auth.validCertCache[:0]

	//
	auth.HistoryCertCache = &HistoryCertData{auth.authConfig, currentHeight, -1}

	return nil
}

// ValidateCerts
func (auth *Authority) ValidateCerts(task []*types.Signature) bool {
	//FIXME               ，
	done := make(chan struct{})
	defer close(done)

	taskes := gen(done, task)

	c := make(chan result)
	var wg sync.WaitGroup
	wg.Add(cpuNum)
	for i := 0; i < cpuNum; i++ {
		go func() {
			auth.task(done, taskes, c)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(c)
	}()

	for r := range c {
		if r.err != nil {
			return false
		}
	}

	return true
}

func gen(done <-chan struct{}, task []*types.Signature) <-chan *types.Signature {
	ch := make(chan *types.Signature)
	go func() {
		defer func() {
			close(ch)
		}()
		for i := 0; i < len(task); i++ {
			select {
			case ch <- task[i]:
			case <-done:
				return
			}
		}
	}()
	return ch
}

type result struct {
	err error
}

func (auth *Authority) task(done <-chan struct{}, taskes <-chan *types.Signature, c chan<- result) {
	for task := range taskes {
		select {
		case c <- result{auth.Validate(task)}:
		case <-done:
			return
		}
	}
}

// Validate
func (auth *Authority) Validate(signature *types.Signature) error {
	//  proto   signature
	cert, err := auth.validator.GetCertFromSignature(signature.Signature)
	if err != nil {
		return err
	}

	//
	for _, v := range auth.validCertCache {
		if bytes.Equal(v, cert) {
			return nil
		}
	}

	//
	err = auth.validator.Validate(cert, signature.GetPubkey())
	if err != nil {
		alog.Error(fmt.Sprintf("validate cert failed. %s", err.Error()))
		return fmt.Errorf("validate cert failed. error:%s", err.Error())
	}
	auth.validCertCache = append(auth.validCertCache, cert)

	return nil
}

// GetSnFromSig
func (auth *Authority) GetSnFromByte(signature *types.Signature) ([]byte, error) {
	return auth.validator.GetCertSnFromSignature(signature.Signature)

}

// ToHistoryCertStore       store
func (certdata *HistoryCertData) ToHistoryCertStore(store *types.HistoryCertStore) {
	if store == nil {
		alog.Error("Convert cert data to cert store failed")
		return
	}

	store.Rootcerts = make([][]byte, len(certdata.CryptoCfg.RootCerts))
	for i, v := range certdata.CryptoCfg.RootCerts {
		store.Rootcerts[i] = append(store.Rootcerts[i], v...)
	}

	store.IntermediateCerts = make([][]byte, len(certdata.CryptoCfg.IntermediateCerts))
	for i, v := range certdata.CryptoCfg.IntermediateCerts {
		store.IntermediateCerts[i] = append(store.IntermediateCerts[i], v...)
	}

	store.RevocationList = make([][]byte, len(certdata.CryptoCfg.RevocationList))
	for i, v := range certdata.CryptoCfg.RevocationList {
		store.RevocationList[i] = append(store.RevocationList[i], v...)
	}

	store.CurHeigth = certdata.CurHeight
	store.NxtHeight = certdata.NxtHeight
}

// User
type User struct {
	ID   string
	Cert []byte
	Key  crypto.PrivKey
}

// UserLoader SKD  user
type UserLoader struct {
	configPath string
	userMap    map[string]*User
	signType   int
}

// Init userloader
func (loader *UserLoader) Init(configPath string, signType string) error {
	loader.configPath = configPath
	loader.userMap = make(map[string]*User)

	sign := types.GetSignType("cert", signType)
	if sign == types.Invalid {
		alog.Error(fmt.Sprintf("Invalid sign type:%s", signType))
		return types.ErrInvalidParam
	}
	loader.signType = sign

	return loader.loadUsers()
}

func (loader *UserLoader) loadUsers() error {
	certDir := path.Join(loader.configPath, "signcerts")
	dir, err := ioutil.ReadDir(certDir)
	if err != nil {
		return err
	}

	keyDir := path.Join(loader.configPath, "keystore")
	for _, file := range dir {
		filePath := path.Join(certDir, file.Name())
		certBytes, err := utils.ReadFile(filePath)
		if err != nil {
			continue
		}

		ski, err := utils.GetPublicKeySKIFromCert(certBytes, loader.signType)
		if err != nil {
			alog.Error(err.Error())
			continue
		}
		filePath = path.Join(keyDir, ski+"_sk")
		keyBytes, err := utils.ReadFile(filePath)
		if err != nil {
			continue
		}

		priv, err := loader.genCryptoPriv(keyBytes)
		if err != nil {
			alog.Error(fmt.Sprintf("Generate crypto private failed. error:%s", err.Error()))
			continue
		}

		loader.userMap[file.Name()] = &User{file.Name(), certBytes, priv}
	}

	return nil
}

func (loader *UserLoader) genCryptoPriv(keyBytes []byte) (crypto.PrivKey, error) {
	cr, err := crypto.New(types.GetSignName("cert", loader.signType))
	if err != nil {
		return nil, fmt.Errorf("create crypto %s failed, error:%s", types.GetSignName("cert", loader.signType), err)
	}
	privKeyByte, err := utils.PrivKeyByteFromRaw(keyBytes, loader.signType)
	if err != nil {
		return nil, err
	}

	priv, err := cr.PrivKeyFromBytes(privKeyByte)
	if err != nil {
		return nil, fmt.Errorf("get private key failed, error:%s", err)
	}

	return priv, nil
}

// Get        user
func (loader *UserLoader) Get(userName string) (*User, error) {
	keyvalue := fmt.Sprintf("%s@%s-cert.pem", userName, OrgName)
	user, ok := loader.userMap[keyvalue]
	if !ok {
		return nil, types.ErrInvalidParam
	}

	resp := &User{}
	resp.Cert = append(resp.Cert, user.Cert...)
	resp.Key = user.Key

	return resp, nil
}
