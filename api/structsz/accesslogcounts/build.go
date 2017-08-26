package accesslogcounts

import (
	"tls-server/api/middlewares"

	flatbuffers "github.com/google/flatbuffers/go"
)

// MakeAccessLogCounts create AccessLogCount bytes list
func MakeAccessLogCounts(li *[]middlewares.AccessLogCount) []byte {
	b := flatbuffers.NewBuilder(0)
	ln := len(*li)
	arrLi := make([]flatbuffers.UOffsetT, ln)

	for i, l := range *li {
		idP := b.CreateByteString([]byte(l.ID.Hex()))
		pathP := b.CreateByteString([]byte(l.Path))

		AccessLogCountStart(b)
		AccessLogCountAddID(b, idP)
		AccessLogCountAddPath(b, pathP)
		AccessLogCountAddCount(b, uint64(l.Count))
		AccessLogCountAddTimed(b, uint64(l.Timed.Unix()))

		arrLi[i] = AccessLogCountEnd(b)
	}

	AccessLogCountsStartListVector(b, ln)
	for i := ln - 1; i >= 0; i-- {
		b.PrependUOffsetT(arrLi[i])
	}
	list := b.EndVector(ln)

	AccessLogCountsStart(b)
	AccessLogCountsAddList(b, list)
	b.Finish(AccessLogCountsEnd(b))

	bts := b.Bytes[b.Head():]
	b = nil
	return bts
}
