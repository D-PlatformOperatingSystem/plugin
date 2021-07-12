// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/account"
	"github.com/D-PlatformOperatingSystem/dpos/common/db"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	pt "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/trade/types"
)

/*
           token        trade     ，     symbol   token  symbol，
      symbol     exec.sybol@title, @title    ， (     ,          )。
          exec，            。


            exec = "" symbol = "TEST"
             exec =  "token"  symbol = "token.TEST"

       ,
     exec = "paracross"  symbol = "token.TEST"
     exec = "token"      symbol = "token.TEST"

*/

//GetExecSymbol : return exec, symbol
func GetExecSymbol(order *pt.SellOrder) (string, string) {
	if order.AssetExec == "" {
		return defaultAssetExec, defaultAssetExec + "." + order.TokenSymbol
	}
	return order.AssetExec, order.TokenSymbol
}

func checkAsset(cfg *types.DplatformOSConfig, height int64, exec, symbol string) bool {
	if cfg.IsDappFork(height, pt.TradeX, pt.ForkTradeAssetX) {
		if exec == "" || symbol == "" {
			return false
		}
	} else {
		if exec != "" {
			return false
		}
	}
	return true
}

func checkPrice(cfg *types.DplatformOSConfig, height int64, exec, symbol string) bool {
	if cfg.IsDappFork(height, pt.TradeX, pt.ForkTradePriceX) {
		if exec == "" && symbol != "" || exec != "" && symbol == "" {
			return false
		}
	} else {
		if exec != "" || symbol != "" {
			return false
		}
	}
	return true
}

func notSameAsset(cfg *types.DplatformOSConfig, height int64, assetExec, assetSymbol, priceExec, priceSymbol string) bool {
	if cfg.IsDappFork(height, pt.TradeX, pt.ForkTradePriceX) {
		if assetExec == priceExec && assetSymbol == priceSymbol {
			return false
		}
	}
	return true
}

func createAccountDB(cfg *types.DplatformOSConfig, height int64, db db.KV, exec, symbol string) (*account.DB, error) {
	if cfg.IsDappFork(height, pt.TradeX, pt.ForkTradeFixAssetDBX) {
		if exec == "" {
			exec = defaultAssetExec
		}
		return account.NewAccountDB(cfg, exec, symbol, db)
	} else if cfg.IsDappFork(height, pt.TradeX, pt.ForkTradeAssetX) {
		return account.NewAccountDB(cfg, exec, symbol, db)
	}

	return account.NewAccountDB(cfg, defaultAssetExec, symbol, db)
}

func createPriceDB(cfg *types.DplatformOSConfig, height int64, db db.KV, exec, symbol string) (*account.DB, error) {
	if cfg.IsDappFork(height, pt.TradeX, pt.ForkTradePriceX) {
		//        coins
		if exec == "" {
			acc := account.NewCoinsAccount(cfg)
			acc.SetDB(db)
			return acc, nil
		}
		return account.NewAccountDB(cfg, exec, symbol, db)
	}
	acc := account.NewCoinsAccount(cfg)
	acc.SetDB(db)
	return acc, nil
}
