// source: pkg/model/role.proto
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

var google_protobuf_descriptor_pb = require('google-protobuf/google/protobuf/descriptor_pb.js');
goog.object.extend(proto, google_protobuf_descriptor_pb);
goog.exportSymbol('proto.model.Role', null, global);
goog.exportSymbol('proto.model.Role.ProjectRole', null, global);
goog.exportSymbol('proto.model.role', null, global);
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
proto.model.Role = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.model.Role.repeatedFields_, null);
};
goog.inherits(proto.model.Role, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.Role.displayName = 'proto.model.Role';
}

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.model.Role.repeatedFields_ = [3];



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
proto.model.Role.prototype.toObject = function(opt_includeInstance) {
  return proto.model.Role.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.Role} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.Role.toObject = function(includeInstance, msg) {
  var f, obj = {
    projectId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    projectRole: jspb.Message.getFieldWithDefault(msg, 2, 0),
    projectRbacRoleNamesList: (f = jspb.Message.getRepeatedField(msg, 3)) == null ? undefined : f
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
 * @return {!proto.model.Role}
 */
proto.model.Role.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.Role;
  return proto.model.Role.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.Role} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.Role}
 */
proto.model.Role.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setProjectId(value);
      break;
    case 2:
      var value = /** @type {!proto.model.Role.ProjectRole} */ (reader.readEnum());
      msg.setProjectRole(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.addProjectRbacRoleNames(value);
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
proto.model.Role.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.Role.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.Role} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.Role.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getProjectId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getProjectRole();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getProjectRbacRoleNamesList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      3,
      f
    );
  }
};


/**
 * @enum {number}
 */
proto.model.Role.ProjectRole = {
  VIEWER: 0,
  EDITOR: 1,
  ADMIN: 2
};

/**
 * optional string project_id = 1;
 * @return {string}
 */
proto.model.Role.prototype.getProjectId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.Role} returns this
 */
proto.model.Role.prototype.setProjectId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional ProjectRole project_role = 2;
 * @return {!proto.model.Role.ProjectRole}
 */
proto.model.Role.prototype.getProjectRole = function() {
  return /** @type {!proto.model.Role.ProjectRole} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.model.Role.ProjectRole} value
 * @return {!proto.model.Role} returns this
 */
proto.model.Role.prototype.setProjectRole = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * repeated string project_rbac_role_names = 3;
 * @return {!Array<string>}
 */
proto.model.Role.prototype.getProjectRbacRoleNamesList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 3));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.model.Role} returns this
 */
proto.model.Role.prototype.setProjectRbacRoleNamesList = function(value) {
  return jspb.Message.setField(this, 3, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.model.Role} returns this
 */
proto.model.Role.prototype.addProjectRbacRoleNames = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 3, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.model.Role} returns this
 */
proto.model.Role.prototype.clearProjectRbacRoleNamesList = function() {
  return this.setProjectRbacRoleNamesList([]);
};



/**
 * A tuple of {field number, class constructor} for the extension
 * field named `role`.
 * @type {!jspb.ExtensionFieldInfo<!proto.model.Role>}
 */
proto.model.role = new jspb.ExtensionFieldInfo(
    59090,
    {role: 0},
    proto.model.Role,
     /** @type {?function((boolean|undefined),!jspb.Message=): !Object} */ (
         proto.model.Role.toObject),
    0);

google_protobuf_descriptor_pb.MethodOptions.extensionsBinary[59090] = new jspb.ExtensionFieldBinaryInfo(
    proto.model.role,
    jspb.BinaryReader.prototype.readMessage,
    jspb.BinaryWriter.prototype.writeMessage,
    proto.model.Role.serializeBinaryToWriter,
    proto.model.Role.deserializeBinaryFromReader,
    false);
// This registers the extension field with the extended class, so that
// toObject() will function correctly.
google_protobuf_descriptor_pb.MethodOptions.extensions[59090] = proto.model.role;

goog.object.extend(exports, proto.model);
