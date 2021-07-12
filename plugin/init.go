package plugin

import (
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/consensus/init" //consensus init
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/crypto/init"    //crypto init
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/init"      //dapp init
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/mempool/init"   //mempool init
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/p2p/init"       //p2p init
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/store/init"     //store init
)
