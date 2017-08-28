package accessLogs

import (
	"tls-server/api/middlewares"

	flatbuffers "github.com/google/flatbuffers/go"
)

// MakeListAccessLogs create AccessLog bytes list
func MakeAccessLogs(li *[]middlewares.AccessLog) *[]byte {
	b := flatbuffers.NewBuilder(0)
	ln := len(*li)
	arrLi := make([]flatbuffers.UOffsetT, ln)

	for i, l := range *li {
		idP := b.CreateByteString([]byte(l.ID.Hex()))
		uidP := b.CreateByteString([]byte(l.UID))
		remoteAddrP := b.CreateByteString([]byte(l.RemoteAddr))
		reqContentTypeP := b.CreateByteString([]byte(l.ReqContentType))
		respContentTypeP := b.CreateByteString([]byte(l.RespContentType))
		pathP := b.CreateByteString([]byte(l.Path))
		queryP := b.CreateByteString([]byte(l.Query))
		methodP := b.CreateByteString([]byte(l.Method))
		cachedP := b.CreateByteString([]byte(l.Cached))

		AccessLogStart(b)
		AccessLogAddID(b, idP)
		AccessLogAddRemoteAddr(b, remoteAddrP)
		AccessLogAddUID(b, uidP)
		AccessLogAddReqContentType(b, reqContentTypeP)
		AccessLogAddRespContentType(b, respContentTypeP)
		AccessLogAddReqLength(b, int32(l.ReqLength))
		AccessLogAddRespLength(b, int32(l.RespLength))
		AccessLogAddStatus(b, int32(l.Status))
		AccessLogAddPath(b, pathP)
		AccessLogAddQuery(b, queryP)
		AccessLogAddMethod(b, methodP)
		AccessLogAddCached(b, cachedP)
		AccessLogAddDuration(b, uint64(l.Duration))
		AccessLogAddTimed(b, uint64(l.Timed.Unix()))

		arrLi[i] = AccessLogEnd(b)
	}

	AccessLogsStartListVector(b, ln)
	for i := ln - 1; i >= 0; i-- {
		b.PrependUOffsetT(arrLi[i])
	}
	list := b.EndVector(ln)

	AccessLogsStart(b)
	AccessLogsAddList(b, list)
	b.Finish(AccessLogsEnd(b))

	bts := b.Bytes[b.Head():]
	b = nil
	return &bts
}
