// source: pkg/model/deployment_source.proto
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

goog.exportSymbol('proto.model.DeploymentSource', null, global);
goog.exportSymbol('proto.model.PluginApplicationSpec', null, global);
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
proto.model.DeploymentSource = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.model.DeploymentSource, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.DeploymentSource.displayName = 'proto.model.DeploymentSource';
}
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
proto.model.PluginApplicationSpec = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.model.PluginApplicationSpec, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.PluginApplicationSpec.displayName = 'proto.model.PluginApplicationSpec';
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
proto.model.DeploymentSource.prototype.toObject = function(opt_includeInstance) {
  return proto.model.DeploymentSource.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.DeploymentSource} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.DeploymentSource.toObject = function(includeInstance, msg) {
  var f, obj = {
    applicationDirectory: jspb.Message.getFieldWithDefault(msg, 1, ""),
    revision: jspb.Message.getFieldWithDefault(msg, 2, ""),
    applicationConfig: (f = msg.getApplicationConfig()) && proto.model.PluginApplicationSpec.toObject(includeInstance, f)
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
 * @return {!proto.model.DeploymentSource}
 */
proto.model.DeploymentSource.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.DeploymentSource;
  return proto.model.DeploymentSource.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.DeploymentSource} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.DeploymentSource}
 */
proto.model.DeploymentSource.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setApplicationDirectory(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setRevision(value);
      break;
    case 3:
      var value = new proto.model.PluginApplicationSpec;
      reader.readMessage(value,proto.model.PluginApplicationSpec.deserializeBinaryFromReader);
      msg.setApplicationConfig(value);
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
proto.model.DeploymentSource.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.DeploymentSource.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.DeploymentSource} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.DeploymentSource.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplicationDirectory();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRevision();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getApplicationConfig();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.model.PluginApplicationSpec.serializeBinaryToWriter
    );
  }
};


/**
 * optional string application_directory = 1;
 * @return {string}
 */
proto.model.DeploymentSource.prototype.getApplicationDirectory = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.DeploymentSource} returns this
 */
proto.model.DeploymentSource.prototype.setApplicationDirectory = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string revision = 2;
 * @return {string}
 */
proto.model.DeploymentSource.prototype.getRevision = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.DeploymentSource} returns this
 */
proto.model.DeploymentSource.prototype.setRevision = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional PluginApplicationSpec application_config = 3;
 * @return {?proto.model.PluginApplicationSpec}
 */
proto.model.DeploymentSource.prototype.getApplicationConfig = function() {
  return /** @type{?proto.model.PluginApplicationSpec} */ (
    jspb.Message.getWrapperField(this, proto.model.PluginApplicationSpec, 3));
};


/**
 * @param {?proto.model.PluginApplicationSpec|undefined} value
 * @return {!proto.model.DeploymentSource} returns this
*/
proto.model.DeploymentSource.prototype.setApplicationConfig = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.model.DeploymentSource} returns this
 */
proto.model.DeploymentSource.prototype.clearApplicationConfig = function() {
  return this.setApplicationConfig(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.model.DeploymentSource.prototype.hasApplicationConfig = function() {
  return jspb.Message.getField(this, 3) != null;
};





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
proto.model.PluginApplicationSpec.prototype.toObject = function(opt_includeInstance) {
  return proto.model.PluginApplicationSpec.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.PluginApplicationSpec} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.PluginApplicationSpec.toObject = function(includeInstance, msg) {
  var f, obj = {
    kind: jspb.Message.getFieldWithDefault(msg, 1, ""),
    apiVersion: jspb.Message.getFieldWithDefault(msg, 2, ""),
    spec: msg.getSpec_asB64()
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
 * @return {!proto.model.PluginApplicationSpec}
 */
proto.model.PluginApplicationSpec.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.PluginApplicationSpec;
  return proto.model.PluginApplicationSpec.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.PluginApplicationSpec} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.PluginApplicationSpec}
 */
proto.model.PluginApplicationSpec.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setKind(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setApiVersion(value);
      break;
    case 3:
      var value = /** @type {!Uint8Array} */ (reader.readBytes());
      msg.setSpec(value);
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
proto.model.PluginApplicationSpec.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.PluginApplicationSpec.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.PluginApplicationSpec} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.PluginApplicationSpec.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getKind();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getApiVersion();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getSpec_asU8();
  if (f.length > 0) {
    writer.writeBytes(
      3,
      f
    );
  }
};


/**
 * optional string kind = 1;
 * @return {string}
 */
proto.model.PluginApplicationSpec.prototype.getKind = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.PluginApplicationSpec} returns this
 */
proto.model.PluginApplicationSpec.prototype.setKind = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string api_version = 2;
 * @return {string}
 */
proto.model.PluginApplicationSpec.prototype.getApiVersion = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.PluginApplicationSpec} returns this
 */
proto.model.PluginApplicationSpec.prototype.setApiVersion = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional bytes spec = 3;
 * @return {string}
 */
proto.model.PluginApplicationSpec.prototype.getSpec = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * optional bytes spec = 3;
 * This is a type-conversion wrapper around `getSpec()`
 * @return {string}
 */
proto.model.PluginApplicationSpec.prototype.getSpec_asB64 = function() {
  return /** @type {string} */ (jspb.Message.bytesAsB64(
      this.getSpec()));
};


/**
 * optional bytes spec = 3;
 * Note that Uint8Array is not supported on all browsers.
 * @see http://caniuse.com/Uint8Array
 * This is a type-conversion wrapper around `getSpec()`
 * @return {!Uint8Array}
 */
proto.model.PluginApplicationSpec.prototype.getSpec_asU8 = function() {
  return /** @type {!Uint8Array} */ (jspb.Message.bytesAsU8(
      this.getSpec()));
};


/**
 * @param {!(string|Uint8Array)} value
 * @return {!proto.model.PluginApplicationSpec} returns this
 */
proto.model.PluginApplicationSpec.prototype.setSpec = function(value) {
  return jspb.Message.setProto3BytesField(this, 3, value);
};


goog.object.extend(exports, proto.model);
