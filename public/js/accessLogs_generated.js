// automatically generated by the FlatBuffers compiler, do not modify

/**
 * @const
 * @namespace
 */
var structsz = structsz || {};

/**
 * @constructor
 */
structsz.AccessLog = function() {
  /**
   * @type {flatbuffers.ByteBuffer}
   */
  this.bb = null;

  /**
   * @type {number}
   */
  this.bb_pos = 0;
};

/**
 * @param {number} i
 * @param {flatbuffers.ByteBuffer} bb
 * @returns {structsz.AccessLog}
 */
structsz.AccessLog.prototype.__init = function(i, bb) {
  this.bb_pos = i;
  this.bb = bb;
  return this;
};

/**
 * @param {flatbuffers.ByteBuffer} bb
 * @param {structsz.AccessLog=} obj
 * @returns {structsz.AccessLog}
 */
structsz.AccessLog.getRootAsAccessLog = function(bb, obj) {
  return (obj || new structsz.AccessLog).__init(bb.readInt32(bb.position()) + bb.position(), bb);
};

/**
 * @param {flatbuffers.Encoding=} optionalEncoding
 * @returns {string|Uint8Array|null}
 */
structsz.AccessLog.prototype.ID = function(optionalEncoding) {
  var offset = this.bb.__offset(this.bb_pos, 4);
  return offset ? this.bb.__string(this.bb_pos + offset, optionalEncoding) : null;
};

/**
 * @param {flatbuffers.Encoding=} optionalEncoding
 * @returns {string|Uint8Array|null}
 */
structsz.AccessLog.prototype.RemoteAddr = function(optionalEncoding) {
  var offset = this.bb.__offset(this.bb_pos, 6);
  return offset ? this.bb.__string(this.bb_pos + offset, optionalEncoding) : null;
};

/**
 * @param {flatbuffers.Encoding=} optionalEncoding
 * @returns {string|Uint8Array|null}
 */
structsz.AccessLog.prototype.UID = function(optionalEncoding) {
  var offset = this.bb.__offset(this.bb_pos, 8);
  return offset ? this.bb.__string(this.bb_pos + offset, optionalEncoding) : null;
};

/**
 * @param {flatbuffers.Encoding=} optionalEncoding
 * @returns {string|Uint8Array|null}
 */
structsz.AccessLog.prototype.ReqContentType = function(optionalEncoding) {
  var offset = this.bb.__offset(this.bb_pos, 10);
  return offset ? this.bb.__string(this.bb_pos + offset, optionalEncoding) : null;
};

/**
 * @param {flatbuffers.Encoding=} optionalEncoding
 * @returns {string|Uint8Array|null}
 */
structsz.AccessLog.prototype.RespContentType = function(optionalEncoding) {
  var offset = this.bb.__offset(this.bb_pos, 12);
  return offset ? this.bb.__string(this.bb_pos + offset, optionalEncoding) : null;
};

/**
 * @returns {number}
 */
structsz.AccessLog.prototype.ReqLength = function() {
  var offset = this.bb.__offset(this.bb_pos, 14);
  return offset ? this.bb.readInt32(this.bb_pos + offset) : 0;
};

/**
 * @returns {number}
 */
structsz.AccessLog.prototype.RespLength = function() {
  var offset = this.bb.__offset(this.bb_pos, 16);
  return offset ? this.bb.readInt32(this.bb_pos + offset) : 0;
};

/**
 * @returns {number}
 */
structsz.AccessLog.prototype.Status = function() {
  var offset = this.bb.__offset(this.bb_pos, 18);
  return offset ? this.bb.readInt32(this.bb_pos + offset) : 0;
};

/**
 * @param {flatbuffers.Encoding=} optionalEncoding
 * @returns {string|Uint8Array|null}
 */
structsz.AccessLog.prototype.Path = function(optionalEncoding) {
  var offset = this.bb.__offset(this.bb_pos, 20);
  return offset ? this.bb.__string(this.bb_pos + offset, optionalEncoding) : null;
};

/**
 * @param {flatbuffers.Encoding=} optionalEncoding
 * @returns {string|Uint8Array|null}
 */
structsz.AccessLog.prototype.Query = function(optionalEncoding) {
  var offset = this.bb.__offset(this.bb_pos, 22);
  return offset ? this.bb.__string(this.bb_pos + offset, optionalEncoding) : null;
};

/**
 * @param {flatbuffers.Encoding=} optionalEncoding
 * @returns {string|Uint8Array|null}
 */
structsz.AccessLog.prototype.Method = function(optionalEncoding) {
  var offset = this.bb.__offset(this.bb_pos, 24);
  return offset ? this.bb.__string(this.bb_pos + offset, optionalEncoding) : null;
};

/**
 * @param {flatbuffers.Encoding=} optionalEncoding
 * @returns {string|Uint8Array|null}
 */
structsz.AccessLog.prototype.Cached = function(optionalEncoding) {
  var offset = this.bb.__offset(this.bb_pos, 26);
  return offset ? this.bb.__string(this.bb_pos + offset, optionalEncoding) : null;
};

/**
 * @returns {flatbuffers.Long}
 */
structsz.AccessLog.prototype.Duration = function() {
  var offset = this.bb.__offset(this.bb_pos, 28);
  return offset ? this.bb.readUint64(this.bb_pos + offset) : this.bb.createLong(0, 0);
};

/**
 * @returns {flatbuffers.Long}
 */
structsz.AccessLog.prototype.Timed = function() {
  var offset = this.bb.__offset(this.bb_pos, 30);
  return offset ? this.bb.readUint64(this.bb_pos + offset) : this.bb.createLong(0, 0);
};

/**
 * @param {flatbuffers.Builder} builder
 */
structsz.AccessLog.startAccessLog = function(builder) {
  builder.startObject(14);
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {flatbuffers.Offset} IDOffset
 */
structsz.AccessLog.addID = function(builder, IDOffset) {
  builder.addFieldOffset(0, IDOffset, 0);
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {flatbuffers.Offset} RemoteAddrOffset
 */
structsz.AccessLog.addRemoteAddr = function(builder, RemoteAddrOffset) {
  builder.addFieldOffset(1, RemoteAddrOffset, 0);
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {flatbuffers.Offset} UIDOffset
 */
structsz.AccessLog.addUID = function(builder, UIDOffset) {
  builder.addFieldOffset(2, UIDOffset, 0);
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {flatbuffers.Offset} ReqContentTypeOffset
 */
structsz.AccessLog.addReqContentType = function(builder, ReqContentTypeOffset) {
  builder.addFieldOffset(3, ReqContentTypeOffset, 0);
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {flatbuffers.Offset} RespContentTypeOffset
 */
structsz.AccessLog.addRespContentType = function(builder, RespContentTypeOffset) {
  builder.addFieldOffset(4, RespContentTypeOffset, 0);
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {number} ReqLength
 */
structsz.AccessLog.addReqLength = function(builder, ReqLength) {
  builder.addFieldInt32(5, ReqLength, 0);
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {number} RespLength
 */
structsz.AccessLog.addRespLength = function(builder, RespLength) {
  builder.addFieldInt32(6, RespLength, 0);
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {number} Status
 */
structsz.AccessLog.addStatus = function(builder, Status) {
  builder.addFieldInt32(7, Status, 0);
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {flatbuffers.Offset} PathOffset
 */
structsz.AccessLog.addPath = function(builder, PathOffset) {
  builder.addFieldOffset(8, PathOffset, 0);
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {flatbuffers.Offset} QueryOffset
 */
structsz.AccessLog.addQuery = function(builder, QueryOffset) {
  builder.addFieldOffset(9, QueryOffset, 0);
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {flatbuffers.Offset} MethodOffset
 */
structsz.AccessLog.addMethod = function(builder, MethodOffset) {
  builder.addFieldOffset(10, MethodOffset, 0);
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {flatbuffers.Offset} CachedOffset
 */
structsz.AccessLog.addCached = function(builder, CachedOffset) {
  builder.addFieldOffset(11, CachedOffset, 0);
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {flatbuffers.Long} Duration
 */
structsz.AccessLog.addDuration = function(builder, Duration) {
  builder.addFieldInt64(12, Duration, builder.createLong(0, 0));
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {flatbuffers.Long} Timed
 */
structsz.AccessLog.addTimed = function(builder, Timed) {
  builder.addFieldInt64(13, Timed, builder.createLong(0, 0));
};

/**
 * @param {flatbuffers.Builder} builder
 * @returns {flatbuffers.Offset}
 */
structsz.AccessLog.endAccessLog = function(builder) {
  var offset = builder.endObject();
  return offset;
};

/**
 * @constructor
 */
structsz.AccessLogs = function() {
  /**
   * @type {flatbuffers.ByteBuffer}
   */
  this.bb = null;

  /**
   * @type {number}
   */
  this.bb_pos = 0;
};

/**
 * @param {number} i
 * @param {flatbuffers.ByteBuffer} bb
 * @returns {structsz.AccessLogs}
 */
structsz.AccessLogs.prototype.__init = function(i, bb) {
  this.bb_pos = i;
  this.bb = bb;
  return this;
};

/**
 * @param {flatbuffers.ByteBuffer} bb
 * @param {structsz.AccessLogs=} obj
 * @returns {structsz.AccessLogs}
 */
structsz.AccessLogs.getRootAsAccessLogs = function(bb, obj) {
  return (obj || new structsz.AccessLogs).__init(bb.readInt32(bb.position()) + bb.position(), bb);
};

/**
 * @param {number} index
 * @param {structsz.AccessLog=} obj
 * @returns {structsz.AccessLog}
 */
structsz.AccessLogs.prototype.List = function(index, obj) {
  var offset = this.bb.__offset(this.bb_pos, 4);
  return offset ? (obj || new structsz.AccessLog).__init(this.bb.__indirect(this.bb.__vector(this.bb_pos + offset) + index * 4), this.bb) : null;
};

/**
 * @returns {number}
 */
structsz.AccessLogs.prototype.ListLength = function() {
  var offset = this.bb.__offset(this.bb_pos, 4);
  return offset ? this.bb.__vector_len(this.bb_pos + offset) : 0;
};

/**
 * @param {flatbuffers.Builder} builder
 */
structsz.AccessLogs.startAccessLogs = function(builder) {
  builder.startObject(1);
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {flatbuffers.Offset} ListOffset
 */
structsz.AccessLogs.addList = function(builder, ListOffset) {
  builder.addFieldOffset(0, ListOffset, 0);
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {Array.<flatbuffers.Offset>} data
 * @returns {flatbuffers.Offset}
 */
structsz.AccessLogs.createListVector = function(builder, data) {
  builder.startVector(4, data.length, 4);
  for (var i = data.length - 1; i >= 0; i--) {
    builder.addOffset(data[i]);
  }
  return builder.endVector();
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {number} numElems
 */
structsz.AccessLogs.startListVector = function(builder, numElems) {
  builder.startVector(4, numElems, 4);
};

/**
 * @param {flatbuffers.Builder} builder
 * @returns {flatbuffers.Offset}
 */
structsz.AccessLogs.endAccessLogs = function(builder) {
  var offset = builder.endObject();
  return offset;
};

/**
 * @param {flatbuffers.Builder} builder
 * @param {flatbuffers.Offset} offset
 */
structsz.AccessLogs.finishAccessLogsBuffer = function(builder, offset) {
  builder.finish(offset);
};

