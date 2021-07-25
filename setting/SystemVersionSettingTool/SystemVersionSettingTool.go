package SystemVersionSettingTool

import "helloworld-blockchain-go/setting/SystemVersionSetting"

func CheckSystemVersion(blockHeight uint64) bool {
	return blockHeight <= SystemVersionSetting.BLOCK_CHAIN_VERSION_LIST[len(SystemVersionSetting.BLOCK_CHAIN_VERSION_LIST)-1]
}
