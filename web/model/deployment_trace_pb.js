// source: pkg/model/deployment_trace.proto
/**
 * @fileoverview
 * @enhanceable
 * @suppress {missingRequire} reports error on implicit type usages.
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!
/* eslint-disable */
// @ts-nocheck

var jspb = require('google-protobuf');
var goog = jspb;
var global =
    (typeof globalThis !== 'undefined' && globalThis) ||
    (typeof window !== 'undefined' && window) ||
    (typeof global !== 'undefined' && global) ||
    (typeof self !== 'undefined' && self) ||
    (function () { return this; }).call(null) ||
    Function('return this')();



goog.exportSymbol('proto.model.DeploymentTrace', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.model.DeploymentTrace = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.model.DeploymentTrace, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.DeploymentTrace.displayName = 'proto.model.DeploymentTrace';
}



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.model.DeploymentTrace.prototype.toObject = function(opt_includeInstance) {
  return proto.model.DeploymentTrace.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.DeploymentTrace} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.DeploymentTrace.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    title: jspb.Message.getFieldWithDefault(msg, 2, ""),
    commitMessage: jspb.Message.getFieldWithDefault(msg, 3, ""),
    commitHash: jspb.Message.getFieldWithDefault(msg, 4, ""),
    commitUrl: jspb.Message.getFieldWithDefault(msg, 5, ""),
    commitTimestamp: jspb.Message.getFieldWithDefault(msg, 6, 0),
    author: jspb.Message.getFieldWithDefault(msg, 7, ""),
    createdAt: jspb.Message.getFieldWithDefault(msg, 101, 0),
    updatedAt: jspb.Message.getFieldWithDefault(msg, 102, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.model.DeploymentTrace}
 */
proto.model.DeploymentTrace.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.DeploymentTrace;
  return proto.model.DeploymentTrace.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.DeploymentTrace} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.DeploymentTrace}
 */
proto.model.DeploymentTrace.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setTitle(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setCommitMessage(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setCommitHash(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setCommitUrl(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCommitTimestamp(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setAuthor(value);
      break;
    case 101:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCreatedAt(value);
      break;
    case 102:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setUpdatedAt(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.model.DeploymentTrace.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.DeploymentTrace.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.DeploymentTrace} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.DeploymentTrace.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getTitle();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getCommitMessage();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getCommitHash();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getCommitUrl();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getCommitTimestamp();
  if (f !== 0) {
    writer.writeInt64(
      6,
      f
    );
  }
  f = message.getAuthor();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(
      101,
      f
    );
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      102,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.model.DeploymentTrace.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.DeploymentTrace} returns this
 */
proto.model.DeploymentTrace.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string title = 2;
 * @return {string}
 */
proto.model.DeploymentTrace.prototype.getTitle = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.DeploymentTrace} returns this
 */
proto.model.DeploymentTrace.prototype.setTitle = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string commit_message = 3;
 * @return {string}
 */
proto.model.DeploymentTrace.prototype.getCommitMessage = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.DeploymentTrace} returns this
 */
proto.model.DeploymentTrace.prototype.setCommitMessage = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string commit_hash = 4;
 * @return {string}
 */
proto.model.DeploymentTrace.prototype.getCommitHash = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.DeploymentTrace} returns this
 */
proto.model.DeploymentTrace.prototype.setCommitHash = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string commit_url = 5;
 * @return {string}
 */
proto.model.DeploymentTrace.prototype.getCommitUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.DeploymentTrace} returns this
 */
proto.model.DeploymentTrace.prototype.setCommitUrl = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional int64 commit_timestamp = 6;
 * @return {number}
 */
proto.model.DeploymentTrace.prototype.getCommitTimestamp = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.DeploymentTrace} returns this
 */
proto.model.DeploymentTrace.prototype.setCommitTimestamp = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * optional string author = 7;
 * @return {string}
 */
proto.model.DeploymentTrace.prototype.getAuthor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.DeploymentTrace} returns this
 */
proto.model.DeploymentTrace.prototype.setAuthor = function(value) {
  return jspb.Message.setProto3StringField(this, 7, value);
};


/**
 * optional int64 created_at = 101;
 * @return {number}
 */
proto.model.DeploymentTrace.prototype.getCreatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 101, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.DeploymentTrace} returns this
 */
proto.model.DeploymentTrace.prototype.setCreatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 101, value);
};


/**
 * optional int64 updated_at = 102;
 * @return {number}
 */
proto.model.DeploymentTrace.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 102, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.DeploymentTrace} returns this
 */
proto.model.DeploymentTrace.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 102, value);
};


goog.object.extend(exports, proto.model);
