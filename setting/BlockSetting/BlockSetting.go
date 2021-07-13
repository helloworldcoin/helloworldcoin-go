package BlockSetting

import "helloworldcoin-go/setting/IncentiveSetting"

const BLOCK_MAX_TRANSACTION_COUNT uint64 = IncentiveSetting.BLOCK_TIME / uint64(1000)
const BLOCK_MAX_CHARACTER_COUNT uint64 = uint64(1024) * uint64(1024)
const NONCE_CHARACTER_COUNT uint64 = uint64(64)
