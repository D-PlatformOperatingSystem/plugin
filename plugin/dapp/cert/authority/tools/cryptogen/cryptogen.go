// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/D-PlatformOperatingSystem/dpos/types"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/crypto/init"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/cert/authority/tools/cryptogen/generator"
	ca "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/cert/authority/tools/cryptogen/generator/impl"
	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

const (
	// CANAME   CA
	CANAME = "ca"
	// CONFIGFILENAME
	CONFIGFILENAME = "dplatformos.cryptogen.toml"
	// OUTPUTDIR
	OUTPUTDIR = "./authdir/crypto"
	// ORGNAME
	ORGNAME = "DplatformOS"
)

// Config
type Config struct {
	Name     []string
	SignType string
}

var (
	cmd = &cobra.Command{
		Use:   "cryptogen [-f configfile] [-o output directory]",
		Short: "dplatformos crypto tool for generating key and certificate",
		Run:   generate,
	}
	cfg Config
)

func initCfg(path string) *Config {
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	return &cfg
}

func main() {
	cmd.Flags().StringP("configfile", "f", CONFIGFILENAME, "config file for users")
	cmd.Flags().StringP("outputdir", "o", OUTPUTDIR, "output diraction for key and certificate")

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generate(cmd *cobra.Command, args []string) {
	configfile, _ := cmd.Flags().GetString("configfile")
	outputdir, _ := cmd.Flags().GetString("outputdir")

	initCfg(configfile)
	fmt.Println(cfg.Name)

	generateUsers(outputdir, ORGNAME)
}

func generateUsers(baseDir string, orgName string) {
	fmt.Printf("generateUsers\n")
	fmt.Println(baseDir)

	err := os.RemoveAll(baseDir)
	if err != nil {
		fmt.Printf("Clean directory %s error", baseDir)
		os.Exit(1)
	}

	caDir := filepath.Join(baseDir, "cacerts")

	signType := types.GetSignType("cert", cfg.SignType)
	if signType == types.Invalid {
		fmt.Printf("Invalid sign type:%s", cfg.SignType)
		return
	}

	signCA, err := ca.NewCA(caDir, CANAME, signType)
	if err != nil {
		fmt.Printf("Error generating signCA:%s", err.Error())
		os.Exit(1)
	}

	generateNodes(baseDir, signCA, orgName)
}

func generateNodes(baseDir string, signCA generator.CAGenerator, orgName string) {
	for _, name := range cfg.Name {
		userDir := filepath.Join(baseDir, name)
		fileName := fmt.Sprintf("%s@%s", name, orgName)
		err := signCA.GenerateLocalUser(userDir, fileName)
		if err != nil {
			fmt.Printf("Error generating local user")
			os.Exit(1)
		}
	}
}
