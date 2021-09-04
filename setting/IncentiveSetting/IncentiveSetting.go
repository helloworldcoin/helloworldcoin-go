package IncentiveSetting

/*
 @author king 409060350@qq.com
*/

const BLOCK_TIME uint64 = uint64(1000) * uint64(60) * uint64(10)
const INTERVAL_BLOCK_COUNT uint64 = uint64(6 * 24 * 7 * 2)
const INTERVAL_TIME uint64 = BLOCK_TIME * INTERVAL_BLOCK_COUNT
const BLOCK_INIT_INCENTIVE uint64 = uint64(50) * uint64(100000000)
const INCENTIVE_HALVING_INTERVAL uint64 = uint64(210000)
