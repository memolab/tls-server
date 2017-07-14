// automatically generated by the FlatBuffers compiler, do not modify

package accessLogs

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type AccessLog struct {
	_tab flatbuffers.Table
}

func GetRootAsAccessLog(buf []byte, offset flatbuffers.UOffsetT) *AccessLog {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &AccessLog{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *AccessLog) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *AccessLog) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *AccessLog) ID() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *AccessLog) RemoteAddr() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *AccessLog) UID() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *AccessLog) ReqContentType() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *AccessLog) RespContentType() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *AccessLog) ReqLength() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *AccessLog) MutateReqLength(n int32) bool {
	return rcv._tab.MutateInt32Slot(14, n)
}

func (rcv *AccessLog) RespLength() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *AccessLog) MutateRespLength(n int32) bool {
	return rcv._tab.MutateInt32Slot(16, n)
}

func (rcv *AccessLog) Status() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *AccessLog) MutateStatus(n int32) bool {
	return rcv._tab.MutateInt32Slot(18, n)
}

func (rcv *AccessLog) Path() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *AccessLog) Query() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *AccessLog) Method() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *AccessLog) Cached() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *AccessLog) Duration() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(28))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *AccessLog) MutateDuration(n uint64) bool {
	return rcv._tab.MutateUint64Slot(28, n)
}

func (rcv *AccessLog) Timed() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(30))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *AccessLog) MutateTimed(n uint64) bool {
	return rcv._tab.MutateUint64Slot(30, n)
}

func AccessLogStart(builder *flatbuffers.Builder) {
	builder.StartObject(14)
}
func AccessLogAddID(builder *flatbuffers.Builder, ID flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(ID), 0)
}
func AccessLogAddRemoteAddr(builder *flatbuffers.Builder, RemoteAddr flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(RemoteAddr), 0)
}
func AccessLogAddUID(builder *flatbuffers.Builder, UID flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(UID), 0)
}
func AccessLogAddReqContentType(builder *flatbuffers.Builder, ReqContentType flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(ReqContentType), 0)
}
func AccessLogAddRespContentType(builder *flatbuffers.Builder, RespContentType flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(RespContentType), 0)
}
func AccessLogAddReqLength(builder *flatbuffers.Builder, ReqLength int32) {
	builder.PrependInt32Slot(5, ReqLength, 0)
}
func AccessLogAddRespLength(builder *flatbuffers.Builder, RespLength int32) {
	builder.PrependInt32Slot(6, RespLength, 0)
}
func AccessLogAddStatus(builder *flatbuffers.Builder, Status int32) {
	builder.PrependInt32Slot(7, Status, 0)
}
func AccessLogAddPath(builder *flatbuffers.Builder, Path flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(8, flatbuffers.UOffsetT(Path), 0)
}
func AccessLogAddQuery(builder *flatbuffers.Builder, Query flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(9, flatbuffers.UOffsetT(Query), 0)
}
func AccessLogAddMethod(builder *flatbuffers.Builder, Method flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(10, flatbuffers.UOffsetT(Method), 0)
}
func AccessLogAddCached(builder *flatbuffers.Builder, Cached flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(11, flatbuffers.UOffsetT(Cached), 0)
}
func AccessLogAddDuration(builder *flatbuffers.Builder, Duration uint64) {
	builder.PrependUint64Slot(12, Duration, 0)
}
func AccessLogAddTimed(builder *flatbuffers.Builder, Timed uint64) {
	builder.PrependUint64Slot(13, Timed, 0)
}
func AccessLogEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}