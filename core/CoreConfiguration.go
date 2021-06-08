package core

type CoreConfiguration struct {
	CorePath string
}

func (c *CoreConfiguration) getCorePath() string {
	return c.CorePath
}
