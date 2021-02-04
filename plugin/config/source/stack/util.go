package stack

import (
	"time"

	"github.com/stack-labs/stack/pkg/config/source"
	proto "github.com/stack-labs/stack/plugin/config/source/stack/proto"
)

func toChangeSet(c *proto.ChangeSet) *source.ChangeSet {
	return &source.ChangeSet{
		Data:      c.Data,
		Checksum:  c.Checksum,
		Format:    c.Format,
		Timestamp: time.Unix(c.Timestamp, 0),
		Source:    c.Source,
	}
}
