// source: pkg/model/event.proto
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
var global = Function('return this')();



goog.exportSymbol('proto.model.Event', null, global);
goog.exportSymbol('proto.model.EventStatus', null, global);
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
proto.model.Event = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.model.Event, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.Event.displayName = 'proto.model.Event';
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
proto.model.Event.prototype.toObject = function(opt_includeInstance) {
  return proto.model.Event.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.Event} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.Event.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    data: jspb.Message.getFieldWithDefault(msg, 3, ""),
    projectId: jspb.Message.getFieldWithDefault(msg, 4, ""),
    labelsMap: (f = msg.getLabelsMap()) ? f.toObject(includeInstance, undefined) : [],
    eventKey: jspb.Message.getFieldWithDefault(msg, 6, ""),
    status: jspb.Message.getFieldWithDefault(msg, 8, 0),
    statusDescription: jspb.Message.getFieldWithDefault(msg, 9, ""),
    handledAt: jspb.Message.getFieldWithDefault(msg, 13, 0),
    createdAt: jspb.Message.getFieldWithDefault(msg, 14, 0),
    updatedAt: jspb.Message.getFieldWithDefault(msg, 15, 0)
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
 * @return {!proto.model.Event}
 */
proto.model.Event.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.Event;
  return proto.model.Event.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.Event} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.Event}
 */
proto.model.Event.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setName(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setData(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setProjectId(value);
      break;
    case 5:
      var value = msg.getLabelsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "", "");
         });
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setEventKey(value);
      break;
    case 8:
      var value = /** @type {!proto.model.EventStatus} */ (reader.readEnum());
      msg.setStatus(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setStatusDescription(value);
      break;
    case 13:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setHandledAt(value);
      break;
    case 14:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCreatedAt(value);
      break;
    case 15:
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
proto.model.Event.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.Event.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.Event} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.Event.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getData();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getProjectId();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getLabelsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(5, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
  f = message.getEventKey();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getStatus();
  if (f !== 0.0) {
    writer.writeEnum(
      8,
      f
    );
  }
  f = message.getStatusDescription();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
  f = message.getHandledAt();
  if (f !== 0) {
    writer.writeInt64(
      13,
      f
    );
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(
      14,
      f
    );
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      15,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.model.Event.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Event} returns this
 */
proto.model.Event.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.model.Event.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Event} returns this
 */
proto.model.Event.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string data = 3;
 * @return {string}
 */
proto.model.Event.prototype.getData = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Event} returns this
 */
proto.model.Event.prototype.setData = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string project_id = 4;
 * @return {string}
 */
proto.model.Event.prototype.getProjectId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Event} returns this
 */
proto.model.Event.prototype.setProjectId = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * map<string, string> labels = 5;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.model.Event.prototype.getLabelsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 5, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.model.Event} returns this
 */
proto.model.Event.prototype.clearLabelsMap = function() {
  this.getLabelsMap().clear();
  return this;};


/**
 * optional string event_key = 6;
 * @return {string}
 */
proto.model.Event.prototype.getEventKey = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Event} returns this
 */
proto.model.Event.prototype.setEventKey = function(value) {
  return jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional EventStatus status = 8;
 * @return {!proto.model.EventStatus}
 */
proto.model.Event.prototype.getStatus = function() {
  return /** @type {!proto.model.EventStatus} */ (jspb.Message.getFieldWithDefault(this, 8, 0));
};


/**
 * @param {!proto.model.EventStatus} value
 * @return {!proto.model.Event} returns this
 */
proto.model.Event.prototype.setStatus = function(value) {
  return jspb.Message.setProto3EnumField(this, 8, value);
};


/**
 * optional string status_description = 9;
 * @return {string}
 */
proto.model.Event.prototype.getStatusDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Event} returns this
 */
proto.model.Event.prototype.setStatusDescription = function(value) {
  return jspb.Message.setProto3StringField(this, 9, value);
};


/**
 * optional int64 handled_at = 13;
 * @return {number}
 */
proto.model.Event.prototype.getHandledAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 13, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.Event} returns this
 */
proto.model.Event.prototype.setHandledAt = function(value) {
  return jspb.Message.setProto3IntField(this, 13, value);
};


/**
 * optional int64 created_at = 14;
 * @return {number}
 */
proto.model.Event.prototype.getCreatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 14, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.Event} returns this
 */
proto.model.Event.prototype.setCreatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 14, value);
};


/**
 * optional int64 updated_at = 15;
 * @return {number}
 */
proto.model.Event.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 15, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.Event} returns this
 */
proto.model.Event.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 15, value);
};


/**
 * @enum {number}
 */
proto.model.EventStatus = {
  EVENT_NOT_HANDLED: 0,
  EVENT_SUCCESS: 1,
  EVENT_FAILURE: 2,
  EVENT_OUTDATED: 3
};

goog.object.extend(exports, proto.model);
