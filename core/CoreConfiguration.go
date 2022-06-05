package core

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/util/ByteUtil"
	"helloworldcoin-go/util/FileUtil"
	"helloworldcoin-go/util/KvDbUtil"
)

type CoreConfiguration struct {
	corePath string
}

func NewCoreConfiguration(corePath string) *CoreConfiguration {
	var coreConfiguration CoreConfiguration
	coreConfiguration.corePath = corePath
	return &coreConfiguration
}

const CONFIGURATION_DATABASE_NAME = "ConfigurationDatabase"

const MINE_OPTION_KEY = "IS_MINER_ACTIVE"

const MINER_MINE_MAX_BLOCK_HEIGHT_KEY = "MINER_MINE_MAX_BLOCK_HEIGHT"

const MINE_OPTION_DEFAULT_VALUE = false

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

func (c *CoreConfiguration) setMinerMineMaxBlockHeight(maxHeight uint64) {
	c.addOrUpdateConfiguration(ByteUtil.StringToUtf8Bytes(MINER_MINE_MAX_BLOCK_HEIGHT_KEY), ByteUtil.Uint64ToBytes(maxHeight))
}

func (c *CoreConfiguration) getMinerMineMaxBlockHeight() uint64 {
	bytesMineMaxBlockHeight := c.getConfigurationValue(ByteUtil.StringToUtf8Bytes(MINER_MINE_MAX_BLOCK_HEIGHT_KEY))
	if bytesMineMaxBlockHeight == nil {
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
