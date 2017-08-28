package middcachez

import (
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
)

// MakeCacheHandlersObj create return CacheHandlersObj bytes
func MakeCacheHandlersObj(status int, ContentType []byte, data *[]byte) *[]byte {
	b := flatbuffers.NewBuilder(0)
	b.Reset()

	contentTypeP := b.CreateByteString(ContentType)
	dataP := b.CreateByteString(*data)

	CacheHandlersObjStart(b)
	CacheHandlersObjAddContentType(b, contentTypeP)
	CacheHandlersObjAddStatus(b, int32(status))
	CacheHandlersObjAddTimed(b, time.Now().Unix())
	CacheHandlersObjAddBody(b, dataP)

	bp := CacheHandlersObjEnd(b)
	b.Finish(bp)
	bts := b.Bytes[b.Head():]
	b = nil
	return &bts
}
