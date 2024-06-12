// source: pkg/model/deployment_chain.proto
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



var pkg_model_deployment_pb = require('pipecd/web/model/deployment_pb.js');
goog.object.extend(proto, pkg_model_deployment_pb);
goog.exportSymbol('proto.model.ChainApplicationRef', null, global);
goog.exportSymbol('proto.model.ChainBlock', null, global);
goog.exportSymbol('proto.model.ChainBlockStatus', null, global);
goog.exportSymbol('proto.model.ChainDeploymentRef', null, global);
goog.exportSymbol('proto.model.ChainNode', null, global);
goog.exportSymbol('proto.model.ChainStatus', null, global);
goog.exportSymbol('proto.model.DeploymentChain', null, global);
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
proto.model.DeploymentChain = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.model.DeploymentChain.repeatedFields_, null);
};
goog.inherits(proto.model.DeploymentChain, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.DeploymentChain.displayName = 'proto.model.DeploymentChain';
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
proto.model.ChainApplicationRef = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.model.ChainApplicationRef, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.ChainApplicationRef.displayName = 'proto.model.ChainApplicationRef';
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
proto.model.ChainDeploymentRef = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.model.ChainDeploymentRef, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.ChainDeploymentRef.displayName = 'proto.model.ChainDeploymentRef';
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
proto.model.ChainNode = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.model.ChainNode, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.ChainNode.displayName = 'proto.model.ChainNode';
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
proto.model.ChainBlock = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.model.ChainBlock.repeatedFields_, null);
};
goog.inherits(proto.model.ChainBlock, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.ChainBlock.displayName = 'proto.model.ChainBlock';
}

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.model.DeploymentChain.repeatedFields_ = [4];



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
proto.model.DeploymentChain.prototype.toObject = function(opt_includeInstance) {
  return proto.model.DeploymentChain.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.DeploymentChain} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.DeploymentChain.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    projectId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    status: jspb.Message.getFieldWithDefault(msg, 3, 0),
    blocksList: jspb.Message.toObjectList(msg.getBlocksList(),
    proto.model.ChainBlock.toObject, includeInstance),
    completedAt: jspb.Message.getFieldWithDefault(msg, 100, 0),
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
 * @return {!proto.model.DeploymentChain}
 */
proto.model.DeploymentChain.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.DeploymentChain;
  return proto.model.DeploymentChain.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.DeploymentChain} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.DeploymentChain}
 */
proto.model.DeploymentChain.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setProjectId(value);
      break;
    case 3:
      var value = /** @type {!proto.model.ChainStatus} */ (reader.readEnum());
      msg.setStatus(value);
      break;
    case 4:
      var value = new proto.model.ChainBlock;
      reader.readMessage(value,proto.model.ChainBlock.deserializeBinaryFromReader);
      msg.addBlocks(value);
      break;
    case 100:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCompletedAt(value);
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
proto.model.DeploymentChain.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.DeploymentChain.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.DeploymentChain} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.DeploymentChain.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getProjectId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getStatus();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
  f = message.getBlocksList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      proto.model.ChainBlock.serializeBinaryToWriter
    );
  }
  f = message.getCompletedAt();
  if (f !== 0) {
    writer.writeInt64(
      100,
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
proto.model.DeploymentChain.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.DeploymentChain} returns this
 */
proto.model.DeploymentChain.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string project_id = 2;
 * @return {string}
 */
proto.model.DeploymentChain.prototype.getProjectId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.DeploymentChain} returns this
 */
proto.model.DeploymentChain.prototype.setProjectId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional ChainStatus status = 3;
 * @return {!proto.model.ChainStatus}
 */
proto.model.DeploymentChain.prototype.getStatus = function() {
  return /** @type {!proto.model.ChainStatus} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {!proto.model.ChainStatus} value
 * @return {!proto.model.DeploymentChain} returns this
 */
proto.model.DeploymentChain.prototype.setStatus = function(value) {
  return jspb.Message.setProto3EnumField(this, 3, value);
};


/**
 * repeated ChainBlock blocks = 4;
 * @return {!Array<!proto.model.ChainBlock>}
 */
proto.model.DeploymentChain.prototype.getBlocksList = function() {
  return /** @type{!Array<!proto.model.ChainBlock>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.model.ChainBlock, 4));
};


/**
 * @param {!Array<!proto.model.ChainBlock>} value
 * @return {!proto.model.DeploymentChain} returns this
*/
proto.model.DeploymentChain.prototype.setBlocksList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.model.ChainBlock=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.ChainBlock}
 */
proto.model.DeploymentChain.prototype.addBlocks = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.model.ChainBlock, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.model.DeploymentChain} returns this
 */
proto.model.DeploymentChain.prototype.clearBlocksList = function() {
  return this.setBlocksList([]);
};


/**
 * optional int64 completed_at = 100;
 * @return {number}
 */
proto.model.DeploymentChain.prototype.getCompletedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 100, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.DeploymentChain} returns this
 */
proto.model.DeploymentChain.prototype.setCompletedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 100, value);
};


/**
 * optional int64 created_at = 101;
 * @return {number}
 */
proto.model.DeploymentChain.prototype.getCreatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 101, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.DeploymentChain} returns this
 */
proto.model.DeploymentChain.prototype.setCreatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 101, value);
};


/**
 * optional int64 updated_at = 102;
 * @return {number}
 */
proto.model.DeploymentChain.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 102, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.DeploymentChain} returns this
 */
proto.model.DeploymentChain.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 102, value);
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
proto.model.ChainApplicationRef.prototype.toObject = function(opt_includeInstance) {
  return proto.model.ChainApplicationRef.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.ChainApplicationRef} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.ChainApplicationRef.toObject = function(includeInstance, msg) {
  var f, obj = {
    applicationId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    applicationName: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.model.ChainApplicationRef}
 */
proto.model.ChainApplicationRef.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.ChainApplicationRef;
  return proto.model.ChainApplicationRef.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.ChainApplicationRef} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.ChainApplicationRef}
 */
proto.model.ChainApplicationRef.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setApplicationId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setApplicationName(value);
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
proto.model.ChainApplicationRef.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.ChainApplicationRef.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.ChainApplicationRef} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.ChainApplicationRef.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplicationId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getApplicationName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string application_id = 1;
 * @return {string}
 */
proto.model.ChainApplicationRef.prototype.getApplicationId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.ChainApplicationRef} returns this
 */
proto.model.ChainApplicationRef.prototype.setApplicationId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string application_name = 2;
 * @return {string}
 */
proto.model.ChainApplicationRef.prototype.getApplicationName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.ChainApplicationRef} returns this
 */
proto.model.ChainApplicationRef.prototype.setApplicationName = function(value) {
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
proto.model.ChainDeploymentRef.prototype.toObject = function(opt_includeInstance) {
  return proto.model.ChainDeploymentRef.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.ChainDeploymentRef} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.ChainDeploymentRef.toObject = function(includeInstance, msg) {
  var f, obj = {
    deploymentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    status: jspb.Message.getFieldWithDefault(msg, 2, 0),
    statusReason: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.model.ChainDeploymentRef}
 */
proto.model.ChainDeploymentRef.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.ChainDeploymentRef;
  return proto.model.ChainDeploymentRef.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.ChainDeploymentRef} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.ChainDeploymentRef}
 */
proto.model.ChainDeploymentRef.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setDeploymentId(value);
      break;
    case 2:
      var value = /** @type {!proto.model.DeploymentStatus} */ (reader.readEnum());
      msg.setStatus(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setStatusReason(value);
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
proto.model.ChainDeploymentRef.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.ChainDeploymentRef.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.ChainDeploymentRef} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.ChainDeploymentRef.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeploymentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStatus();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getStatusReason();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string deployment_id = 1;
 * @return {string}
 */
proto.model.ChainDeploymentRef.prototype.getDeploymentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.ChainDeploymentRef} returns this
 */
proto.model.ChainDeploymentRef.prototype.setDeploymentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional DeploymentStatus status = 2;
 * @return {!proto.model.DeploymentStatus}
 */
proto.model.ChainDeploymentRef.prototype.getStatus = function() {
  return /** @type {!proto.model.DeploymentStatus} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.model.DeploymentStatus} value
 * @return {!proto.model.ChainDeploymentRef} returns this
 */
proto.model.ChainDeploymentRef.prototype.setStatus = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional string status_reason = 3;
 * @return {string}
 */
proto.model.ChainDeploymentRef.prototype.getStatusReason = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.ChainDeploymentRef} returns this
 */
proto.model.ChainDeploymentRef.prototype.setStatusReason = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
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
proto.model.ChainNode.prototype.toObject = function(opt_includeInstance) {
  return proto.model.ChainNode.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.ChainNode} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.ChainNode.toObject = function(includeInstance, msg) {
  var f, obj = {
    applicationRef: (f = msg.getApplicationRef()) && proto.model.ChainApplicationRef.toObject(includeInstance, f),
    deploymentRef: (f = msg.getDeploymentRef()) && proto.model.ChainDeploymentRef.toObject(includeInstance, f)
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
 * @return {!proto.model.ChainNode}
 */
proto.model.ChainNode.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.ChainNode;
  return proto.model.ChainNode.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.ChainNode} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.ChainNode}
 */
proto.model.ChainNode.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.model.ChainApplicationRef;
      reader.readMessage(value,proto.model.ChainApplicationRef.deserializeBinaryFromReader);
      msg.setApplicationRef(value);
      break;
    case 2:
      var value = new proto.model.ChainDeploymentRef;
      reader.readMessage(value,proto.model.ChainDeploymentRef.deserializeBinaryFromReader);
      msg.setDeploymentRef(value);
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
proto.model.ChainNode.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.ChainNode.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.ChainNode} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.ChainNode.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplicationRef();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.model.ChainApplicationRef.serializeBinaryToWriter
    );
  }
  f = message.getDeploymentRef();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.model.ChainDeploymentRef.serializeBinaryToWriter
    );
  }
};


/**
 * optional ChainApplicationRef application_ref = 1;
 * @return {?proto.model.ChainApplicationRef}
 */
proto.model.ChainNode.prototype.getApplicationRef = function() {
  return /** @type{?proto.model.ChainApplicationRef} */ (
    jspb.Message.getWrapperField(this, proto.model.ChainApplicationRef, 1));
};


/**
 * @param {?proto.model.ChainApplicationRef|undefined} value
 * @return {!proto.model.ChainNode} returns this
*/
proto.model.ChainNode.prototype.setApplicationRef = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.model.ChainNode} returns this
 */
proto.model.ChainNode.prototype.clearApplicationRef = function() {
  return this.setApplicationRef(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.model.ChainNode.prototype.hasApplicationRef = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional ChainDeploymentRef deployment_ref = 2;
 * @return {?proto.model.ChainDeploymentRef}
 */
proto.model.ChainNode.prototype.getDeploymentRef = function() {
  return /** @type{?proto.model.ChainDeploymentRef} */ (
    jspb.Message.getWrapperField(this, proto.model.ChainDeploymentRef, 2));
};


/**
 * @param {?proto.model.ChainDeploymentRef|undefined} value
 * @return {!proto.model.ChainNode} returns this
*/
proto.model.ChainNode.prototype.setDeploymentRef = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.model.ChainNode} returns this
 */
proto.model.ChainNode.prototype.clearDeploymentRef = function() {
  return this.setDeploymentRef(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.model.ChainNode.prototype.hasDeploymentRef = function() {
  return jspb.Message.getField(this, 2) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.model.ChainBlock.repeatedFields_ = [1];



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
proto.model.ChainBlock.prototype.toObject = function(opt_includeInstance) {
  return proto.model.ChainBlock.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.ChainBlock} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.ChainBlock.toObject = function(includeInstance, msg) {
  var f, obj = {
    nodesList: jspb.Message.toObjectList(msg.getNodesList(),
    proto.model.ChainNode.toObject, includeInstance),
    status: jspb.Message.getFieldWithDefault(msg, 2, 0),
    startedAt: jspb.Message.getFieldWithDefault(msg, 100, 0),
    completedAt: jspb.Message.getFieldWithDefault(msg, 101, 0)
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
 * @return {!proto.model.ChainBlock}
 */
proto.model.ChainBlock.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.ChainBlock;
  return proto.model.ChainBlock.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.ChainBlock} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.ChainBlock}
 */
proto.model.ChainBlock.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.model.ChainNode;
      reader.readMessage(value,proto.model.ChainNode.deserializeBinaryFromReader);
      msg.addNodes(value);
      break;
    case 2:
      var value = /** @type {!proto.model.ChainBlockStatus} */ (reader.readEnum());
      msg.setStatus(value);
      break;
    case 100:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setStartedAt(value);
      break;
    case 101:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCompletedAt(value);
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
proto.model.ChainBlock.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.ChainBlock.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.ChainBlock} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.ChainBlock.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getNodesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.model.ChainNode.serializeBinaryToWriter
    );
  }
  f = message.getStatus();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getStartedAt();
  if (f !== 0) {
    writer.writeInt64(
      100,
      f
    );
  }
  f = message.getCompletedAt();
  if (f !== 0) {
    writer.writeInt64(
      101,
      f
    );
  }
};


/**
 * repeated ChainNode nodes = 1;
 * @return {!Array<!proto.model.ChainNode>}
 */
proto.model.ChainBlock.prototype.getNodesList = function() {
  return /** @type{!Array<!proto.model.ChainNode>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.model.ChainNode, 1));
};


/**
 * @param {!Array<!proto.model.ChainNode>} value
 * @return {!proto.model.ChainBlock} returns this
*/
proto.model.ChainBlock.prototype.setNodesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.model.ChainNode=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.ChainNode}
 */
proto.model.ChainBlock.prototype.addNodes = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.model.ChainNode, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.model.ChainBlock} returns this
 */
proto.model.ChainBlock.prototype.clearNodesList = function() {
  return this.setNodesList([]);
};


/**
 * optional ChainBlockStatus status = 2;
 * @return {!proto.model.ChainBlockStatus}
 */
proto.model.ChainBlock.prototype.getStatus = function() {
  return /** @type {!proto.model.ChainBlockStatus} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.model.ChainBlockStatus} value
 * @return {!proto.model.ChainBlock} returns this
 */
proto.model.ChainBlock.prototype.setStatus = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional int64 started_at = 100;
 * @return {number}
 */
proto.model.ChainBlock.prototype.getStartedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 100, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.ChainBlock} returns this
 */
proto.model.ChainBlock.prototype.setStartedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 100, value);
};


/**
 * optional int64 completed_at = 101;
 * @return {number}
 */
proto.model.ChainBlock.prototype.getCompletedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 101, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.ChainBlock} returns this
 */
proto.model.ChainBlock.prototype.setCompletedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 101, value);
};


/**
 * @enum {number}
 */
proto.model.ChainStatus = {
  DEPLOYMENT_CHAIN_PENDING: 0,
  DEPLOYMENT_CHAIN_RUNNING: 1,
  DEPLOYMENT_CHAIN_SUCCESS: 2,
  DEPLOYMENT_CHAIN_FAILURE: 3,
  DEPLOYMENT_CHAIN_CANCELLED: 4
};

/**
 * @enum {number}
 */
proto.model.ChainBlockStatus = {
  DEPLOYMENT_BLOCK_PENDING: 0,
  DEPLOYMENT_BLOCK_RUNNING: 1,
  DEPLOYMENT_BLOCK_SUCCESS: 2,
  DEPLOYMENT_BLOCK_FAILURE: 3,
  DEPLOYMENT_BLOCK_CANCELLED: 4
};

goog.object.extend(exports, proto.model);
