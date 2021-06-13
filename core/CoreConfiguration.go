package core

type CoreConfiguration struct {
	CorePath string
}

func (c *CoreConfiguration) getCorePath() string {
	return c.CorePath
}
func (c *CoreConfiguration) GetMinerMineTimeInterval() uint64 {
	return uint64(10000)
}
