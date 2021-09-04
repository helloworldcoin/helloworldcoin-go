package core

/*
 @author king 409060350@qq.com
*/

import (
	"helloworld-blockchain-go/crypto/ByteUtil"
	"helloworld-blockchain-go/util/FileUtil"
	"helloworld-blockchain-go/util/KvDbUtil"
)

type CoreConfiguration struct {
	corePath string
}

//配置数据库名字
const CONFIGURATION_DATABASE_NAME = "ConfigurationDatabase"

//'矿工是否是激活状态'存入到数据库时的主键
const MINE_OPTION_KEY = "IS_MINER_ACTIVE"

//'矿工可挖的最高区块高度'存入到数据库时的主键
const MINE_MAX_BLOCK_HEIGHT_KEY = "MAX_BLOCK_HEIGHT"

//'矿工是否是激活状态'的默认值
const MINE_OPTION_DEFAULT_VALUE = false

//这个时间间隔更新一次正在被挖矿的区块的交易。如果时间太长，可能导致新提交的交易延迟被确认。
const MINE_TIMESTAMP_PER_ROUND = uint64(1000) * uint64(10)

func (c *CoreConfiguration) getCorePath() string {
	return c.corePath
}

func (c *CoreConfiguration) IsMinerActive() bool {
	mineOption := c.getConfigurationValue(ByteUtil.StringToUtf8Bytes(MINE_OPTION_KEY))
	if mineOption == nil {
		return MINE_OPTION_DEFAULT_VALUE
	}
	return ByteUtil.Utf8BytesToBoolean(mineOption)
}

func (c *CoreConfiguration) activeMiner() {
	c.addOrUpdateConfiguration(ByteUtil.StringToUtf8Bytes(MINE_OPTION_KEY), ByteUtil.BooleanToUtf8Bytes(true))
}

func (c *CoreConfiguration) deactiveMiner() {
	c.addOrUpdateConfiguration(ByteUtil.StringToUtf8Bytes(MINE_OPTION_KEY), ByteUtil.BooleanToUtf8Bytes(false))
}

func (c *CoreConfiguration) setMaxBlockHeight(maxHeight uint64) {
	c.addOrUpdateConfiguration(ByteUtil.StringToUtf8Bytes(MINE_MAX_BLOCK_HEIGHT_KEY), ByteUtil.Uint64ToBytes(maxHeight))
}

func (c *CoreConfiguration) getMaxBlockHeight() uint64 {
	bytesMineMaxBlockHeight := c.getConfigurationValue(ByteUtil.StringToUtf8Bytes(MINE_MAX_BLOCK_HEIGHT_KEY))
	if bytesMineMaxBlockHeight == nil {
		//设置默认值，这是一个十分巨大的数字，矿工永远挖不到的高度
		return uint64(10000000000000000)
	}
	return ByteUtil.BytesToUint64(bytesMineMaxBlockHeight)
}

func (c *CoreConfiguration) GetMinerMineTimeInterval() uint64 {
	return MINE_TIMESTAMP_PER_ROUND
}

func (c *CoreConfiguration) getConfigurationDatabasePath() string {
	return FileUtil.NewPath(c.corePath, CONFIGURATION_DATABASE_NAME)
}
func (c *CoreConfiguration) getConfigurationValue(configurationKey []byte) []byte {
	return KvDbUtil.Get(c.getConfigurationDatabasePath(), configurationKey)
}
func (c *CoreConfiguration) addOrUpdateConfiguration(configurationKey []byte, configurationValue []byte) {
	KvDbUtil.Put(c.getConfigurationDatabasePath(), configurationKey, configurationValue)
}
