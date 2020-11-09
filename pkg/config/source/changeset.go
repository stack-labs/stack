package source

import (
	//nolint:gosec
	"crypto/md5"
	"fmt"
)

// Sum returns the md5 checksum of the ChangeSet data
func (c *ChangeSet) Sum() string {
	//nolint:gosec
	h := md5.New()
	_, _ = h.Write(c.Data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
