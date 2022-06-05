package configuration

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/util/ByteUtil"
	"helloworldcoin-go/util/FileUtil"
	"helloworldcoin-go/util/KvDbUtil"
)

type NetCoreConfiguration struct {
	NetCorePath string
}

func NewNetCoreConfiguration(netCorePath string) *NetCoreConfiguration {
	FileUtil.MakeDirectory(netCorePath)
	return &NetCoreConfiguration{NetCorePath: netCorePath}
}
func (n *NetCoreConfiguration) getNetCorePath() string {
	return n.NetCorePath
}

const NETCORE_CONFIGURATION_DATABASE_NAME = "NetCoreConfigurationDatabase"

const AUTO_SEARCH_BLOCK_OPTION_KEY = "IS_AUTO_SEARCH_BLOCK"

const AUTO_SEARCH_BLOCK_OPTION_DEFAULT_VALUE = true

const AUTO_SEARCH_NODE_OPTION_KEY = "IS_AUTO_SEARCH_NODE"

const AUTO_SEARCH_NODE_OPTION_DEFAULT_VALUE = true

const SEARCH_NODE_TIME_INTERVAL = 1000 * 60 * 2

const SEARCH_BLOCKCHAIN_HEIGHT_TIME_INTERVAL = 1000 * 60 * 2

const SEARCH_BLOCKS_TIME_INTERVAL = 1000 * 60 * 2

const BLOCKCHAIN_HEIGHT_BROADCASTER_TIME_INTERVAL = 1000 * 20

const BLOCK_BROADCASTER_TIME_INTERVAL = 1000 * 20

const ADD_SEED_NODE_TIME_INTERVAL = 1000 * 60 * 2

const NODE_BROADCAST_TIME_INTERVAL = 1000 * 60 * 2

const NODE_CLEAN_TIME_INTERVAL = 1000 * 60 * 10

const HARD_FORK_BLOCK_COUNT = 100000000

const SEARCH_UNCONFIRMED_TRANSACTIONS_TIME_INTERVAL = 1000 * 60 * 2

func (n *NetCoreConfiguration) GetSeedNodeInitializeTimeInterval() uint64 {
	return ADD_SEED_NODE_TIME_INTERVAL
}
func (n *NetCoreConfiguration) GetNodeSearchTimeInterval() uint64 {
	return SEARCH_NODE_TIME_INTERVAL
}
func (n *NetCoreConfiguration) GetNodeBroadcastTimeInterval() uint64 {
	return NODE_BROADCAST_TIME_INTERVAL
}
func (n *NetCoreConfiguration) GetNodeCleanTimeInterval() uint64 {
	return NODE_CLEAN_TIME_INTERVAL
}

func (n *NetCoreConfiguration) IsAutoSearchBlock() bool {
	bytesConfigurationValue := n.getConfigurationValue(ByteUtil.StringToUtf8Bytes(AUTO_SEARCH_BLOCK_OPTION_KEY))
	if bytesConfigurationValue == nil {
		return AUTO_SEARCH_BLOCK_OPTION_DEFAULT_VALUE
	}
	return ByteUtil.Utf8BytesToBoolean(bytesConfigurationValue)
}
func (n *NetCoreConfiguration) ActiveAutoSearchBlock() {
	n.addOrUpdateConfiguration(ByteUtil.StringToUtf8Bytes(AUTO_SEARCH_BLOCK_OPTION_KEY), ByteUtil.BooleanToUtf8Bytes(true))
}
func (n *NetCoreConfiguration) DeactiveAutoSearchBlock() {
	n.addOrUpdateConfiguration(ByteUtil.StringToUtf8Bytes(AUTO_SEARCH_BLOCK_OPTION_KEY), ByteUtil.BooleanToUtf8Bytes(false))
}

func (n *NetCoreConfiguration) IsAutoSearchNode() bool {
	bytesConfigurationValue := n.getConfigurationValue(ByteUtil.StringToUtf8Bytes(AUTO_SEARCH_NODE_OPTION_KEY))
	if bytesConfigurationValue == nil {
		return AUTO_SEARCH_NODE_OPTION_DEFAULT_VALUE
	}
	return ByteUtil.Utf8BytesToBoolean(bytesConfigurationValue)
}
func (n *NetCoreConfiguration) ActiveAutoSearchNode() {
	n.addOrUpdateConfiguration(ByteUtil.StringToUtf8Bytes(AUTO_SEARCH_NODE_OPTION_KEY), ByteUtil.BooleanToUtf8Bytes(true))
}
func (n *NetCoreConfiguration) DeactiveAutoSearchNode() {
	n.addOrUpdateConfiguration(ByteUtil.StringToUtf8Bytes(AUTO_SEARCH_NODE_OPTION_KEY), ByteUtil.BooleanToUtf8Bytes(false))
}

func (n *NetCoreConfiguration) GetBlockSearchTimeInterval() uint64 {
	return SEARCH_BLOCKS_TIME_INTERVAL
}
func (n *NetCoreConfiguration) GetBlockBroadcastTimeInterval() uint64 {
	return BLOCK_BROADCASTER_TIME_INTERVAL
}

func (n *NetCoreConfiguration) GetBlockchainHeightSearchTimeInterval() uint64 {
	return SEARCH_BLOCKCHAIN_HEIGHT_TIME_INTERVAL
}
func (n *NetCoreConfiguration) GetBlockchainHeightBroadcastTimeInterval() uint64 {
	return BLOCKCHAIN_HEIGHT_BROADCASTER_TIME_INTERVAL
}

func (n *NetCoreConfiguration) GetHardForkBlockCount() uint64 {
	return HARD_FORK_BLOCK_COUNT
}

func (n *NetCoreConfiguration) GetSearchUnconfirmedTransactionsTimeInterval() uint64 {
	return SEARCH_UNCONFIRMED_TRANSACTIONS_TIME_INTERVAL
}

func (n *NetCoreConfiguration) getConfigurationValue(configurationKey []byte) []byte {
	bytesConfigurationValue := KvDbUtil.Get(n.getNetCoreConfigurationDatabasePath(), configurationKey)
	return bytesConfigurationValue
}
func (n *NetCoreConfiguration) addOrUpdateConfiguration(configurationKey []byte, configurationValue []byte) {
	KvDbUtil.Put(n.getNetCoreConfigurationDatabasePath(), configurationKey, configurationValue)
}
func (n *NetCoreConfiguration) getNetCoreConfigurationDatabasePath() string {
	return FileUtil.NewPath(n.NetCorePath, NETCORE_CONFIGURATION_DATABASE_NAME)
}
