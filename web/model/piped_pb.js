// source: pkg/model/piped.proto
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
var global = (function() {
  if (this) { return this; }
  if (typeof window !== 'undefined') { return window; }
  if (typeof global !== 'undefined') { return global; }
  if (typeof self !== 'undefined') { return self; }
  return Function('return this')();
}.call(null));



var pkg_model_common_pb = require('pipecd/web/model/common_pb.js');
goog.object.extend(proto, pkg_model_common_pb);
goog.exportSymbol('proto.model.Piped', null, global);
goog.exportSymbol('proto.model.Piped.CloudProvider', null, global);
goog.exportSymbol('proto.model.Piped.ConnectionStatus', null, global);
goog.exportSymbol('proto.model.Piped.SecretEncryption', null, global);
goog.exportSymbol('proto.model.PipedKey', null, global);
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
proto.model.Piped = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.model.Piped.repeatedFields_, null);
};
goog.inherits(proto.model.Piped, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.Piped.displayName = 'proto.model.Piped';
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
proto.model.Piped.CloudProvider = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.model.Piped.CloudProvider, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.Piped.CloudProvider.displayName = 'proto.model.Piped.CloudProvider';
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
proto.model.Piped.SecretEncryption = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.model.Piped.SecretEncryption, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.Piped.SecretEncryption.displayName = 'proto.model.Piped.SecretEncryption';
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
proto.model.PipedKey = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.model.PipedKey, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.PipedKey.displayName = 'proto.model.PipedKey';
}

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.model.Piped.repeatedFields_ = [9,10,20];



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
proto.model.Piped.prototype.toObject = function(opt_includeInstance) {
  return proto.model.Piped.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.Piped} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.Piped.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    desc: jspb.Message.getFieldWithDefault(msg, 3, ""),
    keyHash: jspb.Message.getFieldWithDefault(msg, 4, ""),
    projectId: jspb.Message.getFieldWithDefault(msg, 5, ""),
    version: jspb.Message.getFieldWithDefault(msg, 7, ""),
    startedAt: jspb.Message.getFieldWithDefault(msg, 8, 0),
    cloudProvidersList: jspb.Message.toObjectList(msg.getCloudProvidersList(),
    proto.model.Piped.CloudProvider.toObject, includeInstance),
    repositoriesList: jspb.Message.toObjectList(msg.getRepositoriesList(),
    pkg_model_common_pb.ApplicationGitRepository.toObject, includeInstance),
    status: jspb.Message.getFieldWithDefault(msg, 11, 0),
    config: jspb.Message.getFieldWithDefault(msg, 12, ""),
    secretEncryption: (f = msg.getSecretEncryption()) && proto.model.Piped.SecretEncryption.toObject(includeInstance, f),
    keysList: jspb.Message.toObjectList(msg.getKeysList(),
    proto.model.PipedKey.toObject, includeInstance),
    desiredVersion: jspb.Message.getFieldWithDefault(msg, 30, ""),
    disabled: jspb.Message.getBooleanFieldWithDefault(msg, 13, false),
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
 * @return {!proto.model.Piped}
 */
proto.model.Piped.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.Piped;
  return proto.model.Piped.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.Piped} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.Piped}
 */
proto.model.Piped.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setDesc(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setKeyHash(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setProjectId(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setVersion(value);
      break;
    case 8:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setStartedAt(value);
      break;
    case 9:
      var value = new proto.model.Piped.CloudProvider;
      reader.readMessage(value,proto.model.Piped.CloudProvider.deserializeBinaryFromReader);
      msg.addCloudProviders(value);
      break;
    case 10:
      var value = new pkg_model_common_pb.ApplicationGitRepository;
      reader.readMessage(value,pkg_model_common_pb.ApplicationGitRepository.deserializeBinaryFromReader);
      msg.addRepositories(value);
      break;
    case 11:
      var value = /** @type {!proto.model.Piped.ConnectionStatus} */ (reader.readEnum());
      msg.setStatus(value);
      break;
    case 12:
      var value = /** @type {string} */ (reader.readString());
      msg.setConfig(value);
      break;
    case 21:
      var value = new proto.model.Piped.SecretEncryption;
      reader.readMessage(value,proto.model.Piped.SecretEncryption.deserializeBinaryFromReader);
      msg.setSecretEncryption(value);
      break;
    case 20:
      var value = new proto.model.PipedKey;
      reader.readMessage(value,proto.model.PipedKey.deserializeBinaryFromReader);
      msg.addKeys(value);
      break;
    case 30:
      var value = /** @type {string} */ (reader.readString());
      msg.setDesiredVersion(value);
      break;
    case 13:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDisabled(value);
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
proto.model.Piped.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.Piped.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.Piped} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.Piped.serializeBinaryToWriter = function(message, writer) {
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
  f = message.getDesc();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getKeyHash();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getProjectId();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getVersion();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
  f = message.getStartedAt();
  if (f !== 0) {
    writer.writeInt64(
      8,
      f
    );
  }
  f = message.getCloudProvidersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      9,
      f,
      proto.model.Piped.CloudProvider.serializeBinaryToWriter
    );
  }
  f = message.getRepositoriesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      10,
      f,
      pkg_model_common_pb.ApplicationGitRepository.serializeBinaryToWriter
    );
  }
  f = message.getStatus();
  if (f !== 0.0) {
    writer.writeEnum(
      11,
      f
    );
  }
  f = message.getConfig();
  if (f.length > 0) {
    writer.writeString(
      12,
      f
    );
  }
  f = message.getSecretEncryption();
  if (f != null) {
    writer.writeMessage(
      21,
      f,
      proto.model.Piped.SecretEncryption.serializeBinaryToWriter
    );
  }
  f = message.getKeysList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      20,
      f,
      proto.model.PipedKey.serializeBinaryToWriter
    );
  }
  f = message.getDesiredVersion();
  if (f.length > 0) {
    writer.writeString(
      30,
      f
    );
  }
  f = message.getDisabled();
  if (f) {
    writer.writeBool(
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
 * @enum {number}
 */
proto.model.Piped.ConnectionStatus = {
  UNKNOWN: 0,
  ONLINE: 1,
  OFFLINE: 2
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
proto.model.Piped.CloudProvider.prototype.toObject = function(opt_includeInstance) {
  return proto.model.Piped.CloudProvider.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.Piped.CloudProvider} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.Piped.CloudProvider.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    type: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.model.Piped.CloudProvider}
 */
proto.model.Piped.CloudProvider.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.Piped.CloudProvider;
  return proto.model.Piped.CloudProvider.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.Piped.CloudProvider} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.Piped.CloudProvider}
 */
proto.model.Piped.CloudProvider.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setType(value);
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
proto.model.Piped.CloudProvider.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.Piped.CloudProvider.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.Piped.CloudProvider} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.Piped.CloudProvider.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getType();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.model.Piped.CloudProvider.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Piped.CloudProvider} returns this
 */
proto.model.Piped.CloudProvider.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string type = 2;
 * @return {string}
 */
proto.model.Piped.CloudProvider.prototype.getType = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Piped.CloudProvider} returns this
 */
proto.model.Piped.CloudProvider.prototype.setType = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
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
proto.model.Piped.SecretEncryption.prototype.toObject = function(opt_includeInstance) {
  return proto.model.Piped.SecretEncryption.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.Piped.SecretEncryption} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.Piped.SecretEncryption.toObject = function(includeInstance, msg) {
  var f, obj = {
    type: jspb.Message.getFieldWithDefault(msg, 1, ""),
    publicKey: jspb.Message.getFieldWithDefault(msg, 2, ""),
    encryptServiceAccount: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.model.Piped.SecretEncryption}
 */
proto.model.Piped.SecretEncryption.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.Piped.SecretEncryption;
  return proto.model.Piped.SecretEncryption.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.Piped.SecretEncryption} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.Piped.SecretEncryption}
 */
proto.model.Piped.SecretEncryption.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setType(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setPublicKey(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setEncryptServiceAccount(value);
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
proto.model.Piped.SecretEncryption.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.Piped.SecretEncryption.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.Piped.SecretEncryption} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.Piped.SecretEncryption.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getType();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPublicKey();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getEncryptServiceAccount();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string type = 1;
 * @return {string}
 */
proto.model.Piped.SecretEncryption.prototype.getType = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Piped.SecretEncryption} returns this
 */
proto.model.Piped.SecretEncryption.prototype.setType = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string public_key = 2;
 * @return {string}
 */
proto.model.Piped.SecretEncryption.prototype.getPublicKey = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Piped.SecretEncryption} returns this
 */
proto.model.Piped.SecretEncryption.prototype.setPublicKey = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string encrypt_service_account = 3;
 * @return {string}
 */
proto.model.Piped.SecretEncryption.prototype.getEncryptServiceAccount = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Piped.SecretEncryption} returns this
 */
proto.model.Piped.SecretEncryption.prototype.setEncryptServiceAccount = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.model.Piped.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Piped} returns this
 */
proto.model.Piped.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.model.Piped.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Piped} returns this
 */
proto.model.Piped.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string desc = 3;
 * @return {string}
 */
proto.model.Piped.prototype.getDesc = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Piped} returns this
 */
proto.model.Piped.prototype.setDesc = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string key_hash = 4;
 * @return {string}
 */
proto.model.Piped.prototype.getKeyHash = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Piped} returns this
 */
proto.model.Piped.prototype.setKeyHash = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string project_id = 5;
 * @return {string}
 */
proto.model.Piped.prototype.getProjectId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Piped} returns this
 */
proto.model.Piped.prototype.setProjectId = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string version = 7;
 * @return {string}
 */
proto.model.Piped.prototype.getVersion = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Piped} returns this
 */
proto.model.Piped.prototype.setVersion = function(value) {
  return jspb.Message.setProto3StringField(this, 7, value);
};


/**
 * optional int64 started_at = 8;
 * @return {number}
 */
proto.model.Piped.prototype.getStartedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 8, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.Piped} returns this
 */
proto.model.Piped.prototype.setStartedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 8, value);
};


/**
 * repeated CloudProvider cloud_providers = 9;
 * @return {!Array<!proto.model.Piped.CloudProvider>}
 */
proto.model.Piped.prototype.getCloudProvidersList = function() {
  return /** @type{!Array<!proto.model.Piped.CloudProvider>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.model.Piped.CloudProvider, 9));
};


/**
 * @param {!Array<!proto.model.Piped.CloudProvider>} value
 * @return {!proto.model.Piped} returns this
*/
proto.model.Piped.prototype.setCloudProvidersList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 9, value);
};


/**
 * @param {!proto.model.Piped.CloudProvider=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.Piped.CloudProvider}
 */
proto.model.Piped.prototype.addCloudProviders = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 9, opt_value, proto.model.Piped.CloudProvider, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.model.Piped} returns this
 */
proto.model.Piped.prototype.clearCloudProvidersList = function() {
  return this.setCloudProvidersList([]);
};


/**
 * repeated ApplicationGitRepository repositories = 10;
 * @return {!Array<!proto.model.ApplicationGitRepository>}
 */
proto.model.Piped.prototype.getRepositoriesList = function() {
  return /** @type{!Array<!proto.model.ApplicationGitRepository>} */ (
    jspb.Message.getRepeatedWrapperField(this, pkg_model_common_pb.ApplicationGitRepository, 10));
};


/**
 * @param {!Array<!proto.model.ApplicationGitRepository>} value
 * @return {!proto.model.Piped} returns this
*/
proto.model.Piped.prototype.setRepositoriesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 10, value);
};


/**
 * @param {!proto.model.ApplicationGitRepository=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.ApplicationGitRepository}
 */
proto.model.Piped.prototype.addRepositories = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 10, opt_value, proto.model.ApplicationGitRepository, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.model.Piped} returns this
 */
proto.model.Piped.prototype.clearRepositoriesList = function() {
  return this.setRepositoriesList([]);
};


/**
 * optional ConnectionStatus status = 11;
 * @return {!proto.model.Piped.ConnectionStatus}
 */
proto.model.Piped.prototype.getStatus = function() {
  return /** @type {!proto.model.Piped.ConnectionStatus} */ (jspb.Message.getFieldWithDefault(this, 11, 0));
};


/**
 * @param {!proto.model.Piped.ConnectionStatus} value
 * @return {!proto.model.Piped} returns this
 */
proto.model.Piped.prototype.setStatus = function(value) {
  return jspb.Message.setProto3EnumField(this, 11, value);
};


/**
 * optional string config = 12;
 * @return {string}
 */
proto.model.Piped.prototype.getConfig = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 12, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Piped} returns this
 */
proto.model.Piped.prototype.setConfig = function(value) {
  return jspb.Message.setProto3StringField(this, 12, value);
};


/**
 * optional SecretEncryption secret_encryption = 21;
 * @return {?proto.model.Piped.SecretEncryption}
 */
proto.model.Piped.prototype.getSecretEncryption = function() {
  return /** @type{?proto.model.Piped.SecretEncryption} */ (
    jspb.Message.getWrapperField(this, proto.model.Piped.SecretEncryption, 21));
};


/**
 * @param {?proto.model.Piped.SecretEncryption|undefined} value
 * @return {!proto.model.Piped} returns this
*/
proto.model.Piped.prototype.setSecretEncryption = function(value) {
  return jspb.Message.setWrapperField(this, 21, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.model.Piped} returns this
 */
proto.model.Piped.prototype.clearSecretEncryption = function() {
  return this.setSecretEncryption(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.model.Piped.prototype.hasSecretEncryption = function() {
  return jspb.Message.getField(this, 21) != null;
};


/**
 * repeated PipedKey keys = 20;
 * @return {!Array<!proto.model.PipedKey>}
 */
proto.model.Piped.prototype.getKeysList = function() {
  return /** @type{!Array<!proto.model.PipedKey>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.model.PipedKey, 20));
};


/**
 * @param {!Array<!proto.model.PipedKey>} value
 * @return {!proto.model.Piped} returns this
*/
proto.model.Piped.prototype.setKeysList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 20, value);
};


/**
 * @param {!proto.model.PipedKey=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.PipedKey}
 */
proto.model.Piped.prototype.addKeys = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 20, opt_value, proto.model.PipedKey, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.model.Piped} returns this
 */
proto.model.Piped.prototype.clearKeysList = function() {
  return this.setKeysList([]);
};


/**
 * optional string desired_version = 30;
 * @return {string}
 */
proto.model.Piped.prototype.getDesiredVersion = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 30, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Piped} returns this
 */
proto.model.Piped.prototype.setDesiredVersion = function(value) {
  return jspb.Message.setProto3StringField(this, 30, value);
};


/**
 * optional bool disabled = 13;
 * @return {boolean}
 */
proto.model.Piped.prototype.getDisabled = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 13, false));
};


/**
 * @param {boolean} value
 * @return {!proto.model.Piped} returns this
 */
proto.model.Piped.prototype.setDisabled = function(value) {
  return jspb.Message.setProto3BooleanField(this, 13, value);
};


/**
 * optional int64 created_at = 14;
 * @return {number}
 */
proto.model.Piped.prototype.getCreatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 14, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.Piped} returns this
 */
proto.model.Piped.prototype.setCreatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 14, value);
};


/**
 * optional int64 updated_at = 15;
 * @return {number}
 */
proto.model.Piped.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 15, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.Piped} returns this
 */
proto.model.Piped.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 15, value);
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
proto.model.PipedKey.prototype.toObject = function(opt_includeInstance) {
  return proto.model.PipedKey.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.PipedKey} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.PipedKey.toObject = function(includeInstance, msg) {
  var f, obj = {
    hash: jspb.Message.getFieldWithDefault(msg, 1, ""),
    creator: jspb.Message.getFieldWithDefault(msg, 2, ""),
    createdAt: jspb.Message.getFieldWithDefault(msg, 10, 0)
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
 * @return {!proto.model.PipedKey}
 */
proto.model.PipedKey.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.PipedKey;
  return proto.model.PipedKey.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.PipedKey} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.PipedKey}
 */
proto.model.PipedKey.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setHash(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setCreator(value);
      break;
    case 10:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCreatedAt(value);
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
proto.model.PipedKey.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.PipedKey.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.PipedKey} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.PipedKey.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getHash();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCreator();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(
      10,
      f
    );
  }
};


/**
 * optional string hash = 1;
 * @return {string}
 */
proto.model.PipedKey.prototype.getHash = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.PipedKey} returns this
 */
proto.model.PipedKey.prototype.setHash = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string creator = 2;
 * @return {string}
 */
proto.model.PipedKey.prototype.getCreator = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.PipedKey} returns this
 */
proto.model.PipedKey.prototype.setCreator = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional int64 created_at = 10;
 * @return {number}
 */
proto.model.PipedKey.prototype.getCreatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 10, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.PipedKey} returns this
 */
proto.model.PipedKey.prototype.setCreatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 10, value);
};


goog.object.extend(exports, proto.model);
