package FileUtil

import (
	"path/filepath"
)

func NewPath(parent string, child string) string {
	return filepath.Join(parent, child)
}
