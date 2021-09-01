package configuration

import (
	"helloworld-blockchain-go/crypto/ByteUtil"
	"helloworld-blockchain-go/util/FileUtil"
	"helloworld-blockchain-go/util/KvDbUtil"
)

type NetCoreConfiguration struct {
	NetCorePath string
}

const NETCORE_CONFIGURATION_DATABASE_NAME = "NetCoreConfigurationDatabase"

//节点搜索器"是否是自动搜索新区块"状态存入到数据库时的主键
const AUTO_SEARCH_BLOCK_OPTION_KEY = "IS_AUTO_SEARCH_BLOCK"

//节点搜索器"是否是自动搜索新区块"开关的默认状态
const AUTO_SEARCH_BLOCK_OPTION_DEFAULT_VALUE = true

//节点搜索器'是否自动搜索节点'状态存入到数据库时的主键
const AUTO_SEARCH_NODE_OPTION_KEY = "IS_AUTO_SEARCH_NODE"

//节点搜索器'是否自动搜索节点'开关的默认状态
const AUTO_SEARCH_NODE_OPTION_DEFAULT_VALUE = true

//在区块链网络中自动搜寻新的节点的间隔时间
const SEARCH_NODE_TIME_INTERVAL = 1000 * 60 * 2

//在区块链网络中自动搜索节点的区块链高度
const SEARCH_BLOCKCHAIN_HEIGHT_TIME_INTERVAL = 1000 * 60 * 2

//在区块链网络中自动搜寻新的区块的间隔时间。
const SEARCH_BLOCKS_TIME_INTERVAL = 1000 * 60 * 2

//区块高度广播时间间隔
const BLOCKCHAIN_HEIGHT_BROADCASTER_TIME_INTERVAL = 1000 * 20

//区块广播时间间隔。
const BLOCK_BROADCASTER_TIME_INTERVAL = 1000 * 20

//定时将种子节点加入本地区块链网络的时间间隔。
const ADD_SEED_NODE_TIME_INTERVAL = 1000 * 60 * 2

//广播自己节点的时间间隔。
const NODE_BROADCAST_TIME_INTERVAL = 1000 * 60 * 2

//清理死亡节点的时间间隔。
const NODE_CLEAN_TIME_INTERVAL = 1000 * 60 * 10

//两个区块链有分叉时，区块差异数量大于这个值，则硬分叉了。
const HARD_FORK_BLOCK_COUNT = 100000000

//在区块链网络中搜寻未确认交易的间隔时间。
const SEARCH_UNCONFIRMED_TRANSACTIONS_TIME_INTERVAL = 1000 * 60 * 2

func NewNetCoreConfiguration(netCorePath string) NetCoreConfiguration {
	FileUtil.MakeDirectory(netCorePath)
	return NetCoreConfiguration{NetCorePath: netCorePath}
}

func (n NetCoreConfiguration) getNetCorePath() string {
	return n.NetCorePath
}

func (n NetCoreConfiguration) IsAutoSearchBlock() bool {
	bytesConfigurationValue := n.getConfigurationValue(ByteUtil.StringToUtf8Bytes(AUTO_SEARCH_BLOCK_OPTION_KEY))
	if bytesConfigurationValue == nil {
		return AUTO_SEARCH_BLOCK_OPTION_DEFAULT_VALUE
	}
	return ByteUtil.Utf8BytesToBoolean(bytesConfigurationValue)
}

func (n NetCoreConfiguration) ActiveAutoSearchBlock() {
	n.addOrUpdateConfiguration(ByteUtil.StringToUtf8Bytes(AUTO_SEARCH_BLOCK_OPTION_KEY), ByteUtil.BooleanToUtf8Bytes(true))
}

func (n NetCoreConfiguration) DeactiveAutoSearchBlock() {
	n.addOrUpdateConfiguration(ByteUtil.StringToUtf8Bytes(AUTO_SEARCH_BLOCK_OPTION_KEY), ByteUtil.BooleanToUtf8Bytes(false))
}

func (n NetCoreConfiguration) IsAutoSearchNode() bool {
	bytesConfigurationValue := n.getConfigurationValue(ByteUtil.StringToUtf8Bytes(AUTO_SEARCH_NODE_OPTION_KEY))
	if bytesConfigurationValue == nil {
		return AUTO_SEARCH_NODE_OPTION_DEFAULT_VALUE
	}
	return ByteUtil.Utf8BytesToBoolean(bytesConfigurationValue)
}

func (n NetCoreConfiguration) ActiveAutoSearchNode() {
	n.addOrUpdateConfiguration(ByteUtil.StringToUtf8Bytes(AUTO_SEARCH_NODE_OPTION_KEY), ByteUtil.BooleanToUtf8Bytes(true))
}

func (n NetCoreConfiguration) DeactiveAutoSearchNode() {
	n.addOrUpdateConfiguration(ByteUtil.StringToUtf8Bytes(AUTO_SEARCH_NODE_OPTION_KEY), ByteUtil.BooleanToUtf8Bytes(false))
}

func (n NetCoreConfiguration) GetSearchNodeTimeInterval() uint64 {
	return SEARCH_NODE_TIME_INTERVAL
}

func (n NetCoreConfiguration) GetSearchBlockchainHeightTimeInterval() uint64 {
	return SEARCH_BLOCKCHAIN_HEIGHT_TIME_INTERVAL
}

func (n NetCoreConfiguration) GetSearchBlockTimeInterval() uint64 {
	return SEARCH_BLOCKS_TIME_INTERVAL
}

func (n NetCoreConfiguration) GetBlockchainHeightBroadcastTimeInterval() uint64 {
	return BLOCKCHAIN_HEIGHT_BROADCASTER_TIME_INTERVAL
}

func (n NetCoreConfiguration) GetBlockBroadcastTimeInterval() uint64 {
	return BLOCK_BROADCASTER_TIME_INTERVAL
}

func (n NetCoreConfiguration) GetAddSeedNodeTimeInterval() uint64 {
	return ADD_SEED_NODE_TIME_INTERVAL
}

func (n NetCoreConfiguration) GetNodeBroadcastTimeInterval() uint64 {
	return NODE_BROADCAST_TIME_INTERVAL
}

func (n NetCoreConfiguration) GetHardForkBlockCount() uint64 {
	return HARD_FORK_BLOCK_COUNT
}

func (n NetCoreConfiguration) GetSearchUnconfirmedTransactionsTimeInterval() uint64 {
	return SEARCH_UNCONFIRMED_TRANSACTIONS_TIME_INTERVAL
}

func (n NetCoreConfiguration) GetNodeCleanTimeInterval() uint64 {
	return NODE_CLEAN_TIME_INTERVAL
}

func (n NetCoreConfiguration) getConfigurationValue(configurationKey []byte) []byte {
	bytesConfigurationValue := KvDbUtil.Get(n.getNetCoreConfigurationDatabasePath(), configurationKey)
	return bytesConfigurationValue
}

func (n NetCoreConfiguration) addOrUpdateConfiguration(configurationKey []byte, configurationValue []byte) {
	KvDbUtil.Put(n.getNetCoreConfigurationDatabasePath(), configurationKey, configurationValue)
}

func (n NetCoreConfiguration) getNetCoreConfigurationDatabasePath() string {
	return FileUtil.NewPath(n.NetCorePath, NETCORE_CONFIGURATION_DATABASE_NAME)
}
