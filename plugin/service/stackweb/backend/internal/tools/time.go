package tools

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/micro/go-log"
	"time"
)

func TimeFixRange(start, end string, startFixed, endFixed time.Duration) (s, e time.Time) {
	s, err := time.Parse(time.RFC3339, start)
	if err != nil {
		s = time.Now()
		log.Log(err)
		s = s.Add(startFixed)
	}

	e, err = time.Parse(time.RFC3339, end)
	if err != nil {
		e = time.Now()
		log.Log(err)
		e = e.Add(endFixed)
	}

	return
}

func PTimestamp(ts *timestamp.Timestamp) (t time.Time) {
	t, _ = ptypes.Timestamp(ts)
	return t
}
