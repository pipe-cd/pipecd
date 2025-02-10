// source: pkg/app/server/service/webservice/service.proto
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



var pkg_model_common_pb = require('pipecd/web/model/common_pb.js');
goog.object.extend(proto, pkg_model_common_pb);
var pkg_model_insight_pb = require('pipecd/web/model/insight_pb.js');
goog.object.extend(proto, pkg_model_insight_pb);
var pkg_model_application_pb = require('pipecd/web/model/application_pb.js');
goog.object.extend(proto, pkg_model_application_pb);
var pkg_model_application_live_state_pb = require('pipecd/web/model/application_live_state_pb.js');
goog.object.extend(proto, pkg_model_application_live_state_pb);
var pkg_model_command_pb = require('pipecd/web/model/command_pb.js');
goog.object.extend(proto, pkg_model_command_pb);
var pkg_model_deployment_pb = require('pipecd/web/model/deployment_pb.js');
goog.object.extend(proto, pkg_model_deployment_pb);
var pkg_model_deployment_chain_pb = require('pipecd/web/model/deployment_chain_pb.js');
goog.object.extend(proto, pkg_model_deployment_chain_pb);
var pkg_model_logblock_pb = require('pipecd/web/model/logblock_pb.js');
goog.object.extend(proto, pkg_model_logblock_pb);
var pkg_model_piped_pb = require('pipecd/web/model/piped_pb.js');
goog.object.extend(proto, pkg_model_piped_pb);
var pkg_model_rbac_pb = require('pipecd/web/model/rbac_pb.js');
goog.object.extend(proto, pkg_model_rbac_pb);
var pkg_model_project_pb = require('pipecd/web/model/project_pb.js');
goog.object.extend(proto, pkg_model_project_pb);
var pkg_model_apikey_pb = require('pipecd/web/model/apikey_pb.js');
goog.object.extend(proto, pkg_model_apikey_pb);
var pkg_model_event_pb = require('pipecd/web/model/event_pb.js');
goog.object.extend(proto, pkg_model_event_pb);
var google_protobuf_wrappers_pb = require('google-protobuf/google/protobuf/wrappers_pb.js');
goog.object.extend(proto, google_protobuf_wrappers_pb);
var google_protobuf_descriptor_pb = require('google-protobuf/google/protobuf/descriptor_pb.js');
goog.object.extend(proto, google_protobuf_descriptor_pb);
goog.exportSymbol('proto.grpc.service.webservice.AddApplicationRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.AddApplicationResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.AddProjectRBACRoleRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.AddProjectRBACRoleResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.AddProjectUserGroupRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.AddProjectUserGroupResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ApproveStageRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ApproveStageResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.CancelDeploymentRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.CancelDeploymentResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.DeleteApplicationRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.DeleteApplicationResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.DeleteOldPipedKeysRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.DeleteOldPipedKeysResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.DeleteProjectRBACRoleRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.DeleteProjectRBACRoleResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.DeleteProjectUserGroupRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.DeleteProjectUserGroupResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.DisableAPIKeyRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.DisableAPIKeyResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.DisableApplicationRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.DisableApplicationResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.DisablePipedRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.DisablePipedResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.DisableStaticAdminRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.DisableStaticAdminResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.EnableApplicationRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.EnableApplicationResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.EnablePipedRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.EnablePipedResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.EnableStaticAdminRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.EnableStaticAdminResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GenerateAPIKeyRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GenerateAPIKeyResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetApplicationLiveStateRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetApplicationLiveStateResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetApplicationRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetApplicationResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetCommandRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetCommandResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetDeploymentChainRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetDeploymentChainResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetDeploymentRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetDeploymentResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetInsightApplicationCountRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetInsightApplicationCountResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetInsightDataRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetInsightDataResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetMeRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetMeResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetPipedRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetPipedResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetProjectRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetProjectResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetStageLogRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.GetStageLogResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListAPIKeysRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListAPIKeysRequest.Options', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListAPIKeysResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListApplicationsRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListApplicationsRequest.Options', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListApplicationsResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListDeploymentChainsRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListDeploymentChainsRequest.Options', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListDeploymentChainsResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListDeploymentsRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListDeploymentsRequest.Options', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListDeploymentsResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListDeprecatedNotesRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListDeprecatedNotesResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListEventsRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListEventsRequest.Options', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListEventsResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListPipedsRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListPipedsRequest.Options', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListPipedsResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListReleasedVersionsRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListReleasedVersionsResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListUnregisteredApplicationsRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.ListUnregisteredApplicationsResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.RecreatePipedKeyRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.RecreatePipedKeyResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.RegisterPipedRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.RegisterPipedResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.RestartPipedRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.RestartPipedResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.SkipStageRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.SkipStageResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.SyncApplicationRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.SyncApplicationResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.UpdateApplicationRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.UpdateApplicationResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.UpdatePipedRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.UpdatePipedResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.UpdateProjectRBACConfigRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.UpdateProjectRBACConfigResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.UpdateProjectRBACRoleRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.UpdateProjectRBACRoleResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.UpdateProjectSSOConfigRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.UpdateProjectSSOConfigResponse', null, global);
goog.exportSymbol('proto.grpc.service.webservice.UpdateProjectStaticAdminRequest', null, global);
goog.exportSymbol('proto.grpc.service.webservice.UpdateProjectStaticAdminResponse', null, global);
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
proto.grpc.service.webservice.RegisterPipedRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.RegisterPipedRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.RegisterPipedRequest.displayName = 'proto.grpc.service.webservice.RegisterPipedRequest';
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
proto.grpc.service.webservice.RegisterPipedResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.RegisterPipedResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.RegisterPipedResponse.displayName = 'proto.grpc.service.webservice.RegisterPipedResponse';
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
proto.grpc.service.webservice.UpdatePipedRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.UpdatePipedRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.UpdatePipedRequest.displayName = 'proto.grpc.service.webservice.UpdatePipedRequest';
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
proto.grpc.service.webservice.UpdatePipedResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.UpdatePipedResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.UpdatePipedResponse.displayName = 'proto.grpc.service.webservice.UpdatePipedResponse';
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
proto.grpc.service.webservice.RecreatePipedKeyRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.RecreatePipedKeyRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.RecreatePipedKeyRequest.displayName = 'proto.grpc.service.webservice.RecreatePipedKeyRequest';
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
proto.grpc.service.webservice.RecreatePipedKeyResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.RecreatePipedKeyResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.RecreatePipedKeyResponse.displayName = 'proto.grpc.service.webservice.RecreatePipedKeyResponse';
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
proto.grpc.service.webservice.DeleteOldPipedKeysRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.DeleteOldPipedKeysRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.DeleteOldPipedKeysRequest.displayName = 'proto.grpc.service.webservice.DeleteOldPipedKeysRequest';
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
proto.grpc.service.webservice.DeleteOldPipedKeysResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.DeleteOldPipedKeysResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.DeleteOldPipedKeysResponse.displayName = 'proto.grpc.service.webservice.DeleteOldPipedKeysResponse';
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
proto.grpc.service.webservice.EnablePipedRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.EnablePipedRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.EnablePipedRequest.displayName = 'proto.grpc.service.webservice.EnablePipedRequest';
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
proto.grpc.service.webservice.EnablePipedResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.EnablePipedResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.EnablePipedResponse.displayName = 'proto.grpc.service.webservice.EnablePipedResponse';
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
proto.grpc.service.webservice.DisablePipedRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.DisablePipedRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.DisablePipedRequest.displayName = 'proto.grpc.service.webservice.DisablePipedRequest';
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
proto.grpc.service.webservice.DisablePipedResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.DisablePipedResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.DisablePipedResponse.displayName = 'proto.grpc.service.webservice.DisablePipedResponse';
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
proto.grpc.service.webservice.ListPipedsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.ListPipedsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListPipedsRequest.displayName = 'proto.grpc.service.webservice.ListPipedsRequest';
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
proto.grpc.service.webservice.ListPipedsRequest.Options = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.ListPipedsRequest.Options, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListPipedsRequest.Options.displayName = 'proto.grpc.service.webservice.ListPipedsRequest.Options';
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
proto.grpc.service.webservice.ListPipedsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.grpc.service.webservice.ListPipedsResponse.repeatedFields_, null);
};
goog.inherits(proto.grpc.service.webservice.ListPipedsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListPipedsResponse.displayName = 'proto.grpc.service.webservice.ListPipedsResponse';
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
proto.grpc.service.webservice.GetPipedRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetPipedRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetPipedRequest.displayName = 'proto.grpc.service.webservice.GetPipedRequest';
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
proto.grpc.service.webservice.GetPipedResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetPipedResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetPipedResponse.displayName = 'proto.grpc.service.webservice.GetPipedResponse';
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
proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.repeatedFields_, null);
};
goog.inherits(proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.displayName = 'proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest';
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
proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse.displayName = 'proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse';
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
proto.grpc.service.webservice.RestartPipedRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.RestartPipedRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.RestartPipedRequest.displayName = 'proto.grpc.service.webservice.RestartPipedRequest';
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
proto.grpc.service.webservice.RestartPipedResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.RestartPipedResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.RestartPipedResponse.displayName = 'proto.grpc.service.webservice.RestartPipedResponse';
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
proto.grpc.service.webservice.ListReleasedVersionsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.ListReleasedVersionsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListReleasedVersionsRequest.displayName = 'proto.grpc.service.webservice.ListReleasedVersionsRequest';
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
proto.grpc.service.webservice.ListReleasedVersionsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.grpc.service.webservice.ListReleasedVersionsResponse.repeatedFields_, null);
};
goog.inherits(proto.grpc.service.webservice.ListReleasedVersionsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListReleasedVersionsResponse.displayName = 'proto.grpc.service.webservice.ListReleasedVersionsResponse';
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
proto.grpc.service.webservice.ListDeprecatedNotesRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.ListDeprecatedNotesRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListDeprecatedNotesRequest.displayName = 'proto.grpc.service.webservice.ListDeprecatedNotesRequest';
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
proto.grpc.service.webservice.ListDeprecatedNotesResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.ListDeprecatedNotesResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListDeprecatedNotesResponse.displayName = 'proto.grpc.service.webservice.ListDeprecatedNotesResponse';
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
proto.grpc.service.webservice.AddApplicationRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.AddApplicationRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.AddApplicationRequest.displayName = 'proto.grpc.service.webservice.AddApplicationRequest';
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
proto.grpc.service.webservice.AddApplicationResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.AddApplicationResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.AddApplicationResponse.displayName = 'proto.grpc.service.webservice.AddApplicationResponse';
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
proto.grpc.service.webservice.UpdateApplicationRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.UpdateApplicationRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.UpdateApplicationRequest.displayName = 'proto.grpc.service.webservice.UpdateApplicationRequest';
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
proto.grpc.service.webservice.UpdateApplicationResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.UpdateApplicationResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.UpdateApplicationResponse.displayName = 'proto.grpc.service.webservice.UpdateApplicationResponse';
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
proto.grpc.service.webservice.EnableApplicationRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.EnableApplicationRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.EnableApplicationRequest.displayName = 'proto.grpc.service.webservice.EnableApplicationRequest';
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
proto.grpc.service.webservice.EnableApplicationResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.EnableApplicationResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.EnableApplicationResponse.displayName = 'proto.grpc.service.webservice.EnableApplicationResponse';
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
proto.grpc.service.webservice.DisableApplicationRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.DisableApplicationRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.DisableApplicationRequest.displayName = 'proto.grpc.service.webservice.DisableApplicationRequest';
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
proto.grpc.service.webservice.DisableApplicationResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.DisableApplicationResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.DisableApplicationResponse.displayName = 'proto.grpc.service.webservice.DisableApplicationResponse';
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
proto.grpc.service.webservice.DeleteApplicationRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.DeleteApplicationRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.DeleteApplicationRequest.displayName = 'proto.grpc.service.webservice.DeleteApplicationRequest';
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
proto.grpc.service.webservice.DeleteApplicationResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.DeleteApplicationResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.DeleteApplicationResponse.displayName = 'proto.grpc.service.webservice.DeleteApplicationResponse';
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
proto.grpc.service.webservice.ListApplicationsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.ListApplicationsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListApplicationsRequest.displayName = 'proto.grpc.service.webservice.ListApplicationsRequest';
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
proto.grpc.service.webservice.ListApplicationsRequest.Options = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.grpc.service.webservice.ListApplicationsRequest.Options.repeatedFields_, null);
};
goog.inherits(proto.grpc.service.webservice.ListApplicationsRequest.Options, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListApplicationsRequest.Options.displayName = 'proto.grpc.service.webservice.ListApplicationsRequest.Options';
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
proto.grpc.service.webservice.ListApplicationsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.grpc.service.webservice.ListApplicationsResponse.repeatedFields_, null);
};
goog.inherits(proto.grpc.service.webservice.ListApplicationsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListApplicationsResponse.displayName = 'proto.grpc.service.webservice.ListApplicationsResponse';
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
proto.grpc.service.webservice.SyncApplicationRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.SyncApplicationRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.SyncApplicationRequest.displayName = 'proto.grpc.service.webservice.SyncApplicationRequest';
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
proto.grpc.service.webservice.SyncApplicationResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.SyncApplicationResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.SyncApplicationResponse.displayName = 'proto.grpc.service.webservice.SyncApplicationResponse';
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
proto.grpc.service.webservice.GetApplicationRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetApplicationRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetApplicationRequest.displayName = 'proto.grpc.service.webservice.GetApplicationRequest';
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
proto.grpc.service.webservice.GetApplicationResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetApplicationResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetApplicationResponse.displayName = 'proto.grpc.service.webservice.GetApplicationResponse';
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
proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest.displayName = 'proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest';
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
proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse.displayName = 'proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse';
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
proto.grpc.service.webservice.ListUnregisteredApplicationsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.ListUnregisteredApplicationsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListUnregisteredApplicationsRequest.displayName = 'proto.grpc.service.webservice.ListUnregisteredApplicationsRequest';
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
proto.grpc.service.webservice.ListUnregisteredApplicationsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.grpc.service.webservice.ListUnregisteredApplicationsResponse.repeatedFields_, null);
};
goog.inherits(proto.grpc.service.webservice.ListUnregisteredApplicationsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListUnregisteredApplicationsResponse.displayName = 'proto.grpc.service.webservice.ListUnregisteredApplicationsResponse';
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
proto.grpc.service.webservice.ListDeploymentsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.ListDeploymentsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListDeploymentsRequest.displayName = 'proto.grpc.service.webservice.ListDeploymentsRequest';
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
proto.grpc.service.webservice.ListDeploymentsRequest.Options = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.grpc.service.webservice.ListDeploymentsRequest.Options.repeatedFields_, null);
};
goog.inherits(proto.grpc.service.webservice.ListDeploymentsRequest.Options, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListDeploymentsRequest.Options.displayName = 'proto.grpc.service.webservice.ListDeploymentsRequest.Options';
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
proto.grpc.service.webservice.ListDeploymentsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.grpc.service.webservice.ListDeploymentsResponse.repeatedFields_, null);
};
goog.inherits(proto.grpc.service.webservice.ListDeploymentsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListDeploymentsResponse.displayName = 'proto.grpc.service.webservice.ListDeploymentsResponse';
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
proto.grpc.service.webservice.GetDeploymentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetDeploymentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetDeploymentRequest.displayName = 'proto.grpc.service.webservice.GetDeploymentRequest';
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
proto.grpc.service.webservice.GetDeploymentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetDeploymentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetDeploymentResponse.displayName = 'proto.grpc.service.webservice.GetDeploymentResponse';
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
proto.grpc.service.webservice.GetStageLogRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetStageLogRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetStageLogRequest.displayName = 'proto.grpc.service.webservice.GetStageLogRequest';
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
proto.grpc.service.webservice.GetStageLogResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.grpc.service.webservice.GetStageLogResponse.repeatedFields_, null);
};
goog.inherits(proto.grpc.service.webservice.GetStageLogResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetStageLogResponse.displayName = 'proto.grpc.service.webservice.GetStageLogResponse';
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
proto.grpc.service.webservice.CancelDeploymentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.CancelDeploymentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.CancelDeploymentRequest.displayName = 'proto.grpc.service.webservice.CancelDeploymentRequest';
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
proto.grpc.service.webservice.CancelDeploymentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.CancelDeploymentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.CancelDeploymentResponse.displayName = 'proto.grpc.service.webservice.CancelDeploymentResponse';
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
proto.grpc.service.webservice.SkipStageRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.SkipStageRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.SkipStageRequest.displayName = 'proto.grpc.service.webservice.SkipStageRequest';
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
proto.grpc.service.webservice.SkipStageResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.SkipStageResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.SkipStageResponse.displayName = 'proto.grpc.service.webservice.SkipStageResponse';
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
proto.grpc.service.webservice.ApproveStageRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.ApproveStageRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ApproveStageRequest.displayName = 'proto.grpc.service.webservice.ApproveStageRequest';
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
proto.grpc.service.webservice.ApproveStageResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.ApproveStageResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ApproveStageResponse.displayName = 'proto.grpc.service.webservice.ApproveStageResponse';
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
proto.grpc.service.webservice.GetApplicationLiveStateRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetApplicationLiveStateRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetApplicationLiveStateRequest.displayName = 'proto.grpc.service.webservice.GetApplicationLiveStateRequest';
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
proto.grpc.service.webservice.GetApplicationLiveStateResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetApplicationLiveStateResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetApplicationLiveStateResponse.displayName = 'proto.grpc.service.webservice.GetApplicationLiveStateResponse';
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
proto.grpc.service.webservice.GetProjectRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetProjectRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetProjectRequest.displayName = 'proto.grpc.service.webservice.GetProjectRequest';
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
proto.grpc.service.webservice.GetProjectResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetProjectResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetProjectResponse.displayName = 'proto.grpc.service.webservice.GetProjectResponse';
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
proto.grpc.service.webservice.UpdateProjectStaticAdminRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.UpdateProjectStaticAdminRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.UpdateProjectStaticAdminRequest.displayName = 'proto.grpc.service.webservice.UpdateProjectStaticAdminRequest';
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
proto.grpc.service.webservice.UpdateProjectStaticAdminResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.UpdateProjectStaticAdminResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.UpdateProjectStaticAdminResponse.displayName = 'proto.grpc.service.webservice.UpdateProjectStaticAdminResponse';
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
proto.grpc.service.webservice.UpdateProjectSSOConfigRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.UpdateProjectSSOConfigRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.UpdateProjectSSOConfigRequest.displayName = 'proto.grpc.service.webservice.UpdateProjectSSOConfigRequest';
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
proto.grpc.service.webservice.UpdateProjectSSOConfigResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.UpdateProjectSSOConfigResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.UpdateProjectSSOConfigResponse.displayName = 'proto.grpc.service.webservice.UpdateProjectSSOConfigResponse';
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
proto.grpc.service.webservice.UpdateProjectRBACConfigRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.UpdateProjectRBACConfigRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.UpdateProjectRBACConfigRequest.displayName = 'proto.grpc.service.webservice.UpdateProjectRBACConfigRequest';
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
proto.grpc.service.webservice.UpdateProjectRBACConfigResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.UpdateProjectRBACConfigResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.UpdateProjectRBACConfigResponse.displayName = 'proto.grpc.service.webservice.UpdateProjectRBACConfigResponse';
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
proto.grpc.service.webservice.EnableStaticAdminRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.EnableStaticAdminRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.EnableStaticAdminRequest.displayName = 'proto.grpc.service.webservice.EnableStaticAdminRequest';
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
proto.grpc.service.webservice.EnableStaticAdminResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.EnableStaticAdminResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.EnableStaticAdminResponse.displayName = 'proto.grpc.service.webservice.EnableStaticAdminResponse';
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
proto.grpc.service.webservice.DisableStaticAdminRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.DisableStaticAdminRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.DisableStaticAdminRequest.displayName = 'proto.grpc.service.webservice.DisableStaticAdminRequest';
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
proto.grpc.service.webservice.DisableStaticAdminResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.DisableStaticAdminResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.DisableStaticAdminResponse.displayName = 'proto.grpc.service.webservice.DisableStaticAdminResponse';
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
proto.grpc.service.webservice.GetMeRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetMeRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetMeRequest.displayName = 'proto.grpc.service.webservice.GetMeRequest';
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
proto.grpc.service.webservice.GetMeResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetMeResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetMeResponse.displayName = 'proto.grpc.service.webservice.GetMeResponse';
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
proto.grpc.service.webservice.AddProjectRBACRoleRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.grpc.service.webservice.AddProjectRBACRoleRequest.repeatedFields_, null);
};
goog.inherits(proto.grpc.service.webservice.AddProjectRBACRoleRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.AddProjectRBACRoleRequest.displayName = 'proto.grpc.service.webservice.AddProjectRBACRoleRequest';
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
proto.grpc.service.webservice.AddProjectRBACRoleResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.AddProjectRBACRoleResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.AddProjectRBACRoleResponse.displayName = 'proto.grpc.service.webservice.AddProjectRBACRoleResponse';
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
proto.grpc.service.webservice.UpdateProjectRBACRoleRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.repeatedFields_, null);
};
goog.inherits(proto.grpc.service.webservice.UpdateProjectRBACRoleRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.displayName = 'proto.grpc.service.webservice.UpdateProjectRBACRoleRequest';
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
proto.grpc.service.webservice.UpdateProjectRBACRoleResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.UpdateProjectRBACRoleResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.UpdateProjectRBACRoleResponse.displayName = 'proto.grpc.service.webservice.UpdateProjectRBACRoleResponse';
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
proto.grpc.service.webservice.DeleteProjectRBACRoleRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.DeleteProjectRBACRoleRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.DeleteProjectRBACRoleRequest.displayName = 'proto.grpc.service.webservice.DeleteProjectRBACRoleRequest';
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
proto.grpc.service.webservice.DeleteProjectRBACRoleResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.DeleteProjectRBACRoleResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.DeleteProjectRBACRoleResponse.displayName = 'proto.grpc.service.webservice.DeleteProjectRBACRoleResponse';
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
proto.grpc.service.webservice.AddProjectUserGroupRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.AddProjectUserGroupRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.AddProjectUserGroupRequest.displayName = 'proto.grpc.service.webservice.AddProjectUserGroupRequest';
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
proto.grpc.service.webservice.AddProjectUserGroupResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.AddProjectUserGroupResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.AddProjectUserGroupResponse.displayName = 'proto.grpc.service.webservice.AddProjectUserGroupResponse';
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
proto.grpc.service.webservice.DeleteProjectUserGroupRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.DeleteProjectUserGroupRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.DeleteProjectUserGroupRequest.displayName = 'proto.grpc.service.webservice.DeleteProjectUserGroupRequest';
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
proto.grpc.service.webservice.DeleteProjectUserGroupResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.DeleteProjectUserGroupResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.DeleteProjectUserGroupResponse.displayName = 'proto.grpc.service.webservice.DeleteProjectUserGroupResponse';
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
proto.grpc.service.webservice.GetCommandRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetCommandRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetCommandRequest.displayName = 'proto.grpc.service.webservice.GetCommandRequest';
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
proto.grpc.service.webservice.GetCommandResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetCommandResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetCommandResponse.displayName = 'proto.grpc.service.webservice.GetCommandResponse';
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
proto.grpc.service.webservice.GenerateAPIKeyRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GenerateAPIKeyRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GenerateAPIKeyRequest.displayName = 'proto.grpc.service.webservice.GenerateAPIKeyRequest';
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
proto.grpc.service.webservice.GenerateAPIKeyResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GenerateAPIKeyResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GenerateAPIKeyResponse.displayName = 'proto.grpc.service.webservice.GenerateAPIKeyResponse';
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
proto.grpc.service.webservice.DisableAPIKeyRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.DisableAPIKeyRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.DisableAPIKeyRequest.displayName = 'proto.grpc.service.webservice.DisableAPIKeyRequest';
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
proto.grpc.service.webservice.DisableAPIKeyResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.DisableAPIKeyResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.DisableAPIKeyResponse.displayName = 'proto.grpc.service.webservice.DisableAPIKeyResponse';
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
proto.grpc.service.webservice.ListAPIKeysRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.ListAPIKeysRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListAPIKeysRequest.displayName = 'proto.grpc.service.webservice.ListAPIKeysRequest';
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
proto.grpc.service.webservice.ListAPIKeysRequest.Options = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.ListAPIKeysRequest.Options, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListAPIKeysRequest.Options.displayName = 'proto.grpc.service.webservice.ListAPIKeysRequest.Options';
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
proto.grpc.service.webservice.ListAPIKeysResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.grpc.service.webservice.ListAPIKeysResponse.repeatedFields_, null);
};
goog.inherits(proto.grpc.service.webservice.ListAPIKeysResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListAPIKeysResponse.displayName = 'proto.grpc.service.webservice.ListAPIKeysResponse';
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
proto.grpc.service.webservice.GetInsightDataRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetInsightDataRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetInsightDataRequest.displayName = 'proto.grpc.service.webservice.GetInsightDataRequest';
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
proto.grpc.service.webservice.GetInsightDataResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.grpc.service.webservice.GetInsightDataResponse.repeatedFields_, null);
};
goog.inherits(proto.grpc.service.webservice.GetInsightDataResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetInsightDataResponse.displayName = 'proto.grpc.service.webservice.GetInsightDataResponse';
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
proto.grpc.service.webservice.GetInsightApplicationCountRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetInsightApplicationCountRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetInsightApplicationCountRequest.displayName = 'proto.grpc.service.webservice.GetInsightApplicationCountRequest';
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
proto.grpc.service.webservice.GetInsightApplicationCountResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.grpc.service.webservice.GetInsightApplicationCountResponse.repeatedFields_, null);
};
goog.inherits(proto.grpc.service.webservice.GetInsightApplicationCountResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetInsightApplicationCountResponse.displayName = 'proto.grpc.service.webservice.GetInsightApplicationCountResponse';
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
proto.grpc.service.webservice.ListDeploymentChainsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.ListDeploymentChainsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListDeploymentChainsRequest.displayName = 'proto.grpc.service.webservice.ListDeploymentChainsRequest';
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
proto.grpc.service.webservice.ListDeploymentChainsRequest.Options = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.ListDeploymentChainsRequest.Options, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListDeploymentChainsRequest.Options.displayName = 'proto.grpc.service.webservice.ListDeploymentChainsRequest.Options';
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
proto.grpc.service.webservice.ListDeploymentChainsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.grpc.service.webservice.ListDeploymentChainsResponse.repeatedFields_, null);
};
goog.inherits(proto.grpc.service.webservice.ListDeploymentChainsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListDeploymentChainsResponse.displayName = 'proto.grpc.service.webservice.ListDeploymentChainsResponse';
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
proto.grpc.service.webservice.GetDeploymentChainRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetDeploymentChainRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetDeploymentChainRequest.displayName = 'proto.grpc.service.webservice.GetDeploymentChainRequest';
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
proto.grpc.service.webservice.GetDeploymentChainResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.GetDeploymentChainResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.GetDeploymentChainResponse.displayName = 'proto.grpc.service.webservice.GetDeploymentChainResponse';
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
proto.grpc.service.webservice.ListEventsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.grpc.service.webservice.ListEventsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListEventsRequest.displayName = 'proto.grpc.service.webservice.ListEventsRequest';
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
proto.grpc.service.webservice.ListEventsRequest.Options = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.grpc.service.webservice.ListEventsRequest.Options.repeatedFields_, null);
};
goog.inherits(proto.grpc.service.webservice.ListEventsRequest.Options, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListEventsRequest.Options.displayName = 'proto.grpc.service.webservice.ListEventsRequest.Options';
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
proto.grpc.service.webservice.ListEventsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.grpc.service.webservice.ListEventsResponse.repeatedFields_, null);
};
goog.inherits(proto.grpc.service.webservice.ListEventsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.grpc.service.webservice.ListEventsResponse.displayName = 'proto.grpc.service.webservice.ListEventsResponse';
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
proto.grpc.service.webservice.RegisterPipedRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.RegisterPipedRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.RegisterPipedRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.RegisterPipedRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    desc: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.grpc.service.webservice.RegisterPipedRequest}
 */
proto.grpc.service.webservice.RegisterPipedRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.RegisterPipedRequest;
  return proto.grpc.service.webservice.RegisterPipedRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.RegisterPipedRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.RegisterPipedRequest}
 */
proto.grpc.service.webservice.RegisterPipedRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setDesc(value);
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
proto.grpc.service.webservice.RegisterPipedRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.RegisterPipedRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.RegisterPipedRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.RegisterPipedRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getDesc();
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
proto.grpc.service.webservice.RegisterPipedRequest.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.RegisterPipedRequest} returns this
 */
proto.grpc.service.webservice.RegisterPipedRequest.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string desc = 2;
 * @return {string}
 */
proto.grpc.service.webservice.RegisterPipedRequest.prototype.getDesc = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.RegisterPipedRequest} returns this
 */
proto.grpc.service.webservice.RegisterPipedRequest.prototype.setDesc = function(value) {
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
proto.grpc.service.webservice.RegisterPipedResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.RegisterPipedResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.RegisterPipedResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.RegisterPipedResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    key: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.grpc.service.webservice.RegisterPipedResponse}
 */
proto.grpc.service.webservice.RegisterPipedResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.RegisterPipedResponse;
  return proto.grpc.service.webservice.RegisterPipedResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.RegisterPipedResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.RegisterPipedResponse}
 */
proto.grpc.service.webservice.RegisterPipedResponse.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setKey(value);
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
proto.grpc.service.webservice.RegisterPipedResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.RegisterPipedResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.RegisterPipedResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.RegisterPipedResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getKey();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.RegisterPipedResponse.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.RegisterPipedResponse} returns this
 */
proto.grpc.service.webservice.RegisterPipedResponse.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string key = 2;
 * @return {string}
 */
proto.grpc.service.webservice.RegisterPipedResponse.prototype.getKey = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.RegisterPipedResponse} returns this
 */
proto.grpc.service.webservice.RegisterPipedResponse.prototype.setKey = function(value) {
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
proto.grpc.service.webservice.UpdatePipedRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.UpdatePipedRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.UpdatePipedRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdatePipedRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    pipedId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    desc: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.grpc.service.webservice.UpdatePipedRequest}
 */
proto.grpc.service.webservice.UpdatePipedRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.UpdatePipedRequest;
  return proto.grpc.service.webservice.UpdatePipedRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.UpdatePipedRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.UpdatePipedRequest}
 */
proto.grpc.service.webservice.UpdatePipedRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setPipedId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setDesc(value);
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
proto.grpc.service.webservice.UpdatePipedRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.UpdatePipedRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.UpdatePipedRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdatePipedRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPipedId();
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
};


/**
 * optional string piped_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.UpdatePipedRequest.prototype.getPipedId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.UpdatePipedRequest} returns this
 */
proto.grpc.service.webservice.UpdatePipedRequest.prototype.setPipedId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.grpc.service.webservice.UpdatePipedRequest.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.UpdatePipedRequest} returns this
 */
proto.grpc.service.webservice.UpdatePipedRequest.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string desc = 3;
 * @return {string}
 */
proto.grpc.service.webservice.UpdatePipedRequest.prototype.getDesc = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.UpdatePipedRequest} returns this
 */
proto.grpc.service.webservice.UpdatePipedRequest.prototype.setDesc = function(value) {
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
proto.grpc.service.webservice.UpdatePipedResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.UpdatePipedResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.UpdatePipedResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdatePipedResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.UpdatePipedResponse}
 */
proto.grpc.service.webservice.UpdatePipedResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.UpdatePipedResponse;
  return proto.grpc.service.webservice.UpdatePipedResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.UpdatePipedResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.UpdatePipedResponse}
 */
proto.grpc.service.webservice.UpdatePipedResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.UpdatePipedResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.UpdatePipedResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.UpdatePipedResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdatePipedResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.RecreatePipedKeyRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.RecreatePipedKeyRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.RecreatePipedKeyRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.RecreatePipedKeyRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.RecreatePipedKeyRequest}
 */
proto.grpc.service.webservice.RecreatePipedKeyRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.RecreatePipedKeyRequest;
  return proto.grpc.service.webservice.RecreatePipedKeyRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.RecreatePipedKeyRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.RecreatePipedKeyRequest}
 */
proto.grpc.service.webservice.RecreatePipedKeyRequest.deserializeBinaryFromReader = function(msg, reader) {
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
proto.grpc.service.webservice.RecreatePipedKeyRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.RecreatePipedKeyRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.RecreatePipedKeyRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.RecreatePipedKeyRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.RecreatePipedKeyRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.RecreatePipedKeyRequest} returns this
 */
proto.grpc.service.webservice.RecreatePipedKeyRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.RecreatePipedKeyResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.RecreatePipedKeyResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.RecreatePipedKeyResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.RecreatePipedKeyResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    key: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.RecreatePipedKeyResponse}
 */
proto.grpc.service.webservice.RecreatePipedKeyResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.RecreatePipedKeyResponse;
  return proto.grpc.service.webservice.RecreatePipedKeyResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.RecreatePipedKeyResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.RecreatePipedKeyResponse}
 */
proto.grpc.service.webservice.RecreatePipedKeyResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setKey(value);
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
proto.grpc.service.webservice.RecreatePipedKeyResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.RecreatePipedKeyResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.RecreatePipedKeyResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.RecreatePipedKeyResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getKey();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string key = 1;
 * @return {string}
 */
proto.grpc.service.webservice.RecreatePipedKeyResponse.prototype.getKey = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.RecreatePipedKeyResponse} returns this
 */
proto.grpc.service.webservice.RecreatePipedKeyResponse.prototype.setKey = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.DeleteOldPipedKeysRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.DeleteOldPipedKeysRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.DeleteOldPipedKeysRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DeleteOldPipedKeysRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    pipedId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.DeleteOldPipedKeysRequest}
 */
proto.grpc.service.webservice.DeleteOldPipedKeysRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.DeleteOldPipedKeysRequest;
  return proto.grpc.service.webservice.DeleteOldPipedKeysRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.DeleteOldPipedKeysRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.DeleteOldPipedKeysRequest}
 */
proto.grpc.service.webservice.DeleteOldPipedKeysRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setPipedId(value);
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
proto.grpc.service.webservice.DeleteOldPipedKeysRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.DeleteOldPipedKeysRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.DeleteOldPipedKeysRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DeleteOldPipedKeysRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPipedId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string piped_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.DeleteOldPipedKeysRequest.prototype.getPipedId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.DeleteOldPipedKeysRequest} returns this
 */
proto.grpc.service.webservice.DeleteOldPipedKeysRequest.prototype.setPipedId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.DeleteOldPipedKeysResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.DeleteOldPipedKeysResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.DeleteOldPipedKeysResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DeleteOldPipedKeysResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.DeleteOldPipedKeysResponse}
 */
proto.grpc.service.webservice.DeleteOldPipedKeysResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.DeleteOldPipedKeysResponse;
  return proto.grpc.service.webservice.DeleteOldPipedKeysResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.DeleteOldPipedKeysResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.DeleteOldPipedKeysResponse}
 */
proto.grpc.service.webservice.DeleteOldPipedKeysResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.DeleteOldPipedKeysResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.DeleteOldPipedKeysResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.DeleteOldPipedKeysResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DeleteOldPipedKeysResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.EnablePipedRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.EnablePipedRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.EnablePipedRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.EnablePipedRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    pipedId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.EnablePipedRequest}
 */
proto.grpc.service.webservice.EnablePipedRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.EnablePipedRequest;
  return proto.grpc.service.webservice.EnablePipedRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.EnablePipedRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.EnablePipedRequest}
 */
proto.grpc.service.webservice.EnablePipedRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setPipedId(value);
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
proto.grpc.service.webservice.EnablePipedRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.EnablePipedRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.EnablePipedRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.EnablePipedRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPipedId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string piped_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.EnablePipedRequest.prototype.getPipedId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.EnablePipedRequest} returns this
 */
proto.grpc.service.webservice.EnablePipedRequest.prototype.setPipedId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.EnablePipedResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.EnablePipedResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.EnablePipedResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.EnablePipedResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.EnablePipedResponse}
 */
proto.grpc.service.webservice.EnablePipedResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.EnablePipedResponse;
  return proto.grpc.service.webservice.EnablePipedResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.EnablePipedResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.EnablePipedResponse}
 */
proto.grpc.service.webservice.EnablePipedResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.EnablePipedResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.EnablePipedResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.EnablePipedResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.EnablePipedResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.DisablePipedRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.DisablePipedRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.DisablePipedRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DisablePipedRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    pipedId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.DisablePipedRequest}
 */
proto.grpc.service.webservice.DisablePipedRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.DisablePipedRequest;
  return proto.grpc.service.webservice.DisablePipedRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.DisablePipedRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.DisablePipedRequest}
 */
proto.grpc.service.webservice.DisablePipedRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setPipedId(value);
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
proto.grpc.service.webservice.DisablePipedRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.DisablePipedRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.DisablePipedRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DisablePipedRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPipedId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string piped_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.DisablePipedRequest.prototype.getPipedId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.DisablePipedRequest} returns this
 */
proto.grpc.service.webservice.DisablePipedRequest.prototype.setPipedId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.DisablePipedResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.DisablePipedResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.DisablePipedResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DisablePipedResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.DisablePipedResponse}
 */
proto.grpc.service.webservice.DisablePipedResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.DisablePipedResponse;
  return proto.grpc.service.webservice.DisablePipedResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.DisablePipedResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.DisablePipedResponse}
 */
proto.grpc.service.webservice.DisablePipedResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.DisablePipedResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.DisablePipedResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.DisablePipedResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DisablePipedResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.ListPipedsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListPipedsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListPipedsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListPipedsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    withStatus: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    options: (f = msg.getOptions()) && proto.grpc.service.webservice.ListPipedsRequest.Options.toObject(includeInstance, f)
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
 * @return {!proto.grpc.service.webservice.ListPipedsRequest}
 */
proto.grpc.service.webservice.ListPipedsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListPipedsRequest;
  return proto.grpc.service.webservice.ListPipedsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListPipedsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListPipedsRequest}
 */
proto.grpc.service.webservice.ListPipedsRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setWithStatus(value);
      break;
    case 2:
      var value = new proto.grpc.service.webservice.ListPipedsRequest.Options;
      reader.readMessage(value,proto.grpc.service.webservice.ListPipedsRequest.Options.deserializeBinaryFromReader);
      msg.setOptions(value);
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
proto.grpc.service.webservice.ListPipedsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListPipedsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListPipedsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListPipedsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWithStatus();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getOptions();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.grpc.service.webservice.ListPipedsRequest.Options.serializeBinaryToWriter
    );
  }
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
proto.grpc.service.webservice.ListPipedsRequest.Options.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListPipedsRequest.Options.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListPipedsRequest.Options} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListPipedsRequest.Options.toObject = function(includeInstance, msg) {
  var f, obj = {
    enabled: (f = msg.getEnabled()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f)
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
 * @return {!proto.grpc.service.webservice.ListPipedsRequest.Options}
 */
proto.grpc.service.webservice.ListPipedsRequest.Options.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListPipedsRequest.Options;
  return proto.grpc.service.webservice.ListPipedsRequest.Options.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListPipedsRequest.Options} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListPipedsRequest.Options}
 */
proto.grpc.service.webservice.ListPipedsRequest.Options.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setEnabled(value);
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
proto.grpc.service.webservice.ListPipedsRequest.Options.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListPipedsRequest.Options.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListPipedsRequest.Options} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListPipedsRequest.Options.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEnabled();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
};


/**
 * optional google.protobuf.BoolValue enabled = 1;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.grpc.service.webservice.ListPipedsRequest.Options.prototype.getEnabled = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 1));
};


/**
 * @param {?proto.google.protobuf.BoolValue|undefined} value
 * @return {!proto.grpc.service.webservice.ListPipedsRequest.Options} returns this
*/
proto.grpc.service.webservice.ListPipedsRequest.Options.prototype.setEnabled = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.ListPipedsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListPipedsRequest.Options.prototype.clearEnabled = function() {
  return this.setEnabled(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.ListPipedsRequest.Options.prototype.hasEnabled = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional bool with_status = 1;
 * @return {boolean}
 */
proto.grpc.service.webservice.ListPipedsRequest.prototype.getWithStatus = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.grpc.service.webservice.ListPipedsRequest} returns this
 */
proto.grpc.service.webservice.ListPipedsRequest.prototype.setWithStatus = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional Options options = 2;
 * @return {?proto.grpc.service.webservice.ListPipedsRequest.Options}
 */
proto.grpc.service.webservice.ListPipedsRequest.prototype.getOptions = function() {
  return /** @type{?proto.grpc.service.webservice.ListPipedsRequest.Options} */ (
    jspb.Message.getWrapperField(this, proto.grpc.service.webservice.ListPipedsRequest.Options, 2));
};


/**
 * @param {?proto.grpc.service.webservice.ListPipedsRequest.Options|undefined} value
 * @return {!proto.grpc.service.webservice.ListPipedsRequest} returns this
*/
proto.grpc.service.webservice.ListPipedsRequest.prototype.setOptions = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.ListPipedsRequest} returns this
 */
proto.grpc.service.webservice.ListPipedsRequest.prototype.clearOptions = function() {
  return this.setOptions(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.ListPipedsRequest.prototype.hasOptions = function() {
  return jspb.Message.getField(this, 2) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.grpc.service.webservice.ListPipedsResponse.repeatedFields_ = [1];



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
proto.grpc.service.webservice.ListPipedsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListPipedsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListPipedsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListPipedsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    pipedsList: jspb.Message.toObjectList(msg.getPipedsList(),
    pkg_model_piped_pb.Piped.toObject, includeInstance)
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
 * @return {!proto.grpc.service.webservice.ListPipedsResponse}
 */
proto.grpc.service.webservice.ListPipedsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListPipedsResponse;
  return proto.grpc.service.webservice.ListPipedsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListPipedsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListPipedsResponse}
 */
proto.grpc.service.webservice.ListPipedsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new pkg_model_piped_pb.Piped;
      reader.readMessage(value,pkg_model_piped_pb.Piped.deserializeBinaryFromReader);
      msg.addPipeds(value);
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
proto.grpc.service.webservice.ListPipedsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListPipedsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListPipedsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListPipedsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPipedsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      pkg_model_piped_pb.Piped.serializeBinaryToWriter
    );
  }
};


/**
 * repeated model.Piped pipeds = 1;
 * @return {!Array<!proto.model.Piped>}
 */
proto.grpc.service.webservice.ListPipedsResponse.prototype.getPipedsList = function() {
  return /** @type{!Array<!proto.model.Piped>} */ (
    jspb.Message.getRepeatedWrapperField(this, pkg_model_piped_pb.Piped, 1));
};


/**
 * @param {!Array<!proto.model.Piped>} value
 * @return {!proto.grpc.service.webservice.ListPipedsResponse} returns this
*/
proto.grpc.service.webservice.ListPipedsResponse.prototype.setPipedsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.model.Piped=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.Piped}
 */
proto.grpc.service.webservice.ListPipedsResponse.prototype.addPipeds = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.model.Piped, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.ListPipedsResponse} returns this
 */
proto.grpc.service.webservice.ListPipedsResponse.prototype.clearPipedsList = function() {
  return this.setPipedsList([]);
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
proto.grpc.service.webservice.GetPipedRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetPipedRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetPipedRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetPipedRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    pipedId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.GetPipedRequest}
 */
proto.grpc.service.webservice.GetPipedRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetPipedRequest;
  return proto.grpc.service.webservice.GetPipedRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetPipedRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetPipedRequest}
 */
proto.grpc.service.webservice.GetPipedRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setPipedId(value);
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
proto.grpc.service.webservice.GetPipedRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetPipedRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetPipedRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetPipedRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPipedId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string piped_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.GetPipedRequest.prototype.getPipedId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.GetPipedRequest} returns this
 */
proto.grpc.service.webservice.GetPipedRequest.prototype.setPipedId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.GetPipedResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetPipedResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetPipedResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetPipedResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    piped: (f = msg.getPiped()) && pkg_model_piped_pb.Piped.toObject(includeInstance, f)
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
 * @return {!proto.grpc.service.webservice.GetPipedResponse}
 */
proto.grpc.service.webservice.GetPipedResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetPipedResponse;
  return proto.grpc.service.webservice.GetPipedResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetPipedResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetPipedResponse}
 */
proto.grpc.service.webservice.GetPipedResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new pkg_model_piped_pb.Piped;
      reader.readMessage(value,pkg_model_piped_pb.Piped.deserializeBinaryFromReader);
      msg.setPiped(value);
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
proto.grpc.service.webservice.GetPipedResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetPipedResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetPipedResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetPipedResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPiped();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      pkg_model_piped_pb.Piped.serializeBinaryToWriter
    );
  }
};


/**
 * optional model.Piped piped = 1;
 * @return {?proto.model.Piped}
 */
proto.grpc.service.webservice.GetPipedResponse.prototype.getPiped = function() {
  return /** @type{?proto.model.Piped} */ (
    jspb.Message.getWrapperField(this, pkg_model_piped_pb.Piped, 1));
};


/**
 * @param {?proto.model.Piped|undefined} value
 * @return {!proto.grpc.service.webservice.GetPipedResponse} returns this
*/
proto.grpc.service.webservice.GetPipedResponse.prototype.setPiped = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.GetPipedResponse} returns this
 */
proto.grpc.service.webservice.GetPipedResponse.prototype.clearPiped = function() {
  return this.setPiped(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.GetPipedResponse.prototype.hasPiped = function() {
  return jspb.Message.getField(this, 1) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.repeatedFields_ = [2];



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
proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    version: jspb.Message.getFieldWithDefault(msg, 1, ""),
    pipedIdsList: (f = jspb.Message.getRepeatedField(msg, 2)) == null ? undefined : f
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
 * @return {!proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest}
 */
proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest;
  return proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest}
 */
proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setVersion(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.addPipedIds(value);
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
proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getVersion();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPipedIdsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      2,
      f
    );
  }
};


/**
 * optional string version = 1;
 * @return {string}
 */
proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.prototype.getVersion = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest} returns this
 */
proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.prototype.setVersion = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * repeated string piped_ids = 2;
 * @return {!Array<string>}
 */
proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.prototype.getPipedIdsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 2));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest} returns this
 */
proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.prototype.setPipedIdsList = function(value) {
  return jspb.Message.setField(this, 2, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest} returns this
 */
proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.prototype.addPipedIds = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 2, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest} returns this
 */
proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest.prototype.clearPipedIdsList = function() {
  return this.setPipedIdsList([]);
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
proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse}
 */
proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse;
  return proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse}
 */
proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.RestartPipedRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.RestartPipedRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.RestartPipedRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.RestartPipedRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    pipedId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.RestartPipedRequest}
 */
proto.grpc.service.webservice.RestartPipedRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.RestartPipedRequest;
  return proto.grpc.service.webservice.RestartPipedRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.RestartPipedRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.RestartPipedRequest}
 */
proto.grpc.service.webservice.RestartPipedRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setPipedId(value);
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
proto.grpc.service.webservice.RestartPipedRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.RestartPipedRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.RestartPipedRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.RestartPipedRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPipedId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string piped_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.RestartPipedRequest.prototype.getPipedId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.RestartPipedRequest} returns this
 */
proto.grpc.service.webservice.RestartPipedRequest.prototype.setPipedId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.RestartPipedResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.RestartPipedResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.RestartPipedResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.RestartPipedResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    commandId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.RestartPipedResponse}
 */
proto.grpc.service.webservice.RestartPipedResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.RestartPipedResponse;
  return proto.grpc.service.webservice.RestartPipedResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.RestartPipedResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.RestartPipedResponse}
 */
proto.grpc.service.webservice.RestartPipedResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setCommandId(value);
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
proto.grpc.service.webservice.RestartPipedResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.RestartPipedResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.RestartPipedResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.RestartPipedResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCommandId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string command_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.RestartPipedResponse.prototype.getCommandId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.RestartPipedResponse} returns this
 */
proto.grpc.service.webservice.RestartPipedResponse.prototype.setCommandId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.ListReleasedVersionsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListReleasedVersionsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListReleasedVersionsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListReleasedVersionsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.ListReleasedVersionsRequest}
 */
proto.grpc.service.webservice.ListReleasedVersionsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListReleasedVersionsRequest;
  return proto.grpc.service.webservice.ListReleasedVersionsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListReleasedVersionsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListReleasedVersionsRequest}
 */
proto.grpc.service.webservice.ListReleasedVersionsRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.ListReleasedVersionsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListReleasedVersionsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListReleasedVersionsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListReleasedVersionsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.grpc.service.webservice.ListReleasedVersionsResponse.repeatedFields_ = [1];



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
proto.grpc.service.webservice.ListReleasedVersionsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListReleasedVersionsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListReleasedVersionsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListReleasedVersionsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    versionsList: (f = jspb.Message.getRepeatedField(msg, 1)) == null ? undefined : f
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
 * @return {!proto.grpc.service.webservice.ListReleasedVersionsResponse}
 */
proto.grpc.service.webservice.ListReleasedVersionsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListReleasedVersionsResponse;
  return proto.grpc.service.webservice.ListReleasedVersionsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListReleasedVersionsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListReleasedVersionsResponse}
 */
proto.grpc.service.webservice.ListReleasedVersionsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.addVersions(value);
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
proto.grpc.service.webservice.ListReleasedVersionsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListReleasedVersionsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListReleasedVersionsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListReleasedVersionsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getVersionsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      1,
      f
    );
  }
};


/**
 * repeated string versions = 1;
 * @return {!Array<string>}
 */
proto.grpc.service.webservice.ListReleasedVersionsResponse.prototype.getVersionsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 1));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.grpc.service.webservice.ListReleasedVersionsResponse} returns this
 */
proto.grpc.service.webservice.ListReleasedVersionsResponse.prototype.setVersionsList = function(value) {
  return jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.grpc.service.webservice.ListReleasedVersionsResponse} returns this
 */
proto.grpc.service.webservice.ListReleasedVersionsResponse.prototype.addVersions = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.ListReleasedVersionsResponse} returns this
 */
proto.grpc.service.webservice.ListReleasedVersionsResponse.prototype.clearVersionsList = function() {
  return this.setVersionsList([]);
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
proto.grpc.service.webservice.ListDeprecatedNotesRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListDeprecatedNotesRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListDeprecatedNotesRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListDeprecatedNotesRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    projectId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.ListDeprecatedNotesRequest}
 */
proto.grpc.service.webservice.ListDeprecatedNotesRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListDeprecatedNotesRequest;
  return proto.grpc.service.webservice.ListDeprecatedNotesRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListDeprecatedNotesRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListDeprecatedNotesRequest}
 */
proto.grpc.service.webservice.ListDeprecatedNotesRequest.deserializeBinaryFromReader = function(msg, reader) {
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
proto.grpc.service.webservice.ListDeprecatedNotesRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListDeprecatedNotesRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListDeprecatedNotesRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListDeprecatedNotesRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getProjectId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string project_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.ListDeprecatedNotesRequest.prototype.getProjectId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.ListDeprecatedNotesRequest} returns this
 */
proto.grpc.service.webservice.ListDeprecatedNotesRequest.prototype.setProjectId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.ListDeprecatedNotesResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListDeprecatedNotesResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListDeprecatedNotesResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListDeprecatedNotesResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    notes: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.ListDeprecatedNotesResponse}
 */
proto.grpc.service.webservice.ListDeprecatedNotesResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListDeprecatedNotesResponse;
  return proto.grpc.service.webservice.ListDeprecatedNotesResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListDeprecatedNotesResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListDeprecatedNotesResponse}
 */
proto.grpc.service.webservice.ListDeprecatedNotesResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setNotes(value);
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
proto.grpc.service.webservice.ListDeprecatedNotesResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListDeprecatedNotesResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListDeprecatedNotesResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListDeprecatedNotesResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getNotes();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string notes = 1;
 * @return {string}
 */
proto.grpc.service.webservice.ListDeprecatedNotesResponse.prototype.getNotes = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.ListDeprecatedNotesResponse} returns this
 */
proto.grpc.service.webservice.ListDeprecatedNotesResponse.prototype.setNotes = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.AddApplicationRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.AddApplicationRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.AddApplicationRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.AddApplicationRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    pipedId: jspb.Message.getFieldWithDefault(msg, 3, ""),
    gitPath: (f = msg.getGitPath()) && pkg_model_common_pb.ApplicationGitPath.toObject(includeInstance, f),
    kind: jspb.Message.getFieldWithDefault(msg, 5, 0),
    platformProvider: jspb.Message.getFieldWithDefault(msg, 9, ""),
    deployTargetsByPluginMap: (f = msg.getDeployTargetsByPluginMap()) ? f.toObject(includeInstance, proto.model.DeployTargets.toObject) : [],
    description: jspb.Message.getFieldWithDefault(msg, 7, ""),
    labelsMap: (f = msg.getLabelsMap()) ? f.toObject(includeInstance, undefined) : []
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
 * @return {!proto.grpc.service.webservice.AddApplicationRequest}
 */
proto.grpc.service.webservice.AddApplicationRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.AddApplicationRequest;
  return proto.grpc.service.webservice.AddApplicationRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.AddApplicationRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.AddApplicationRequest}
 */
proto.grpc.service.webservice.AddApplicationRequest.deserializeBinaryFromReader = function(msg, reader) {
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
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setPipedId(value);
      break;
    case 4:
      var value = new pkg_model_common_pb.ApplicationGitPath;
      reader.readMessage(value,pkg_model_common_pb.ApplicationGitPath.deserializeBinaryFromReader);
      msg.setGitPath(value);
      break;
    case 5:
      var value = /** @type {!proto.model.ApplicationKind} */ (reader.readEnum());
      msg.setKind(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setPlatformProvider(value);
      break;
    case 10:
      var value = msg.getDeployTargetsByPluginMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readMessage, proto.model.DeployTargets.deserializeBinaryFromReader, "", new proto.model.DeployTargets());
         });
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescription(value);
      break;
    case 8:
      var value = msg.getLabelsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "", "");
         });
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
proto.grpc.service.webservice.AddApplicationRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.AddApplicationRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.AddApplicationRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.AddApplicationRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPipedId();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getGitPath();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      pkg_model_common_pb.ApplicationGitPath.serializeBinaryToWriter
    );
  }
  f = message.getKind();
  if (f !== 0.0) {
    writer.writeEnum(
      5,
      f
    );
  }
  f = message.getPlatformProvider();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
  f = message.getDeployTargetsByPluginMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(10, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeMessage, proto.model.DeployTargets.serializeBinaryToWriter);
  }
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
  f = message.getLabelsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(8, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.grpc.service.webservice.AddApplicationRequest.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.AddApplicationRequest} returns this
 */
proto.grpc.service.webservice.AddApplicationRequest.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string piped_id = 3;
 * @return {string}
 */
proto.grpc.service.webservice.AddApplicationRequest.prototype.getPipedId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.AddApplicationRequest} returns this
 */
proto.grpc.service.webservice.AddApplicationRequest.prototype.setPipedId = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional model.ApplicationGitPath git_path = 4;
 * @return {?proto.model.ApplicationGitPath}
 */
proto.grpc.service.webservice.AddApplicationRequest.prototype.getGitPath = function() {
  return /** @type{?proto.model.ApplicationGitPath} */ (
    jspb.Message.getWrapperField(this, pkg_model_common_pb.ApplicationGitPath, 4));
};


/**
 * @param {?proto.model.ApplicationGitPath|undefined} value
 * @return {!proto.grpc.service.webservice.AddApplicationRequest} returns this
*/
proto.grpc.service.webservice.AddApplicationRequest.prototype.setGitPath = function(value) {
  return jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.AddApplicationRequest} returns this
 */
proto.grpc.service.webservice.AddApplicationRequest.prototype.clearGitPath = function() {
  return this.setGitPath(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.AddApplicationRequest.prototype.hasGitPath = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional model.ApplicationKind kind = 5;
 * @return {!proto.model.ApplicationKind}
 */
proto.grpc.service.webservice.AddApplicationRequest.prototype.getKind = function() {
  return /** @type {!proto.model.ApplicationKind} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {!proto.model.ApplicationKind} value
 * @return {!proto.grpc.service.webservice.AddApplicationRequest} returns this
 */
proto.grpc.service.webservice.AddApplicationRequest.prototype.setKind = function(value) {
  return jspb.Message.setProto3EnumField(this, 5, value);
};


/**
 * optional string platform_provider = 9;
 * @return {string}
 */
proto.grpc.service.webservice.AddApplicationRequest.prototype.getPlatformProvider = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.AddApplicationRequest} returns this
 */
proto.grpc.service.webservice.AddApplicationRequest.prototype.setPlatformProvider = function(value) {
  return jspb.Message.setProto3StringField(this, 9, value);
};


/**
 * map<string, model.DeployTargets> deploy_targets_by_plugin = 10;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,!proto.model.DeployTargets>}
 */
proto.grpc.service.webservice.AddApplicationRequest.prototype.getDeployTargetsByPluginMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,!proto.model.DeployTargets>} */ (
      jspb.Message.getMapField(this, 10, opt_noLazyCreate,
      proto.model.DeployTargets));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.grpc.service.webservice.AddApplicationRequest} returns this
 */
proto.grpc.service.webservice.AddApplicationRequest.prototype.clearDeployTargetsByPluginMap = function() {
  this.getDeployTargetsByPluginMap().clear();
  return this;
};


/**
 * optional string description = 7;
 * @return {string}
 */
proto.grpc.service.webservice.AddApplicationRequest.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.AddApplicationRequest} returns this
 */
proto.grpc.service.webservice.AddApplicationRequest.prototype.setDescription = function(value) {
  return jspb.Message.setProto3StringField(this, 7, value);
};


/**
 * map<string, string> labels = 8;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.grpc.service.webservice.AddApplicationRequest.prototype.getLabelsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 8, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.grpc.service.webservice.AddApplicationRequest} returns this
 */
proto.grpc.service.webservice.AddApplicationRequest.prototype.clearLabelsMap = function() {
  this.getLabelsMap().clear();
  return this;
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
proto.grpc.service.webservice.AddApplicationResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.AddApplicationResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.AddApplicationResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.AddApplicationResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    applicationId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.AddApplicationResponse}
 */
proto.grpc.service.webservice.AddApplicationResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.AddApplicationResponse;
  return proto.grpc.service.webservice.AddApplicationResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.AddApplicationResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.AddApplicationResponse}
 */
proto.grpc.service.webservice.AddApplicationResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.grpc.service.webservice.AddApplicationResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.AddApplicationResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.AddApplicationResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.AddApplicationResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplicationId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string application_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.AddApplicationResponse.prototype.getApplicationId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.AddApplicationResponse} returns this
 */
proto.grpc.service.webservice.AddApplicationResponse.prototype.setApplicationId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.UpdateApplicationRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.UpdateApplicationRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.UpdateApplicationRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateApplicationRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    applicationId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    pipedId: jspb.Message.getFieldWithDefault(msg, 4, ""),
    kind: jspb.Message.getFieldWithDefault(msg, 6, 0),
    platformProvider: jspb.Message.getFieldWithDefault(msg, 9, ""),
    deployTargetsByPluginMap: (f = msg.getDeployTargetsByPluginMap()) ? f.toObject(includeInstance, proto.model.DeployTargets.toObject) : [],
    configFilename: jspb.Message.getFieldWithDefault(msg, 8, "")
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
 * @return {!proto.grpc.service.webservice.UpdateApplicationRequest}
 */
proto.grpc.service.webservice.UpdateApplicationRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.UpdateApplicationRequest;
  return proto.grpc.service.webservice.UpdateApplicationRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.UpdateApplicationRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.UpdateApplicationRequest}
 */
proto.grpc.service.webservice.UpdateApplicationRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setName(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setPipedId(value);
      break;
    case 6:
      var value = /** @type {!proto.model.ApplicationKind} */ (reader.readEnum());
      msg.setKind(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setPlatformProvider(value);
      break;
    case 10:
      var value = msg.getDeployTargetsByPluginMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readMessage, proto.model.DeployTargets.deserializeBinaryFromReader, "", new proto.model.DeployTargets());
         });
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.setConfigFilename(value);
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
proto.grpc.service.webservice.UpdateApplicationRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.UpdateApplicationRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.UpdateApplicationRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateApplicationRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplicationId();
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
  f = message.getPipedId();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getKind();
  if (f !== 0.0) {
    writer.writeEnum(
      6,
      f
    );
  }
  f = message.getPlatformProvider();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
  f = message.getDeployTargetsByPluginMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(10, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeMessage, proto.model.DeployTargets.serializeBinaryToWriter);
  }
  f = message.getConfigFilename();
  if (f.length > 0) {
    writer.writeString(
      8,
      f
    );
  }
};


/**
 * optional string application_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.UpdateApplicationRequest.prototype.getApplicationId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.UpdateApplicationRequest} returns this
 */
proto.grpc.service.webservice.UpdateApplicationRequest.prototype.setApplicationId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.grpc.service.webservice.UpdateApplicationRequest.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.UpdateApplicationRequest} returns this
 */
proto.grpc.service.webservice.UpdateApplicationRequest.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string piped_id = 4;
 * @return {string}
 */
proto.grpc.service.webservice.UpdateApplicationRequest.prototype.getPipedId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.UpdateApplicationRequest} returns this
 */
proto.grpc.service.webservice.UpdateApplicationRequest.prototype.setPipedId = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional model.ApplicationKind kind = 6;
 * @return {!proto.model.ApplicationKind}
 */
proto.grpc.service.webservice.UpdateApplicationRequest.prototype.getKind = function() {
  return /** @type {!proto.model.ApplicationKind} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {!proto.model.ApplicationKind} value
 * @return {!proto.grpc.service.webservice.UpdateApplicationRequest} returns this
 */
proto.grpc.service.webservice.UpdateApplicationRequest.prototype.setKind = function(value) {
  return jspb.Message.setProto3EnumField(this, 6, value);
};


/**
 * optional string platform_provider = 9;
 * @return {string}
 */
proto.grpc.service.webservice.UpdateApplicationRequest.prototype.getPlatformProvider = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.UpdateApplicationRequest} returns this
 */
proto.grpc.service.webservice.UpdateApplicationRequest.prototype.setPlatformProvider = function(value) {
  return jspb.Message.setProto3StringField(this, 9, value);
};


/**
 * map<string, model.DeployTargets> deploy_targets_by_plugin = 10;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,!proto.model.DeployTargets>}
 */
proto.grpc.service.webservice.UpdateApplicationRequest.prototype.getDeployTargetsByPluginMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,!proto.model.DeployTargets>} */ (
      jspb.Message.getMapField(this, 10, opt_noLazyCreate,
      proto.model.DeployTargets));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.grpc.service.webservice.UpdateApplicationRequest} returns this
 */
proto.grpc.service.webservice.UpdateApplicationRequest.prototype.clearDeployTargetsByPluginMap = function() {
  this.getDeployTargetsByPluginMap().clear();
  return this;
};


/**
 * optional string config_filename = 8;
 * @return {string}
 */
proto.grpc.service.webservice.UpdateApplicationRequest.prototype.getConfigFilename = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.UpdateApplicationRequest} returns this
 */
proto.grpc.service.webservice.UpdateApplicationRequest.prototype.setConfigFilename = function(value) {
  return jspb.Message.setProto3StringField(this, 8, value);
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
proto.grpc.service.webservice.UpdateApplicationResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.UpdateApplicationResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.UpdateApplicationResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateApplicationResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.UpdateApplicationResponse}
 */
proto.grpc.service.webservice.UpdateApplicationResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.UpdateApplicationResponse;
  return proto.grpc.service.webservice.UpdateApplicationResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.UpdateApplicationResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.UpdateApplicationResponse}
 */
proto.grpc.service.webservice.UpdateApplicationResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.UpdateApplicationResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.UpdateApplicationResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.UpdateApplicationResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateApplicationResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.EnableApplicationRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.EnableApplicationRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.EnableApplicationRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.EnableApplicationRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    applicationId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.EnableApplicationRequest}
 */
proto.grpc.service.webservice.EnableApplicationRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.EnableApplicationRequest;
  return proto.grpc.service.webservice.EnableApplicationRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.EnableApplicationRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.EnableApplicationRequest}
 */
proto.grpc.service.webservice.EnableApplicationRequest.deserializeBinaryFromReader = function(msg, reader) {
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
proto.grpc.service.webservice.EnableApplicationRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.EnableApplicationRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.EnableApplicationRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.EnableApplicationRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplicationId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string application_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.EnableApplicationRequest.prototype.getApplicationId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.EnableApplicationRequest} returns this
 */
proto.grpc.service.webservice.EnableApplicationRequest.prototype.setApplicationId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.EnableApplicationResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.EnableApplicationResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.EnableApplicationResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.EnableApplicationResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.EnableApplicationResponse}
 */
proto.grpc.service.webservice.EnableApplicationResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.EnableApplicationResponse;
  return proto.grpc.service.webservice.EnableApplicationResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.EnableApplicationResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.EnableApplicationResponse}
 */
proto.grpc.service.webservice.EnableApplicationResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.EnableApplicationResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.EnableApplicationResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.EnableApplicationResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.EnableApplicationResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.DisableApplicationRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.DisableApplicationRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.DisableApplicationRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DisableApplicationRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    applicationId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.DisableApplicationRequest}
 */
proto.grpc.service.webservice.DisableApplicationRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.DisableApplicationRequest;
  return proto.grpc.service.webservice.DisableApplicationRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.DisableApplicationRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.DisableApplicationRequest}
 */
proto.grpc.service.webservice.DisableApplicationRequest.deserializeBinaryFromReader = function(msg, reader) {
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
proto.grpc.service.webservice.DisableApplicationRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.DisableApplicationRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.DisableApplicationRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DisableApplicationRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplicationId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string application_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.DisableApplicationRequest.prototype.getApplicationId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.DisableApplicationRequest} returns this
 */
proto.grpc.service.webservice.DisableApplicationRequest.prototype.setApplicationId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.DisableApplicationResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.DisableApplicationResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.DisableApplicationResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DisableApplicationResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.DisableApplicationResponse}
 */
proto.grpc.service.webservice.DisableApplicationResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.DisableApplicationResponse;
  return proto.grpc.service.webservice.DisableApplicationResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.DisableApplicationResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.DisableApplicationResponse}
 */
proto.grpc.service.webservice.DisableApplicationResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.DisableApplicationResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.DisableApplicationResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.DisableApplicationResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DisableApplicationResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.DeleteApplicationRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.DeleteApplicationRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.DeleteApplicationRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DeleteApplicationRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    applicationId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.DeleteApplicationRequest}
 */
proto.grpc.service.webservice.DeleteApplicationRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.DeleteApplicationRequest;
  return proto.grpc.service.webservice.DeleteApplicationRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.DeleteApplicationRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.DeleteApplicationRequest}
 */
proto.grpc.service.webservice.DeleteApplicationRequest.deserializeBinaryFromReader = function(msg, reader) {
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
proto.grpc.service.webservice.DeleteApplicationRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.DeleteApplicationRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.DeleteApplicationRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DeleteApplicationRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplicationId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string application_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.DeleteApplicationRequest.prototype.getApplicationId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.DeleteApplicationRequest} returns this
 */
proto.grpc.service.webservice.DeleteApplicationRequest.prototype.setApplicationId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.DeleteApplicationResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.DeleteApplicationResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.DeleteApplicationResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DeleteApplicationResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.DeleteApplicationResponse}
 */
proto.grpc.service.webservice.DeleteApplicationResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.DeleteApplicationResponse;
  return proto.grpc.service.webservice.DeleteApplicationResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.DeleteApplicationResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.DeleteApplicationResponse}
 */
proto.grpc.service.webservice.DeleteApplicationResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.DeleteApplicationResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.DeleteApplicationResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.DeleteApplicationResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DeleteApplicationResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.ListApplicationsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListApplicationsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListApplicationsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListApplicationsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    options: (f = msg.getOptions()) && proto.grpc.service.webservice.ListApplicationsRequest.Options.toObject(includeInstance, f)
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
 * @return {!proto.grpc.service.webservice.ListApplicationsRequest}
 */
proto.grpc.service.webservice.ListApplicationsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListApplicationsRequest;
  return proto.grpc.service.webservice.ListApplicationsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListApplicationsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListApplicationsRequest}
 */
proto.grpc.service.webservice.ListApplicationsRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.grpc.service.webservice.ListApplicationsRequest.Options;
      reader.readMessage(value,proto.grpc.service.webservice.ListApplicationsRequest.Options.deserializeBinaryFromReader);
      msg.setOptions(value);
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
proto.grpc.service.webservice.ListApplicationsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListApplicationsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListApplicationsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListApplicationsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOptions();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.grpc.service.webservice.ListApplicationsRequest.Options.serializeBinaryToWriter
    );
  }
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.repeatedFields_ = [2,3];



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
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListApplicationsRequest.Options.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListApplicationsRequest.Options} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.toObject = function(includeInstance, msg) {
  var f, obj = {
    enabled: (f = msg.getEnabled()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
    kindsList: (f = jspb.Message.getRepeatedField(msg, 2)) == null ? undefined : f,
    syncStatusesList: (f = jspb.Message.getRepeatedField(msg, 3)) == null ? undefined : f,
    name: jspb.Message.getFieldWithDefault(msg, 5, ""),
    labelsMap: (f = msg.getLabelsMap()) ? f.toObject(includeInstance, undefined) : [],
    pipedId: jspb.Message.getFieldWithDefault(msg, 7, "")
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
 * @return {!proto.grpc.service.webservice.ListApplicationsRequest.Options}
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListApplicationsRequest.Options;
  return proto.grpc.service.webservice.ListApplicationsRequest.Options.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListApplicationsRequest.Options} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListApplicationsRequest.Options}
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setEnabled(value);
      break;
    case 2:
      var values = /** @type {!Array<!proto.model.ApplicationKind>} */ (reader.isDelimited() ? reader.readPackedEnum() : [reader.readEnum()]);
      for (var i = 0; i < values.length; i++) {
        msg.addKinds(values[i]);
      }
      break;
    case 3:
      var values = /** @type {!Array<!proto.model.ApplicationSyncStatus>} */ (reader.isDelimited() ? reader.readPackedEnum() : [reader.readEnum()]);
      for (var i = 0; i < values.length; i++) {
        msg.addSyncStatuses(values[i]);
      }
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 6:
      var value = msg.getLabelsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "", "");
         });
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setPipedId(value);
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
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListApplicationsRequest.Options.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListApplicationsRequest.Options} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEnabled();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
  f = message.getKindsList();
  if (f.length > 0) {
    writer.writePackedEnum(
      2,
      f
    );
  }
  f = message.getSyncStatusesList();
  if (f.length > 0) {
    writer.writePackedEnum(
      3,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getLabelsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(6, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
  f = message.getPipedId();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
};


/**
 * optional google.protobuf.BoolValue enabled = 1;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.getEnabled = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 1));
};


/**
 * @param {?proto.google.protobuf.BoolValue|undefined} value
 * @return {!proto.grpc.service.webservice.ListApplicationsRequest.Options} returns this
*/
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.setEnabled = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.ListApplicationsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.clearEnabled = function() {
  return this.setEnabled(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.hasEnabled = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * repeated model.ApplicationKind kinds = 2;
 * @return {!Array<!proto.model.ApplicationKind>}
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.getKindsList = function() {
  return /** @type {!Array<!proto.model.ApplicationKind>} */ (jspb.Message.getRepeatedField(this, 2));
};


/**
 * @param {!Array<!proto.model.ApplicationKind>} value
 * @return {!proto.grpc.service.webservice.ListApplicationsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.setKindsList = function(value) {
  return jspb.Message.setField(this, 2, value || []);
};


/**
 * @param {!proto.model.ApplicationKind} value
 * @param {number=} opt_index
 * @return {!proto.grpc.service.webservice.ListApplicationsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.addKinds = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 2, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.ListApplicationsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.clearKindsList = function() {
  return this.setKindsList([]);
};


/**
 * repeated model.ApplicationSyncStatus sync_statuses = 3;
 * @return {!Array<!proto.model.ApplicationSyncStatus>}
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.getSyncStatusesList = function() {
  return /** @type {!Array<!proto.model.ApplicationSyncStatus>} */ (jspb.Message.getRepeatedField(this, 3));
};


/**
 * @param {!Array<!proto.model.ApplicationSyncStatus>} value
 * @return {!proto.grpc.service.webservice.ListApplicationsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.setSyncStatusesList = function(value) {
  return jspb.Message.setField(this, 3, value || []);
};


/**
 * @param {!proto.model.ApplicationSyncStatus} value
 * @param {number=} opt_index
 * @return {!proto.grpc.service.webservice.ListApplicationsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.addSyncStatuses = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 3, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.ListApplicationsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.clearSyncStatusesList = function() {
  return this.setSyncStatusesList([]);
};


/**
 * optional string name = 5;
 * @return {string}
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.ListApplicationsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * map<string, string> labels = 6;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.getLabelsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 6, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.grpc.service.webservice.ListApplicationsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.clearLabelsMap = function() {
  this.getLabelsMap().clear();
  return this;
};


/**
 * optional string piped_id = 7;
 * @return {string}
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.getPipedId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.ListApplicationsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListApplicationsRequest.Options.prototype.setPipedId = function(value) {
  return jspb.Message.setProto3StringField(this, 7, value);
};


/**
 * optional Options options = 1;
 * @return {?proto.grpc.service.webservice.ListApplicationsRequest.Options}
 */
proto.grpc.service.webservice.ListApplicationsRequest.prototype.getOptions = function() {
  return /** @type{?proto.grpc.service.webservice.ListApplicationsRequest.Options} */ (
    jspb.Message.getWrapperField(this, proto.grpc.service.webservice.ListApplicationsRequest.Options, 1));
};


/**
 * @param {?proto.grpc.service.webservice.ListApplicationsRequest.Options|undefined} value
 * @return {!proto.grpc.service.webservice.ListApplicationsRequest} returns this
*/
proto.grpc.service.webservice.ListApplicationsRequest.prototype.setOptions = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.ListApplicationsRequest} returns this
 */
proto.grpc.service.webservice.ListApplicationsRequest.prototype.clearOptions = function() {
  return this.setOptions(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.ListApplicationsRequest.prototype.hasOptions = function() {
  return jspb.Message.getField(this, 1) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.grpc.service.webservice.ListApplicationsResponse.repeatedFields_ = [1];



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
proto.grpc.service.webservice.ListApplicationsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListApplicationsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListApplicationsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListApplicationsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    applicationsList: jspb.Message.toObjectList(msg.getApplicationsList(),
    pkg_model_application_pb.Application.toObject, includeInstance)
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
 * @return {!proto.grpc.service.webservice.ListApplicationsResponse}
 */
proto.grpc.service.webservice.ListApplicationsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListApplicationsResponse;
  return proto.grpc.service.webservice.ListApplicationsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListApplicationsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListApplicationsResponse}
 */
proto.grpc.service.webservice.ListApplicationsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new pkg_model_application_pb.Application;
      reader.readMessage(value,pkg_model_application_pb.Application.deserializeBinaryFromReader);
      msg.addApplications(value);
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
proto.grpc.service.webservice.ListApplicationsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListApplicationsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListApplicationsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListApplicationsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplicationsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      pkg_model_application_pb.Application.serializeBinaryToWriter
    );
  }
};


/**
 * repeated model.Application applications = 1;
 * @return {!Array<!proto.model.Application>}
 */
proto.grpc.service.webservice.ListApplicationsResponse.prototype.getApplicationsList = function() {
  return /** @type{!Array<!proto.model.Application>} */ (
    jspb.Message.getRepeatedWrapperField(this, pkg_model_application_pb.Application, 1));
};


/**
 * @param {!Array<!proto.model.Application>} value
 * @return {!proto.grpc.service.webservice.ListApplicationsResponse} returns this
*/
proto.grpc.service.webservice.ListApplicationsResponse.prototype.setApplicationsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.model.Application=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.Application}
 */
proto.grpc.service.webservice.ListApplicationsResponse.prototype.addApplications = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.model.Application, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.ListApplicationsResponse} returns this
 */
proto.grpc.service.webservice.ListApplicationsResponse.prototype.clearApplicationsList = function() {
  return this.setApplicationsList([]);
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
proto.grpc.service.webservice.SyncApplicationRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.SyncApplicationRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.SyncApplicationRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.SyncApplicationRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    applicationId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    syncStrategy: jspb.Message.getFieldWithDefault(msg, 2, 0)
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
 * @return {!proto.grpc.service.webservice.SyncApplicationRequest}
 */
proto.grpc.service.webservice.SyncApplicationRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.SyncApplicationRequest;
  return proto.grpc.service.webservice.SyncApplicationRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.SyncApplicationRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.SyncApplicationRequest}
 */
proto.grpc.service.webservice.SyncApplicationRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {!proto.model.SyncStrategy} */ (reader.readEnum());
      msg.setSyncStrategy(value);
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
proto.grpc.service.webservice.SyncApplicationRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.SyncApplicationRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.SyncApplicationRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.SyncApplicationRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplicationId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getSyncStrategy();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
};


/**
 * optional string application_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.SyncApplicationRequest.prototype.getApplicationId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.SyncApplicationRequest} returns this
 */
proto.grpc.service.webservice.SyncApplicationRequest.prototype.setApplicationId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional model.SyncStrategy sync_strategy = 2;
 * @return {!proto.model.SyncStrategy}
 */
proto.grpc.service.webservice.SyncApplicationRequest.prototype.getSyncStrategy = function() {
  return /** @type {!proto.model.SyncStrategy} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.model.SyncStrategy} value
 * @return {!proto.grpc.service.webservice.SyncApplicationRequest} returns this
 */
proto.grpc.service.webservice.SyncApplicationRequest.prototype.setSyncStrategy = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
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
proto.grpc.service.webservice.SyncApplicationResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.SyncApplicationResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.SyncApplicationResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.SyncApplicationResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    commandId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.SyncApplicationResponse}
 */
proto.grpc.service.webservice.SyncApplicationResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.SyncApplicationResponse;
  return proto.grpc.service.webservice.SyncApplicationResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.SyncApplicationResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.SyncApplicationResponse}
 */
proto.grpc.service.webservice.SyncApplicationResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setCommandId(value);
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
proto.grpc.service.webservice.SyncApplicationResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.SyncApplicationResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.SyncApplicationResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.SyncApplicationResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCommandId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string command_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.SyncApplicationResponse.prototype.getCommandId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.SyncApplicationResponse} returns this
 */
proto.grpc.service.webservice.SyncApplicationResponse.prototype.setCommandId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.GetApplicationRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetApplicationRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetApplicationRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetApplicationRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    applicationId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.GetApplicationRequest}
 */
proto.grpc.service.webservice.GetApplicationRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetApplicationRequest;
  return proto.grpc.service.webservice.GetApplicationRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetApplicationRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetApplicationRequest}
 */
proto.grpc.service.webservice.GetApplicationRequest.deserializeBinaryFromReader = function(msg, reader) {
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
proto.grpc.service.webservice.GetApplicationRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetApplicationRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetApplicationRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetApplicationRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplicationId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string application_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.GetApplicationRequest.prototype.getApplicationId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.GetApplicationRequest} returns this
 */
proto.grpc.service.webservice.GetApplicationRequest.prototype.setApplicationId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.GetApplicationResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetApplicationResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetApplicationResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetApplicationResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    application: (f = msg.getApplication()) && pkg_model_application_pb.Application.toObject(includeInstance, f)
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
 * @return {!proto.grpc.service.webservice.GetApplicationResponse}
 */
proto.grpc.service.webservice.GetApplicationResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetApplicationResponse;
  return proto.grpc.service.webservice.GetApplicationResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetApplicationResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetApplicationResponse}
 */
proto.grpc.service.webservice.GetApplicationResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new pkg_model_application_pb.Application;
      reader.readMessage(value,pkg_model_application_pb.Application.deserializeBinaryFromReader);
      msg.setApplication(value);
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
proto.grpc.service.webservice.GetApplicationResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetApplicationResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetApplicationResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetApplicationResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplication();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      pkg_model_application_pb.Application.serializeBinaryToWriter
    );
  }
};


/**
 * optional model.Application application = 1;
 * @return {?proto.model.Application}
 */
proto.grpc.service.webservice.GetApplicationResponse.prototype.getApplication = function() {
  return /** @type{?proto.model.Application} */ (
    jspb.Message.getWrapperField(this, pkg_model_application_pb.Application, 1));
};


/**
 * @param {?proto.model.Application|undefined} value
 * @return {!proto.grpc.service.webservice.GetApplicationResponse} returns this
*/
proto.grpc.service.webservice.GetApplicationResponse.prototype.setApplication = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.GetApplicationResponse} returns this
 */
proto.grpc.service.webservice.GetApplicationResponse.prototype.clearApplication = function() {
  return this.setApplication(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.GetApplicationResponse.prototype.hasApplication = function() {
  return jspb.Message.getField(this, 1) != null;
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
proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    pipedId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    data: jspb.Message.getFieldWithDefault(msg, 2, ""),
    base64Encoding: jspb.Message.getBooleanFieldWithDefault(msg, 3, false)
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
 * @return {!proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest}
 */
proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest;
  return proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest}
 */
proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setPipedId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setData(value);
      break;
    case 3:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setBase64Encoding(value);
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
proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPipedId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getData();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getBase64Encoding();
  if (f) {
    writer.writeBool(
      3,
      f
    );
  }
};


/**
 * optional string piped_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest.prototype.getPipedId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest} returns this
 */
proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest.prototype.setPipedId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string data = 2;
 * @return {string}
 */
proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest.prototype.getData = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest} returns this
 */
proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest.prototype.setData = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional bool base64_encoding = 3;
 * @return {boolean}
 */
proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest.prototype.getBase64Encoding = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 3, false));
};


/**
 * @param {boolean} value
 * @return {!proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest} returns this
 */
proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest.prototype.setBase64Encoding = function(value) {
  return jspb.Message.setProto3BooleanField(this, 3, value);
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
proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    data: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse}
 */
proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse;
  return proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse}
 */
proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setData(value);
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
proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getData();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string data = 1;
 * @return {string}
 */
proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse.prototype.getData = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse} returns this
 */
proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse.prototype.setData = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.ListUnregisteredApplicationsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListUnregisteredApplicationsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListUnregisteredApplicationsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListUnregisteredApplicationsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.ListUnregisteredApplicationsRequest}
 */
proto.grpc.service.webservice.ListUnregisteredApplicationsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListUnregisteredApplicationsRequest;
  return proto.grpc.service.webservice.ListUnregisteredApplicationsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListUnregisteredApplicationsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListUnregisteredApplicationsRequest}
 */
proto.grpc.service.webservice.ListUnregisteredApplicationsRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.ListUnregisteredApplicationsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListUnregisteredApplicationsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListUnregisteredApplicationsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListUnregisteredApplicationsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.grpc.service.webservice.ListUnregisteredApplicationsResponse.repeatedFields_ = [1];



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
proto.grpc.service.webservice.ListUnregisteredApplicationsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListUnregisteredApplicationsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListUnregisteredApplicationsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListUnregisteredApplicationsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    applicationsList: jspb.Message.toObjectList(msg.getApplicationsList(),
    pkg_model_common_pb.ApplicationInfo.toObject, includeInstance)
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
 * @return {!proto.grpc.service.webservice.ListUnregisteredApplicationsResponse}
 */
proto.grpc.service.webservice.ListUnregisteredApplicationsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListUnregisteredApplicationsResponse;
  return proto.grpc.service.webservice.ListUnregisteredApplicationsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListUnregisteredApplicationsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListUnregisteredApplicationsResponse}
 */
proto.grpc.service.webservice.ListUnregisteredApplicationsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new pkg_model_common_pb.ApplicationInfo;
      reader.readMessage(value,pkg_model_common_pb.ApplicationInfo.deserializeBinaryFromReader);
      msg.addApplications(value);
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
proto.grpc.service.webservice.ListUnregisteredApplicationsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListUnregisteredApplicationsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListUnregisteredApplicationsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListUnregisteredApplicationsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplicationsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      pkg_model_common_pb.ApplicationInfo.serializeBinaryToWriter
    );
  }
};


/**
 * repeated model.ApplicationInfo applications = 1;
 * @return {!Array<!proto.model.ApplicationInfo>}
 */
proto.grpc.service.webservice.ListUnregisteredApplicationsResponse.prototype.getApplicationsList = function() {
  return /** @type{!Array<!proto.model.ApplicationInfo>} */ (
    jspb.Message.getRepeatedWrapperField(this, pkg_model_common_pb.ApplicationInfo, 1));
};


/**
 * @param {!Array<!proto.model.ApplicationInfo>} value
 * @return {!proto.grpc.service.webservice.ListUnregisteredApplicationsResponse} returns this
*/
proto.grpc.service.webservice.ListUnregisteredApplicationsResponse.prototype.setApplicationsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.model.ApplicationInfo=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.ApplicationInfo}
 */
proto.grpc.service.webservice.ListUnregisteredApplicationsResponse.prototype.addApplications = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.model.ApplicationInfo, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.ListUnregisteredApplicationsResponse} returns this
 */
proto.grpc.service.webservice.ListUnregisteredApplicationsResponse.prototype.clearApplicationsList = function() {
  return this.setApplicationsList([]);
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
proto.grpc.service.webservice.ListDeploymentsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListDeploymentsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListDeploymentsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListDeploymentsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    options: (f = msg.getOptions()) && proto.grpc.service.webservice.ListDeploymentsRequest.Options.toObject(includeInstance, f),
    pageSize: jspb.Message.getFieldWithDefault(msg, 2, 0),
    cursor: jspb.Message.getFieldWithDefault(msg, 3, ""),
    pageMinUpdatedAt: jspb.Message.getFieldWithDefault(msg, 4, 0)
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
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest}
 */
proto.grpc.service.webservice.ListDeploymentsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListDeploymentsRequest;
  return proto.grpc.service.webservice.ListDeploymentsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListDeploymentsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest}
 */
proto.grpc.service.webservice.ListDeploymentsRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.grpc.service.webservice.ListDeploymentsRequest.Options;
      reader.readMessage(value,proto.grpc.service.webservice.ListDeploymentsRequest.Options.deserializeBinaryFromReader);
      msg.setOptions(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setPageSize(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setCursor(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setPageMinUpdatedAt(value);
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
proto.grpc.service.webservice.ListDeploymentsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListDeploymentsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListDeploymentsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListDeploymentsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOptions();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.grpc.service.webservice.ListDeploymentsRequest.Options.serializeBinaryToWriter
    );
  }
  f = message.getPageSize();
  if (f !== 0) {
    writer.writeInt32(
      2,
      f
    );
  }
  f = message.getCursor();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getPageMinUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      4,
      f
    );
  }
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.repeatedFields_ = [1,2,3];



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
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListDeploymentsRequest.Options.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListDeploymentsRequest.Options} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.toObject = function(includeInstance, msg) {
  var f, obj = {
    statusesList: (f = jspb.Message.getRepeatedField(msg, 1)) == null ? undefined : f,
    kindsList: (f = jspb.Message.getRepeatedField(msg, 2)) == null ? undefined : f,
    applicationIdsList: (f = jspb.Message.getRepeatedField(msg, 3)) == null ? undefined : f,
    applicationName: jspb.Message.getFieldWithDefault(msg, 5, ""),
    labelsMap: (f = msg.getLabelsMap()) ? f.toObject(includeInstance, undefined) : []
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
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest.Options}
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListDeploymentsRequest.Options;
  return proto.grpc.service.webservice.ListDeploymentsRequest.Options.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListDeploymentsRequest.Options} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest.Options}
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var values = /** @type {!Array<!proto.model.DeploymentStatus>} */ (reader.isDelimited() ? reader.readPackedEnum() : [reader.readEnum()]);
      for (var i = 0; i < values.length; i++) {
        msg.addStatuses(values[i]);
      }
      break;
    case 2:
      var values = /** @type {!Array<!proto.model.ApplicationKind>} */ (reader.isDelimited() ? reader.readPackedEnum() : [reader.readEnum()]);
      for (var i = 0; i < values.length; i++) {
        msg.addKinds(values[i]);
      }
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.addApplicationIds(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setApplicationName(value);
      break;
    case 6:
      var value = msg.getLabelsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "", "");
         });
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
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListDeploymentsRequest.Options.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListDeploymentsRequest.Options} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getStatusesList();
  if (f.length > 0) {
    writer.writePackedEnum(
      1,
      f
    );
  }
  f = message.getKindsList();
  if (f.length > 0) {
    writer.writePackedEnum(
      2,
      f
    );
  }
  f = message.getApplicationIdsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      3,
      f
    );
  }
  f = message.getApplicationName();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getLabelsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(6, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
};


/**
 * repeated model.DeploymentStatus statuses = 1;
 * @return {!Array<!proto.model.DeploymentStatus>}
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.getStatusesList = function() {
  return /** @type {!Array<!proto.model.DeploymentStatus>} */ (jspb.Message.getRepeatedField(this, 1));
};


/**
 * @param {!Array<!proto.model.DeploymentStatus>} value
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.setStatusesList = function(value) {
  return jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {!proto.model.DeploymentStatus} value
 * @param {number=} opt_index
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.addStatuses = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.clearStatusesList = function() {
  return this.setStatusesList([]);
};


/**
 * repeated model.ApplicationKind kinds = 2;
 * @return {!Array<!proto.model.ApplicationKind>}
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.getKindsList = function() {
  return /** @type {!Array<!proto.model.ApplicationKind>} */ (jspb.Message.getRepeatedField(this, 2));
};


/**
 * @param {!Array<!proto.model.ApplicationKind>} value
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.setKindsList = function(value) {
  return jspb.Message.setField(this, 2, value || []);
};


/**
 * @param {!proto.model.ApplicationKind} value
 * @param {number=} opt_index
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.addKinds = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 2, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.clearKindsList = function() {
  return this.setKindsList([]);
};


/**
 * repeated string application_ids = 3;
 * @return {!Array<string>}
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.getApplicationIdsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 3));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.setApplicationIdsList = function(value) {
  return jspb.Message.setField(this, 3, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.addApplicationIds = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 3, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.clearApplicationIdsList = function() {
  return this.setApplicationIdsList([]);
};


/**
 * optional string application_name = 5;
 * @return {string}
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.getApplicationName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.setApplicationName = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * map<string, string> labels = 6;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.getLabelsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 6, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListDeploymentsRequest.Options.prototype.clearLabelsMap = function() {
  this.getLabelsMap().clear();
  return this;
};


/**
 * optional Options options = 1;
 * @return {?proto.grpc.service.webservice.ListDeploymentsRequest.Options}
 */
proto.grpc.service.webservice.ListDeploymentsRequest.prototype.getOptions = function() {
  return /** @type{?proto.grpc.service.webservice.ListDeploymentsRequest.Options} */ (
    jspb.Message.getWrapperField(this, proto.grpc.service.webservice.ListDeploymentsRequest.Options, 1));
};


/**
 * @param {?proto.grpc.service.webservice.ListDeploymentsRequest.Options|undefined} value
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest} returns this
*/
proto.grpc.service.webservice.ListDeploymentsRequest.prototype.setOptions = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest} returns this
 */
proto.grpc.service.webservice.ListDeploymentsRequest.prototype.clearOptions = function() {
  return this.setOptions(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.ListDeploymentsRequest.prototype.hasOptions = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional int32 page_size = 2;
 * @return {number}
 */
proto.grpc.service.webservice.ListDeploymentsRequest.prototype.getPageSize = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest} returns this
 */
proto.grpc.service.webservice.ListDeploymentsRequest.prototype.setPageSize = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional string cursor = 3;
 * @return {string}
 */
proto.grpc.service.webservice.ListDeploymentsRequest.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest} returns this
 */
proto.grpc.service.webservice.ListDeploymentsRequest.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional int64 page_min_updated_at = 4;
 * @return {number}
 */
proto.grpc.service.webservice.ListDeploymentsRequest.prototype.getPageMinUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.grpc.service.webservice.ListDeploymentsRequest} returns this
 */
proto.grpc.service.webservice.ListDeploymentsRequest.prototype.setPageMinUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.grpc.service.webservice.ListDeploymentsResponse.repeatedFields_ = [1];



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
proto.grpc.service.webservice.ListDeploymentsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListDeploymentsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListDeploymentsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListDeploymentsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    deploymentsList: jspb.Message.toObjectList(msg.getDeploymentsList(),
    pkg_model_deployment_pb.Deployment.toObject, includeInstance),
    cursor: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.grpc.service.webservice.ListDeploymentsResponse}
 */
proto.grpc.service.webservice.ListDeploymentsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListDeploymentsResponse;
  return proto.grpc.service.webservice.ListDeploymentsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListDeploymentsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListDeploymentsResponse}
 */
proto.grpc.service.webservice.ListDeploymentsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new pkg_model_deployment_pb.Deployment;
      reader.readMessage(value,pkg_model_deployment_pb.Deployment.deserializeBinaryFromReader);
      msg.addDeployments(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setCursor(value);
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
proto.grpc.service.webservice.ListDeploymentsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListDeploymentsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListDeploymentsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListDeploymentsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeploymentsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      pkg_model_deployment_pb.Deployment.serializeBinaryToWriter
    );
  }
  f = message.getCursor();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * repeated model.Deployment deployments = 1;
 * @return {!Array<!proto.model.Deployment>}
 */
proto.grpc.service.webservice.ListDeploymentsResponse.prototype.getDeploymentsList = function() {
  return /** @type{!Array<!proto.model.Deployment>} */ (
    jspb.Message.getRepeatedWrapperField(this, pkg_model_deployment_pb.Deployment, 1));
};


/**
 * @param {!Array<!proto.model.Deployment>} value
 * @return {!proto.grpc.service.webservice.ListDeploymentsResponse} returns this
*/
proto.grpc.service.webservice.ListDeploymentsResponse.prototype.setDeploymentsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.model.Deployment=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.Deployment}
 */
proto.grpc.service.webservice.ListDeploymentsResponse.prototype.addDeployments = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.model.Deployment, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.ListDeploymentsResponse} returns this
 */
proto.grpc.service.webservice.ListDeploymentsResponse.prototype.clearDeploymentsList = function() {
  return this.setDeploymentsList([]);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.grpc.service.webservice.ListDeploymentsResponse.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.ListDeploymentsResponse} returns this
 */
proto.grpc.service.webservice.ListDeploymentsResponse.prototype.setCursor = function(value) {
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
proto.grpc.service.webservice.GetDeploymentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetDeploymentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetDeploymentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetDeploymentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    deploymentId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.GetDeploymentRequest}
 */
proto.grpc.service.webservice.GetDeploymentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetDeploymentRequest;
  return proto.grpc.service.webservice.GetDeploymentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetDeploymentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetDeploymentRequest}
 */
proto.grpc.service.webservice.GetDeploymentRequest.deserializeBinaryFromReader = function(msg, reader) {
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
proto.grpc.service.webservice.GetDeploymentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetDeploymentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetDeploymentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetDeploymentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeploymentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string deployment_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.GetDeploymentRequest.prototype.getDeploymentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.GetDeploymentRequest} returns this
 */
proto.grpc.service.webservice.GetDeploymentRequest.prototype.setDeploymentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.GetDeploymentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetDeploymentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetDeploymentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetDeploymentResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    deployment: (f = msg.getDeployment()) && pkg_model_deployment_pb.Deployment.toObject(includeInstance, f)
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
 * @return {!proto.grpc.service.webservice.GetDeploymentResponse}
 */
proto.grpc.service.webservice.GetDeploymentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetDeploymentResponse;
  return proto.grpc.service.webservice.GetDeploymentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetDeploymentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetDeploymentResponse}
 */
proto.grpc.service.webservice.GetDeploymentResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new pkg_model_deployment_pb.Deployment;
      reader.readMessage(value,pkg_model_deployment_pb.Deployment.deserializeBinaryFromReader);
      msg.setDeployment(value);
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
proto.grpc.service.webservice.GetDeploymentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetDeploymentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetDeploymentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetDeploymentResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeployment();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      pkg_model_deployment_pb.Deployment.serializeBinaryToWriter
    );
  }
};


/**
 * optional model.Deployment deployment = 1;
 * @return {?proto.model.Deployment}
 */
proto.grpc.service.webservice.GetDeploymentResponse.prototype.getDeployment = function() {
  return /** @type{?proto.model.Deployment} */ (
    jspb.Message.getWrapperField(this, pkg_model_deployment_pb.Deployment, 1));
};


/**
 * @param {?proto.model.Deployment|undefined} value
 * @return {!proto.grpc.service.webservice.GetDeploymentResponse} returns this
*/
proto.grpc.service.webservice.GetDeploymentResponse.prototype.setDeployment = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.GetDeploymentResponse} returns this
 */
proto.grpc.service.webservice.GetDeploymentResponse.prototype.clearDeployment = function() {
  return this.setDeployment(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.GetDeploymentResponse.prototype.hasDeployment = function() {
  return jspb.Message.getField(this, 1) != null;
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
proto.grpc.service.webservice.GetStageLogRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetStageLogRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetStageLogRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetStageLogRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    deploymentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    stageId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    retriedCount: jspb.Message.getFieldWithDefault(msg, 3, 0),
    offsetIndex: jspb.Message.getFieldWithDefault(msg, 4, 0)
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
 * @return {!proto.grpc.service.webservice.GetStageLogRequest}
 */
proto.grpc.service.webservice.GetStageLogRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetStageLogRequest;
  return proto.grpc.service.webservice.GetStageLogRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetStageLogRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetStageLogRequest}
 */
proto.grpc.service.webservice.GetStageLogRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {string} */ (reader.readString());
      msg.setStageId(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setRetriedCount(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setOffsetIndex(value);
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
proto.grpc.service.webservice.GetStageLogRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetStageLogRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetStageLogRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetStageLogRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeploymentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStageId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getRetriedCount();
  if (f !== 0) {
    writer.writeInt32(
      3,
      f
    );
  }
  f = message.getOffsetIndex();
  if (f !== 0) {
    writer.writeInt64(
      4,
      f
    );
  }
};


/**
 * optional string deployment_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.GetStageLogRequest.prototype.getDeploymentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.GetStageLogRequest} returns this
 */
proto.grpc.service.webservice.GetStageLogRequest.prototype.setDeploymentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string stage_id = 2;
 * @return {string}
 */
proto.grpc.service.webservice.GetStageLogRequest.prototype.getStageId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.GetStageLogRequest} returns this
 */
proto.grpc.service.webservice.GetStageLogRequest.prototype.setStageId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional int32 retried_count = 3;
 * @return {number}
 */
proto.grpc.service.webservice.GetStageLogRequest.prototype.getRetriedCount = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.grpc.service.webservice.GetStageLogRequest} returns this
 */
proto.grpc.service.webservice.GetStageLogRequest.prototype.setRetriedCount = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * optional int64 offset_index = 4;
 * @return {number}
 */
proto.grpc.service.webservice.GetStageLogRequest.prototype.getOffsetIndex = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.grpc.service.webservice.GetStageLogRequest} returns this
 */
proto.grpc.service.webservice.GetStageLogRequest.prototype.setOffsetIndex = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.grpc.service.webservice.GetStageLogResponse.repeatedFields_ = [1];



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
proto.grpc.service.webservice.GetStageLogResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetStageLogResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetStageLogResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetStageLogResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    blocksList: jspb.Message.toObjectList(msg.getBlocksList(),
    pkg_model_logblock_pb.LogBlock.toObject, includeInstance),
    completed: jspb.Message.getBooleanFieldWithDefault(msg, 2, false)
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
 * @return {!proto.grpc.service.webservice.GetStageLogResponse}
 */
proto.grpc.service.webservice.GetStageLogResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetStageLogResponse;
  return proto.grpc.service.webservice.GetStageLogResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetStageLogResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetStageLogResponse}
 */
proto.grpc.service.webservice.GetStageLogResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new pkg_model_logblock_pb.LogBlock;
      reader.readMessage(value,pkg_model_logblock_pb.LogBlock.deserializeBinaryFromReader);
      msg.addBlocks(value);
      break;
    case 2:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setCompleted(value);
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
proto.grpc.service.webservice.GetStageLogResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetStageLogResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetStageLogResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetStageLogResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getBlocksList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      pkg_model_logblock_pb.LogBlock.serializeBinaryToWriter
    );
  }
  f = message.getCompleted();
  if (f) {
    writer.writeBool(
      2,
      f
    );
  }
};


/**
 * repeated model.LogBlock blocks = 1;
 * @return {!Array<!proto.model.LogBlock>}
 */
proto.grpc.service.webservice.GetStageLogResponse.prototype.getBlocksList = function() {
  return /** @type{!Array<!proto.model.LogBlock>} */ (
    jspb.Message.getRepeatedWrapperField(this, pkg_model_logblock_pb.LogBlock, 1));
};


/**
 * @param {!Array<!proto.model.LogBlock>} value
 * @return {!proto.grpc.service.webservice.GetStageLogResponse} returns this
*/
proto.grpc.service.webservice.GetStageLogResponse.prototype.setBlocksList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.model.LogBlock=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.LogBlock}
 */
proto.grpc.service.webservice.GetStageLogResponse.prototype.addBlocks = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.model.LogBlock, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.GetStageLogResponse} returns this
 */
proto.grpc.service.webservice.GetStageLogResponse.prototype.clearBlocksList = function() {
  return this.setBlocksList([]);
};


/**
 * optional bool completed = 2;
 * @return {boolean}
 */
proto.grpc.service.webservice.GetStageLogResponse.prototype.getCompleted = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 2, false));
};


/**
 * @param {boolean} value
 * @return {!proto.grpc.service.webservice.GetStageLogResponse} returns this
 */
proto.grpc.service.webservice.GetStageLogResponse.prototype.setCompleted = function(value) {
  return jspb.Message.setProto3BooleanField(this, 2, value);
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
proto.grpc.service.webservice.CancelDeploymentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.CancelDeploymentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.CancelDeploymentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.CancelDeploymentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    deploymentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    forceRollback: jspb.Message.getBooleanFieldWithDefault(msg, 2, false),
    forceNoRollback: jspb.Message.getBooleanFieldWithDefault(msg, 3, false)
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
 * @return {!proto.grpc.service.webservice.CancelDeploymentRequest}
 */
proto.grpc.service.webservice.CancelDeploymentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.CancelDeploymentRequest;
  return proto.grpc.service.webservice.CancelDeploymentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.CancelDeploymentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.CancelDeploymentRequest}
 */
proto.grpc.service.webservice.CancelDeploymentRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setForceRollback(value);
      break;
    case 3:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setForceNoRollback(value);
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
proto.grpc.service.webservice.CancelDeploymentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.CancelDeploymentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.CancelDeploymentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.CancelDeploymentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeploymentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getForceRollback();
  if (f) {
    writer.writeBool(
      2,
      f
    );
  }
  f = message.getForceNoRollback();
  if (f) {
    writer.writeBool(
      3,
      f
    );
  }
};


/**
 * optional string deployment_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.CancelDeploymentRequest.prototype.getDeploymentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.CancelDeploymentRequest} returns this
 */
proto.grpc.service.webservice.CancelDeploymentRequest.prototype.setDeploymentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bool force_rollback = 2;
 * @return {boolean}
 */
proto.grpc.service.webservice.CancelDeploymentRequest.prototype.getForceRollback = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 2, false));
};


/**
 * @param {boolean} value
 * @return {!proto.grpc.service.webservice.CancelDeploymentRequest} returns this
 */
proto.grpc.service.webservice.CancelDeploymentRequest.prototype.setForceRollback = function(value) {
  return jspb.Message.setProto3BooleanField(this, 2, value);
};


/**
 * optional bool force_no_rollback = 3;
 * @return {boolean}
 */
proto.grpc.service.webservice.CancelDeploymentRequest.prototype.getForceNoRollback = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 3, false));
};


/**
 * @param {boolean} value
 * @return {!proto.grpc.service.webservice.CancelDeploymentRequest} returns this
 */
proto.grpc.service.webservice.CancelDeploymentRequest.prototype.setForceNoRollback = function(value) {
  return jspb.Message.setProto3BooleanField(this, 3, value);
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
proto.grpc.service.webservice.CancelDeploymentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.CancelDeploymentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.CancelDeploymentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.CancelDeploymentResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    commandId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.CancelDeploymentResponse}
 */
proto.grpc.service.webservice.CancelDeploymentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.CancelDeploymentResponse;
  return proto.grpc.service.webservice.CancelDeploymentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.CancelDeploymentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.CancelDeploymentResponse}
 */
proto.grpc.service.webservice.CancelDeploymentResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setCommandId(value);
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
proto.grpc.service.webservice.CancelDeploymentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.CancelDeploymentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.CancelDeploymentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.CancelDeploymentResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCommandId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string command_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.CancelDeploymentResponse.prototype.getCommandId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.CancelDeploymentResponse} returns this
 */
proto.grpc.service.webservice.CancelDeploymentResponse.prototype.setCommandId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.SkipStageRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.SkipStageRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.SkipStageRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.SkipStageRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    deploymentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    stageId: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.grpc.service.webservice.SkipStageRequest}
 */
proto.grpc.service.webservice.SkipStageRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.SkipStageRequest;
  return proto.grpc.service.webservice.SkipStageRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.SkipStageRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.SkipStageRequest}
 */
proto.grpc.service.webservice.SkipStageRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {string} */ (reader.readString());
      msg.setStageId(value);
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
proto.grpc.service.webservice.SkipStageRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.SkipStageRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.SkipStageRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.SkipStageRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeploymentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStageId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string deployment_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.SkipStageRequest.prototype.getDeploymentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.SkipStageRequest} returns this
 */
proto.grpc.service.webservice.SkipStageRequest.prototype.setDeploymentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string stage_id = 2;
 * @return {string}
 */
proto.grpc.service.webservice.SkipStageRequest.prototype.getStageId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.SkipStageRequest} returns this
 */
proto.grpc.service.webservice.SkipStageRequest.prototype.setStageId = function(value) {
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
proto.grpc.service.webservice.SkipStageResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.SkipStageResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.SkipStageResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.SkipStageResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    commandId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.SkipStageResponse}
 */
proto.grpc.service.webservice.SkipStageResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.SkipStageResponse;
  return proto.grpc.service.webservice.SkipStageResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.SkipStageResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.SkipStageResponse}
 */
proto.grpc.service.webservice.SkipStageResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setCommandId(value);
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
proto.grpc.service.webservice.SkipStageResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.SkipStageResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.SkipStageResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.SkipStageResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCommandId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string command_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.SkipStageResponse.prototype.getCommandId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.SkipStageResponse} returns this
 */
proto.grpc.service.webservice.SkipStageResponse.prototype.setCommandId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.ApproveStageRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ApproveStageRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ApproveStageRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ApproveStageRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    deploymentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    stageId: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.grpc.service.webservice.ApproveStageRequest}
 */
proto.grpc.service.webservice.ApproveStageRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ApproveStageRequest;
  return proto.grpc.service.webservice.ApproveStageRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ApproveStageRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ApproveStageRequest}
 */
proto.grpc.service.webservice.ApproveStageRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {string} */ (reader.readString());
      msg.setStageId(value);
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
proto.grpc.service.webservice.ApproveStageRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ApproveStageRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ApproveStageRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ApproveStageRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeploymentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStageId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string deployment_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.ApproveStageRequest.prototype.getDeploymentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.ApproveStageRequest} returns this
 */
proto.grpc.service.webservice.ApproveStageRequest.prototype.setDeploymentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string stage_id = 2;
 * @return {string}
 */
proto.grpc.service.webservice.ApproveStageRequest.prototype.getStageId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.ApproveStageRequest} returns this
 */
proto.grpc.service.webservice.ApproveStageRequest.prototype.setStageId = function(value) {
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
proto.grpc.service.webservice.ApproveStageResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ApproveStageResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ApproveStageResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ApproveStageResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    commandId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.ApproveStageResponse}
 */
proto.grpc.service.webservice.ApproveStageResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ApproveStageResponse;
  return proto.grpc.service.webservice.ApproveStageResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ApproveStageResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ApproveStageResponse}
 */
proto.grpc.service.webservice.ApproveStageResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setCommandId(value);
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
proto.grpc.service.webservice.ApproveStageResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ApproveStageResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ApproveStageResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ApproveStageResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCommandId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string command_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.ApproveStageResponse.prototype.getCommandId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.ApproveStageResponse} returns this
 */
proto.grpc.service.webservice.ApproveStageResponse.prototype.setCommandId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.GetApplicationLiveStateRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetApplicationLiveStateRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetApplicationLiveStateRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetApplicationLiveStateRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    applicationId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.GetApplicationLiveStateRequest}
 */
proto.grpc.service.webservice.GetApplicationLiveStateRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetApplicationLiveStateRequest;
  return proto.grpc.service.webservice.GetApplicationLiveStateRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetApplicationLiveStateRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetApplicationLiveStateRequest}
 */
proto.grpc.service.webservice.GetApplicationLiveStateRequest.deserializeBinaryFromReader = function(msg, reader) {
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
proto.grpc.service.webservice.GetApplicationLiveStateRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetApplicationLiveStateRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetApplicationLiveStateRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetApplicationLiveStateRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplicationId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string application_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.GetApplicationLiveStateRequest.prototype.getApplicationId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.GetApplicationLiveStateRequest} returns this
 */
proto.grpc.service.webservice.GetApplicationLiveStateRequest.prototype.setApplicationId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.GetApplicationLiveStateResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetApplicationLiveStateResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetApplicationLiveStateResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetApplicationLiveStateResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    snapshot: (f = msg.getSnapshot()) && pkg_model_application_live_state_pb.ApplicationLiveStateSnapshot.toObject(includeInstance, f)
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
 * @return {!proto.grpc.service.webservice.GetApplicationLiveStateResponse}
 */
proto.grpc.service.webservice.GetApplicationLiveStateResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetApplicationLiveStateResponse;
  return proto.grpc.service.webservice.GetApplicationLiveStateResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetApplicationLiveStateResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetApplicationLiveStateResponse}
 */
proto.grpc.service.webservice.GetApplicationLiveStateResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new pkg_model_application_live_state_pb.ApplicationLiveStateSnapshot;
      reader.readMessage(value,pkg_model_application_live_state_pb.ApplicationLiveStateSnapshot.deserializeBinaryFromReader);
      msg.setSnapshot(value);
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
proto.grpc.service.webservice.GetApplicationLiveStateResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetApplicationLiveStateResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetApplicationLiveStateResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetApplicationLiveStateResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSnapshot();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      pkg_model_application_live_state_pb.ApplicationLiveStateSnapshot.serializeBinaryToWriter
    );
  }
};


/**
 * optional model.ApplicationLiveStateSnapshot snapshot = 1;
 * @return {?proto.model.ApplicationLiveStateSnapshot}
 */
proto.grpc.service.webservice.GetApplicationLiveStateResponse.prototype.getSnapshot = function() {
  return /** @type{?proto.model.ApplicationLiveStateSnapshot} */ (
    jspb.Message.getWrapperField(this, pkg_model_application_live_state_pb.ApplicationLiveStateSnapshot, 1));
};


/**
 * @param {?proto.model.ApplicationLiveStateSnapshot|undefined} value
 * @return {!proto.grpc.service.webservice.GetApplicationLiveStateResponse} returns this
*/
proto.grpc.service.webservice.GetApplicationLiveStateResponse.prototype.setSnapshot = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.GetApplicationLiveStateResponse} returns this
 */
proto.grpc.service.webservice.GetApplicationLiveStateResponse.prototype.clearSnapshot = function() {
  return this.setSnapshot(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.GetApplicationLiveStateResponse.prototype.hasSnapshot = function() {
  return jspb.Message.getField(this, 1) != null;
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
proto.grpc.service.webservice.GetProjectRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetProjectRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetProjectRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetProjectRequest.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.GetProjectRequest}
 */
proto.grpc.service.webservice.GetProjectRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetProjectRequest;
  return proto.grpc.service.webservice.GetProjectRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetProjectRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetProjectRequest}
 */
proto.grpc.service.webservice.GetProjectRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.GetProjectRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetProjectRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetProjectRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetProjectRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.GetProjectResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetProjectResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetProjectResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetProjectResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    project: (f = msg.getProject()) && pkg_model_project_pb.Project.toObject(includeInstance, f)
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
 * @return {!proto.grpc.service.webservice.GetProjectResponse}
 */
proto.grpc.service.webservice.GetProjectResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetProjectResponse;
  return proto.grpc.service.webservice.GetProjectResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetProjectResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetProjectResponse}
 */
proto.grpc.service.webservice.GetProjectResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new pkg_model_project_pb.Project;
      reader.readMessage(value,pkg_model_project_pb.Project.deserializeBinaryFromReader);
      msg.setProject(value);
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
proto.grpc.service.webservice.GetProjectResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetProjectResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetProjectResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetProjectResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getProject();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      pkg_model_project_pb.Project.serializeBinaryToWriter
    );
  }
};


/**
 * optional model.Project project = 1;
 * @return {?proto.model.Project}
 */
proto.grpc.service.webservice.GetProjectResponse.prototype.getProject = function() {
  return /** @type{?proto.model.Project} */ (
    jspb.Message.getWrapperField(this, pkg_model_project_pb.Project, 1));
};


/**
 * @param {?proto.model.Project|undefined} value
 * @return {!proto.grpc.service.webservice.GetProjectResponse} returns this
*/
proto.grpc.service.webservice.GetProjectResponse.prototype.setProject = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.GetProjectResponse} returns this
 */
proto.grpc.service.webservice.GetProjectResponse.prototype.clearProject = function() {
  return this.setProject(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.GetProjectResponse.prototype.hasProject = function() {
  return jspb.Message.getField(this, 1) != null;
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
proto.grpc.service.webservice.UpdateProjectStaticAdminRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.UpdateProjectStaticAdminRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.UpdateProjectStaticAdminRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateProjectStaticAdminRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    username: jspb.Message.getFieldWithDefault(msg, 1, ""),
    password: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.grpc.service.webservice.UpdateProjectStaticAdminRequest}
 */
proto.grpc.service.webservice.UpdateProjectStaticAdminRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.UpdateProjectStaticAdminRequest;
  return proto.grpc.service.webservice.UpdateProjectStaticAdminRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.UpdateProjectStaticAdminRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.UpdateProjectStaticAdminRequest}
 */
proto.grpc.service.webservice.UpdateProjectStaticAdminRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setUsername(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setPassword(value);
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
proto.grpc.service.webservice.UpdateProjectStaticAdminRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.UpdateProjectStaticAdminRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.UpdateProjectStaticAdminRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateProjectStaticAdminRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUsername();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPassword();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string username = 1;
 * @return {string}
 */
proto.grpc.service.webservice.UpdateProjectStaticAdminRequest.prototype.getUsername = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.UpdateProjectStaticAdminRequest} returns this
 */
proto.grpc.service.webservice.UpdateProjectStaticAdminRequest.prototype.setUsername = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string password = 2;
 * @return {string}
 */
proto.grpc.service.webservice.UpdateProjectStaticAdminRequest.prototype.getPassword = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.UpdateProjectStaticAdminRequest} returns this
 */
proto.grpc.service.webservice.UpdateProjectStaticAdminRequest.prototype.setPassword = function(value) {
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
proto.grpc.service.webservice.UpdateProjectStaticAdminResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.UpdateProjectStaticAdminResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.UpdateProjectStaticAdminResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateProjectStaticAdminResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.UpdateProjectStaticAdminResponse}
 */
proto.grpc.service.webservice.UpdateProjectStaticAdminResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.UpdateProjectStaticAdminResponse;
  return proto.grpc.service.webservice.UpdateProjectStaticAdminResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.UpdateProjectStaticAdminResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.UpdateProjectStaticAdminResponse}
 */
proto.grpc.service.webservice.UpdateProjectStaticAdminResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.UpdateProjectStaticAdminResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.UpdateProjectStaticAdminResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.UpdateProjectStaticAdminResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateProjectStaticAdminResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.UpdateProjectSSOConfigRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.UpdateProjectSSOConfigRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.UpdateProjectSSOConfigRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateProjectSSOConfigRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    sso: (f = msg.getSso()) && pkg_model_project_pb.ProjectSSOConfig.toObject(includeInstance, f)
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
 * @return {!proto.grpc.service.webservice.UpdateProjectSSOConfigRequest}
 */
proto.grpc.service.webservice.UpdateProjectSSOConfigRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.UpdateProjectSSOConfigRequest;
  return proto.grpc.service.webservice.UpdateProjectSSOConfigRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.UpdateProjectSSOConfigRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.UpdateProjectSSOConfigRequest}
 */
proto.grpc.service.webservice.UpdateProjectSSOConfigRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new pkg_model_project_pb.ProjectSSOConfig;
      reader.readMessage(value,pkg_model_project_pb.ProjectSSOConfig.deserializeBinaryFromReader);
      msg.setSso(value);
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
proto.grpc.service.webservice.UpdateProjectSSOConfigRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.UpdateProjectSSOConfigRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.UpdateProjectSSOConfigRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateProjectSSOConfigRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSso();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      pkg_model_project_pb.ProjectSSOConfig.serializeBinaryToWriter
    );
  }
};


/**
 * optional model.ProjectSSOConfig sso = 1;
 * @return {?proto.model.ProjectSSOConfig}
 */
proto.grpc.service.webservice.UpdateProjectSSOConfigRequest.prototype.getSso = function() {
  return /** @type{?proto.model.ProjectSSOConfig} */ (
    jspb.Message.getWrapperField(this, pkg_model_project_pb.ProjectSSOConfig, 1));
};


/**
 * @param {?proto.model.ProjectSSOConfig|undefined} value
 * @return {!proto.grpc.service.webservice.UpdateProjectSSOConfigRequest} returns this
*/
proto.grpc.service.webservice.UpdateProjectSSOConfigRequest.prototype.setSso = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.UpdateProjectSSOConfigRequest} returns this
 */
proto.grpc.service.webservice.UpdateProjectSSOConfigRequest.prototype.clearSso = function() {
  return this.setSso(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.UpdateProjectSSOConfigRequest.prototype.hasSso = function() {
  return jspb.Message.getField(this, 1) != null;
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
proto.grpc.service.webservice.UpdateProjectSSOConfigResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.UpdateProjectSSOConfigResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.UpdateProjectSSOConfigResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateProjectSSOConfigResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.UpdateProjectSSOConfigResponse}
 */
proto.grpc.service.webservice.UpdateProjectSSOConfigResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.UpdateProjectSSOConfigResponse;
  return proto.grpc.service.webservice.UpdateProjectSSOConfigResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.UpdateProjectSSOConfigResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.UpdateProjectSSOConfigResponse}
 */
proto.grpc.service.webservice.UpdateProjectSSOConfigResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.UpdateProjectSSOConfigResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.UpdateProjectSSOConfigResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.UpdateProjectSSOConfigResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateProjectSSOConfigResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.UpdateProjectRBACConfigRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.UpdateProjectRBACConfigRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.UpdateProjectRBACConfigRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateProjectRBACConfigRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    rbac: (f = msg.getRbac()) && pkg_model_project_pb.ProjectRBACConfig.toObject(includeInstance, f)
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
 * @return {!proto.grpc.service.webservice.UpdateProjectRBACConfigRequest}
 */
proto.grpc.service.webservice.UpdateProjectRBACConfigRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.UpdateProjectRBACConfigRequest;
  return proto.grpc.service.webservice.UpdateProjectRBACConfigRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.UpdateProjectRBACConfigRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.UpdateProjectRBACConfigRequest}
 */
proto.grpc.service.webservice.UpdateProjectRBACConfigRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new pkg_model_project_pb.ProjectRBACConfig;
      reader.readMessage(value,pkg_model_project_pb.ProjectRBACConfig.deserializeBinaryFromReader);
      msg.setRbac(value);
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
proto.grpc.service.webservice.UpdateProjectRBACConfigRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.UpdateProjectRBACConfigRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.UpdateProjectRBACConfigRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateProjectRBACConfigRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRbac();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      pkg_model_project_pb.ProjectRBACConfig.serializeBinaryToWriter
    );
  }
};


/**
 * optional model.ProjectRBACConfig rbac = 1;
 * @return {?proto.model.ProjectRBACConfig}
 */
proto.grpc.service.webservice.UpdateProjectRBACConfigRequest.prototype.getRbac = function() {
  return /** @type{?proto.model.ProjectRBACConfig} */ (
    jspb.Message.getWrapperField(this, pkg_model_project_pb.ProjectRBACConfig, 1));
};


/**
 * @param {?proto.model.ProjectRBACConfig|undefined} value
 * @return {!proto.grpc.service.webservice.UpdateProjectRBACConfigRequest} returns this
*/
proto.grpc.service.webservice.UpdateProjectRBACConfigRequest.prototype.setRbac = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.UpdateProjectRBACConfigRequest} returns this
 */
proto.grpc.service.webservice.UpdateProjectRBACConfigRequest.prototype.clearRbac = function() {
  return this.setRbac(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.UpdateProjectRBACConfigRequest.prototype.hasRbac = function() {
  return jspb.Message.getField(this, 1) != null;
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
proto.grpc.service.webservice.UpdateProjectRBACConfigResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.UpdateProjectRBACConfigResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.UpdateProjectRBACConfigResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateProjectRBACConfigResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.UpdateProjectRBACConfigResponse}
 */
proto.grpc.service.webservice.UpdateProjectRBACConfigResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.UpdateProjectRBACConfigResponse;
  return proto.grpc.service.webservice.UpdateProjectRBACConfigResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.UpdateProjectRBACConfigResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.UpdateProjectRBACConfigResponse}
 */
proto.grpc.service.webservice.UpdateProjectRBACConfigResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.UpdateProjectRBACConfigResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.UpdateProjectRBACConfigResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.UpdateProjectRBACConfigResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateProjectRBACConfigResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.EnableStaticAdminRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.EnableStaticAdminRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.EnableStaticAdminRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.EnableStaticAdminRequest.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.EnableStaticAdminRequest}
 */
proto.grpc.service.webservice.EnableStaticAdminRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.EnableStaticAdminRequest;
  return proto.grpc.service.webservice.EnableStaticAdminRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.EnableStaticAdminRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.EnableStaticAdminRequest}
 */
proto.grpc.service.webservice.EnableStaticAdminRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.EnableStaticAdminRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.EnableStaticAdminRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.EnableStaticAdminRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.EnableStaticAdminRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.EnableStaticAdminResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.EnableStaticAdminResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.EnableStaticAdminResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.EnableStaticAdminResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.EnableStaticAdminResponse}
 */
proto.grpc.service.webservice.EnableStaticAdminResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.EnableStaticAdminResponse;
  return proto.grpc.service.webservice.EnableStaticAdminResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.EnableStaticAdminResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.EnableStaticAdminResponse}
 */
proto.grpc.service.webservice.EnableStaticAdminResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.EnableStaticAdminResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.EnableStaticAdminResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.EnableStaticAdminResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.EnableStaticAdminResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.DisableStaticAdminRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.DisableStaticAdminRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.DisableStaticAdminRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DisableStaticAdminRequest.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.DisableStaticAdminRequest}
 */
proto.grpc.service.webservice.DisableStaticAdminRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.DisableStaticAdminRequest;
  return proto.grpc.service.webservice.DisableStaticAdminRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.DisableStaticAdminRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.DisableStaticAdminRequest}
 */
proto.grpc.service.webservice.DisableStaticAdminRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.DisableStaticAdminRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.DisableStaticAdminRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.DisableStaticAdminRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DisableStaticAdminRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.DisableStaticAdminResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.DisableStaticAdminResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.DisableStaticAdminResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DisableStaticAdminResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.DisableStaticAdminResponse}
 */
proto.grpc.service.webservice.DisableStaticAdminResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.DisableStaticAdminResponse;
  return proto.grpc.service.webservice.DisableStaticAdminResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.DisableStaticAdminResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.DisableStaticAdminResponse}
 */
proto.grpc.service.webservice.DisableStaticAdminResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.DisableStaticAdminResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.DisableStaticAdminResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.DisableStaticAdminResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DisableStaticAdminResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.GetMeRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetMeRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetMeRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetMeRequest.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.GetMeRequest}
 */
proto.grpc.service.webservice.GetMeRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetMeRequest;
  return proto.grpc.service.webservice.GetMeRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetMeRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetMeRequest}
 */
proto.grpc.service.webservice.GetMeRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.GetMeRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetMeRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetMeRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetMeRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.GetMeResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetMeResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetMeResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetMeResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    subject: jspb.Message.getFieldWithDefault(msg, 1, ""),
    avatarUrl: jspb.Message.getFieldWithDefault(msg, 2, ""),
    projectId: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.grpc.service.webservice.GetMeResponse}
 */
proto.grpc.service.webservice.GetMeResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetMeResponse;
  return proto.grpc.service.webservice.GetMeResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetMeResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetMeResponse}
 */
proto.grpc.service.webservice.GetMeResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setSubject(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setAvatarUrl(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setProjectId(value);
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
proto.grpc.service.webservice.GetMeResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetMeResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetMeResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetMeResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSubject();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getAvatarUrl();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getProjectId();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string subject = 1;
 * @return {string}
 */
proto.grpc.service.webservice.GetMeResponse.prototype.getSubject = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.GetMeResponse} returns this
 */
proto.grpc.service.webservice.GetMeResponse.prototype.setSubject = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string avatar_url = 2;
 * @return {string}
 */
proto.grpc.service.webservice.GetMeResponse.prototype.getAvatarUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.GetMeResponse} returns this
 */
proto.grpc.service.webservice.GetMeResponse.prototype.setAvatarUrl = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string project_id = 3;
 * @return {string}
 */
proto.grpc.service.webservice.GetMeResponse.prototype.getProjectId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.GetMeResponse} returns this
 */
proto.grpc.service.webservice.GetMeResponse.prototype.setProjectId = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.grpc.service.webservice.AddProjectRBACRoleRequest.repeatedFields_ = [2];



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
proto.grpc.service.webservice.AddProjectRBACRoleRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.AddProjectRBACRoleRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.AddProjectRBACRoleRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.AddProjectRBACRoleRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    policiesList: jspb.Message.toObjectList(msg.getPoliciesList(),
    pkg_model_project_pb.ProjectRBACPolicy.toObject, includeInstance)
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
 * @return {!proto.grpc.service.webservice.AddProjectRBACRoleRequest}
 */
proto.grpc.service.webservice.AddProjectRBACRoleRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.AddProjectRBACRoleRequest;
  return proto.grpc.service.webservice.AddProjectRBACRoleRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.AddProjectRBACRoleRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.AddProjectRBACRoleRequest}
 */
proto.grpc.service.webservice.AddProjectRBACRoleRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new pkg_model_project_pb.ProjectRBACPolicy;
      reader.readMessage(value,pkg_model_project_pb.ProjectRBACPolicy.deserializeBinaryFromReader);
      msg.addPolicies(value);
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
proto.grpc.service.webservice.AddProjectRBACRoleRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.AddProjectRBACRoleRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.AddProjectRBACRoleRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.AddProjectRBACRoleRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPoliciesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      pkg_model_project_pb.ProjectRBACPolicy.serializeBinaryToWriter
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.grpc.service.webservice.AddProjectRBACRoleRequest.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.AddProjectRBACRoleRequest} returns this
 */
proto.grpc.service.webservice.AddProjectRBACRoleRequest.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * repeated model.ProjectRBACPolicy policies = 2;
 * @return {!Array<!proto.model.ProjectRBACPolicy>}
 */
proto.grpc.service.webservice.AddProjectRBACRoleRequest.prototype.getPoliciesList = function() {
  return /** @type{!Array<!proto.model.ProjectRBACPolicy>} */ (
    jspb.Message.getRepeatedWrapperField(this, pkg_model_project_pb.ProjectRBACPolicy, 2));
};


/**
 * @param {!Array<!proto.model.ProjectRBACPolicy>} value
 * @return {!proto.grpc.service.webservice.AddProjectRBACRoleRequest} returns this
*/
proto.grpc.service.webservice.AddProjectRBACRoleRequest.prototype.setPoliciesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.model.ProjectRBACPolicy=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.ProjectRBACPolicy}
 */
proto.grpc.service.webservice.AddProjectRBACRoleRequest.prototype.addPolicies = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.model.ProjectRBACPolicy, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.AddProjectRBACRoleRequest} returns this
 */
proto.grpc.service.webservice.AddProjectRBACRoleRequest.prototype.clearPoliciesList = function() {
  return this.setPoliciesList([]);
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
proto.grpc.service.webservice.AddProjectRBACRoleResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.AddProjectRBACRoleResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.AddProjectRBACRoleResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.AddProjectRBACRoleResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.AddProjectRBACRoleResponse}
 */
proto.grpc.service.webservice.AddProjectRBACRoleResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.AddProjectRBACRoleResponse;
  return proto.grpc.service.webservice.AddProjectRBACRoleResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.AddProjectRBACRoleResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.AddProjectRBACRoleResponse}
 */
proto.grpc.service.webservice.AddProjectRBACRoleResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.AddProjectRBACRoleResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.AddProjectRBACRoleResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.AddProjectRBACRoleResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.AddProjectRBACRoleResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.repeatedFields_ = [2];



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
proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.UpdateProjectRBACRoleRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    policiesList: jspb.Message.toObjectList(msg.getPoliciesList(),
    pkg_model_project_pb.ProjectRBACPolicy.toObject, includeInstance)
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
 * @return {!proto.grpc.service.webservice.UpdateProjectRBACRoleRequest}
 */
proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.UpdateProjectRBACRoleRequest;
  return proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.UpdateProjectRBACRoleRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.UpdateProjectRBACRoleRequest}
 */
proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new pkg_model_project_pb.ProjectRBACPolicy;
      reader.readMessage(value,pkg_model_project_pb.ProjectRBACPolicy.deserializeBinaryFromReader);
      msg.addPolicies(value);
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
proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.UpdateProjectRBACRoleRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPoliciesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      pkg_model_project_pb.ProjectRBACPolicy.serializeBinaryToWriter
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.UpdateProjectRBACRoleRequest} returns this
 */
proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * repeated model.ProjectRBACPolicy policies = 2;
 * @return {!Array<!proto.model.ProjectRBACPolicy>}
 */
proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.prototype.getPoliciesList = function() {
  return /** @type{!Array<!proto.model.ProjectRBACPolicy>} */ (
    jspb.Message.getRepeatedWrapperField(this, pkg_model_project_pb.ProjectRBACPolicy, 2));
};


/**
 * @param {!Array<!proto.model.ProjectRBACPolicy>} value
 * @return {!proto.grpc.service.webservice.UpdateProjectRBACRoleRequest} returns this
*/
proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.prototype.setPoliciesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.model.ProjectRBACPolicy=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.ProjectRBACPolicy}
 */
proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.prototype.addPolicies = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.model.ProjectRBACPolicy, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.UpdateProjectRBACRoleRequest} returns this
 */
proto.grpc.service.webservice.UpdateProjectRBACRoleRequest.prototype.clearPoliciesList = function() {
  return this.setPoliciesList([]);
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
proto.grpc.service.webservice.UpdateProjectRBACRoleResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.UpdateProjectRBACRoleResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.UpdateProjectRBACRoleResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateProjectRBACRoleResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.UpdateProjectRBACRoleResponse}
 */
proto.grpc.service.webservice.UpdateProjectRBACRoleResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.UpdateProjectRBACRoleResponse;
  return proto.grpc.service.webservice.UpdateProjectRBACRoleResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.UpdateProjectRBACRoleResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.UpdateProjectRBACRoleResponse}
 */
proto.grpc.service.webservice.UpdateProjectRBACRoleResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.UpdateProjectRBACRoleResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.UpdateProjectRBACRoleResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.UpdateProjectRBACRoleResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.UpdateProjectRBACRoleResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.DeleteProjectRBACRoleRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.DeleteProjectRBACRoleRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.DeleteProjectRBACRoleRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DeleteProjectRBACRoleRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.DeleteProjectRBACRoleRequest}
 */
proto.grpc.service.webservice.DeleteProjectRBACRoleRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.DeleteProjectRBACRoleRequest;
  return proto.grpc.service.webservice.DeleteProjectRBACRoleRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.DeleteProjectRBACRoleRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.DeleteProjectRBACRoleRequest}
 */
proto.grpc.service.webservice.DeleteProjectRBACRoleRequest.deserializeBinaryFromReader = function(msg, reader) {
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
proto.grpc.service.webservice.DeleteProjectRBACRoleRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.DeleteProjectRBACRoleRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.DeleteProjectRBACRoleRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DeleteProjectRBACRoleRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.grpc.service.webservice.DeleteProjectRBACRoleRequest.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.DeleteProjectRBACRoleRequest} returns this
 */
proto.grpc.service.webservice.DeleteProjectRBACRoleRequest.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.DeleteProjectRBACRoleResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.DeleteProjectRBACRoleResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.DeleteProjectRBACRoleResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DeleteProjectRBACRoleResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.DeleteProjectRBACRoleResponse}
 */
proto.grpc.service.webservice.DeleteProjectRBACRoleResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.DeleteProjectRBACRoleResponse;
  return proto.grpc.service.webservice.DeleteProjectRBACRoleResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.DeleteProjectRBACRoleResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.DeleteProjectRBACRoleResponse}
 */
proto.grpc.service.webservice.DeleteProjectRBACRoleResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.DeleteProjectRBACRoleResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.DeleteProjectRBACRoleResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.DeleteProjectRBACRoleResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DeleteProjectRBACRoleResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.AddProjectUserGroupRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.AddProjectUserGroupRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.AddProjectUserGroupRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.AddProjectUserGroupRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    ssoGroup: jspb.Message.getFieldWithDefault(msg, 1, ""),
    role: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.grpc.service.webservice.AddProjectUserGroupRequest}
 */
proto.grpc.service.webservice.AddProjectUserGroupRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.AddProjectUserGroupRequest;
  return proto.grpc.service.webservice.AddProjectUserGroupRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.AddProjectUserGroupRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.AddProjectUserGroupRequest}
 */
proto.grpc.service.webservice.AddProjectUserGroupRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setSsoGroup(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setRole(value);
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
proto.grpc.service.webservice.AddProjectUserGroupRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.AddProjectUserGroupRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.AddProjectUserGroupRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.AddProjectUserGroupRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSsoGroup();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRole();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string sso_group = 1;
 * @return {string}
 */
proto.grpc.service.webservice.AddProjectUserGroupRequest.prototype.getSsoGroup = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.AddProjectUserGroupRequest} returns this
 */
proto.grpc.service.webservice.AddProjectUserGroupRequest.prototype.setSsoGroup = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string role = 2;
 * @return {string}
 */
proto.grpc.service.webservice.AddProjectUserGroupRequest.prototype.getRole = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.AddProjectUserGroupRequest} returns this
 */
proto.grpc.service.webservice.AddProjectUserGroupRequest.prototype.setRole = function(value) {
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
proto.grpc.service.webservice.AddProjectUserGroupResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.AddProjectUserGroupResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.AddProjectUserGroupResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.AddProjectUserGroupResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.AddProjectUserGroupResponse}
 */
proto.grpc.service.webservice.AddProjectUserGroupResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.AddProjectUserGroupResponse;
  return proto.grpc.service.webservice.AddProjectUserGroupResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.AddProjectUserGroupResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.AddProjectUserGroupResponse}
 */
proto.grpc.service.webservice.AddProjectUserGroupResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.AddProjectUserGroupResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.AddProjectUserGroupResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.AddProjectUserGroupResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.AddProjectUserGroupResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.DeleteProjectUserGroupRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.DeleteProjectUserGroupRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.DeleteProjectUserGroupRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DeleteProjectUserGroupRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    ssoGroup: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.DeleteProjectUserGroupRequest}
 */
proto.grpc.service.webservice.DeleteProjectUserGroupRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.DeleteProjectUserGroupRequest;
  return proto.grpc.service.webservice.DeleteProjectUserGroupRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.DeleteProjectUserGroupRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.DeleteProjectUserGroupRequest}
 */
proto.grpc.service.webservice.DeleteProjectUserGroupRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setSsoGroup(value);
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
proto.grpc.service.webservice.DeleteProjectUserGroupRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.DeleteProjectUserGroupRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.DeleteProjectUserGroupRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DeleteProjectUserGroupRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSsoGroup();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string sso_group = 1;
 * @return {string}
 */
proto.grpc.service.webservice.DeleteProjectUserGroupRequest.prototype.getSsoGroup = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.DeleteProjectUserGroupRequest} returns this
 */
proto.grpc.service.webservice.DeleteProjectUserGroupRequest.prototype.setSsoGroup = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.DeleteProjectUserGroupResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.DeleteProjectUserGroupResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.DeleteProjectUserGroupResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DeleteProjectUserGroupResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.DeleteProjectUserGroupResponse}
 */
proto.grpc.service.webservice.DeleteProjectUserGroupResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.DeleteProjectUserGroupResponse;
  return proto.grpc.service.webservice.DeleteProjectUserGroupResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.DeleteProjectUserGroupResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.DeleteProjectUserGroupResponse}
 */
proto.grpc.service.webservice.DeleteProjectUserGroupResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.DeleteProjectUserGroupResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.DeleteProjectUserGroupResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.DeleteProjectUserGroupResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DeleteProjectUserGroupResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.GetCommandRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetCommandRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetCommandRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetCommandRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    commandId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.GetCommandRequest}
 */
proto.grpc.service.webservice.GetCommandRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetCommandRequest;
  return proto.grpc.service.webservice.GetCommandRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetCommandRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetCommandRequest}
 */
proto.grpc.service.webservice.GetCommandRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setCommandId(value);
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
proto.grpc.service.webservice.GetCommandRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetCommandRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetCommandRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetCommandRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCommandId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string command_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.GetCommandRequest.prototype.getCommandId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.GetCommandRequest} returns this
 */
proto.grpc.service.webservice.GetCommandRequest.prototype.setCommandId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.GetCommandResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetCommandResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetCommandResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetCommandResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    command: (f = msg.getCommand()) && pkg_model_command_pb.Command.toObject(includeInstance, f)
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
 * @return {!proto.grpc.service.webservice.GetCommandResponse}
 */
proto.grpc.service.webservice.GetCommandResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetCommandResponse;
  return proto.grpc.service.webservice.GetCommandResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetCommandResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetCommandResponse}
 */
proto.grpc.service.webservice.GetCommandResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new pkg_model_command_pb.Command;
      reader.readMessage(value,pkg_model_command_pb.Command.deserializeBinaryFromReader);
      msg.setCommand(value);
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
proto.grpc.service.webservice.GetCommandResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetCommandResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetCommandResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetCommandResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCommand();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      pkg_model_command_pb.Command.serializeBinaryToWriter
    );
  }
};


/**
 * optional model.Command command = 1;
 * @return {?proto.model.Command}
 */
proto.grpc.service.webservice.GetCommandResponse.prototype.getCommand = function() {
  return /** @type{?proto.model.Command} */ (
    jspb.Message.getWrapperField(this, pkg_model_command_pb.Command, 1));
};


/**
 * @param {?proto.model.Command|undefined} value
 * @return {!proto.grpc.service.webservice.GetCommandResponse} returns this
*/
proto.grpc.service.webservice.GetCommandResponse.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.GetCommandResponse} returns this
 */
proto.grpc.service.webservice.GetCommandResponse.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.GetCommandResponse.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 1) != null;
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
proto.grpc.service.webservice.GenerateAPIKeyRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GenerateAPIKeyRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GenerateAPIKeyRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GenerateAPIKeyRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    role: jspb.Message.getFieldWithDefault(msg, 2, 0)
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
 * @return {!proto.grpc.service.webservice.GenerateAPIKeyRequest}
 */
proto.grpc.service.webservice.GenerateAPIKeyRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GenerateAPIKeyRequest;
  return proto.grpc.service.webservice.GenerateAPIKeyRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GenerateAPIKeyRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GenerateAPIKeyRequest}
 */
proto.grpc.service.webservice.GenerateAPIKeyRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {!proto.model.APIKey.Role} */ (reader.readEnum());
      msg.setRole(value);
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
proto.grpc.service.webservice.GenerateAPIKeyRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GenerateAPIKeyRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GenerateAPIKeyRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GenerateAPIKeyRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRole();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.grpc.service.webservice.GenerateAPIKeyRequest.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.GenerateAPIKeyRequest} returns this
 */
proto.grpc.service.webservice.GenerateAPIKeyRequest.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional model.APIKey.Role role = 2;
 * @return {!proto.model.APIKey.Role}
 */
proto.grpc.service.webservice.GenerateAPIKeyRequest.prototype.getRole = function() {
  return /** @type {!proto.model.APIKey.Role} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.model.APIKey.Role} value
 * @return {!proto.grpc.service.webservice.GenerateAPIKeyRequest} returns this
 */
proto.grpc.service.webservice.GenerateAPIKeyRequest.prototype.setRole = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
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
proto.grpc.service.webservice.GenerateAPIKeyResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GenerateAPIKeyResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GenerateAPIKeyResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GenerateAPIKeyResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    key: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.GenerateAPIKeyResponse}
 */
proto.grpc.service.webservice.GenerateAPIKeyResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GenerateAPIKeyResponse;
  return proto.grpc.service.webservice.GenerateAPIKeyResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GenerateAPIKeyResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GenerateAPIKeyResponse}
 */
proto.grpc.service.webservice.GenerateAPIKeyResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setKey(value);
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
proto.grpc.service.webservice.GenerateAPIKeyResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GenerateAPIKeyResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GenerateAPIKeyResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GenerateAPIKeyResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getKey();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string key = 1;
 * @return {string}
 */
proto.grpc.service.webservice.GenerateAPIKeyResponse.prototype.getKey = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.GenerateAPIKeyResponse} returns this
 */
proto.grpc.service.webservice.GenerateAPIKeyResponse.prototype.setKey = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.DisableAPIKeyRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.DisableAPIKeyRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.DisableAPIKeyRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DisableAPIKeyRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.DisableAPIKeyRequest}
 */
proto.grpc.service.webservice.DisableAPIKeyRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.DisableAPIKeyRequest;
  return proto.grpc.service.webservice.DisableAPIKeyRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.DisableAPIKeyRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.DisableAPIKeyRequest}
 */
proto.grpc.service.webservice.DisableAPIKeyRequest.deserializeBinaryFromReader = function(msg, reader) {
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
proto.grpc.service.webservice.DisableAPIKeyRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.DisableAPIKeyRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.DisableAPIKeyRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DisableAPIKeyRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.DisableAPIKeyRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.DisableAPIKeyRequest} returns this
 */
proto.grpc.service.webservice.DisableAPIKeyRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.DisableAPIKeyResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.DisableAPIKeyResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.DisableAPIKeyResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DisableAPIKeyResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.DisableAPIKeyResponse}
 */
proto.grpc.service.webservice.DisableAPIKeyResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.DisableAPIKeyResponse;
  return proto.grpc.service.webservice.DisableAPIKeyResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.DisableAPIKeyResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.DisableAPIKeyResponse}
 */
proto.grpc.service.webservice.DisableAPIKeyResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.DisableAPIKeyResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.DisableAPIKeyResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.DisableAPIKeyResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.DisableAPIKeyResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
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
proto.grpc.service.webservice.ListAPIKeysRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListAPIKeysRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListAPIKeysRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListAPIKeysRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    options: (f = msg.getOptions()) && proto.grpc.service.webservice.ListAPIKeysRequest.Options.toObject(includeInstance, f)
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
 * @return {!proto.grpc.service.webservice.ListAPIKeysRequest}
 */
proto.grpc.service.webservice.ListAPIKeysRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListAPIKeysRequest;
  return proto.grpc.service.webservice.ListAPIKeysRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListAPIKeysRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListAPIKeysRequest}
 */
proto.grpc.service.webservice.ListAPIKeysRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 2:
      var value = new proto.grpc.service.webservice.ListAPIKeysRequest.Options;
      reader.readMessage(value,proto.grpc.service.webservice.ListAPIKeysRequest.Options.deserializeBinaryFromReader);
      msg.setOptions(value);
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
proto.grpc.service.webservice.ListAPIKeysRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListAPIKeysRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListAPIKeysRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListAPIKeysRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOptions();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.grpc.service.webservice.ListAPIKeysRequest.Options.serializeBinaryToWriter
    );
  }
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
proto.grpc.service.webservice.ListAPIKeysRequest.Options.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListAPIKeysRequest.Options.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListAPIKeysRequest.Options} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListAPIKeysRequest.Options.toObject = function(includeInstance, msg) {
  var f, obj = {
    enabled: (f = msg.getEnabled()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f)
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
 * @return {!proto.grpc.service.webservice.ListAPIKeysRequest.Options}
 */
proto.grpc.service.webservice.ListAPIKeysRequest.Options.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListAPIKeysRequest.Options;
  return proto.grpc.service.webservice.ListAPIKeysRequest.Options.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListAPIKeysRequest.Options} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListAPIKeysRequest.Options}
 */
proto.grpc.service.webservice.ListAPIKeysRequest.Options.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setEnabled(value);
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
proto.grpc.service.webservice.ListAPIKeysRequest.Options.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListAPIKeysRequest.Options.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListAPIKeysRequest.Options} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListAPIKeysRequest.Options.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEnabled();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
};


/**
 * optional google.protobuf.BoolValue enabled = 1;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.grpc.service.webservice.ListAPIKeysRequest.Options.prototype.getEnabled = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 1));
};


/**
 * @param {?proto.google.protobuf.BoolValue|undefined} value
 * @return {!proto.grpc.service.webservice.ListAPIKeysRequest.Options} returns this
*/
proto.grpc.service.webservice.ListAPIKeysRequest.Options.prototype.setEnabled = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.ListAPIKeysRequest.Options} returns this
 */
proto.grpc.service.webservice.ListAPIKeysRequest.Options.prototype.clearEnabled = function() {
  return this.setEnabled(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.ListAPIKeysRequest.Options.prototype.hasEnabled = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional Options options = 2;
 * @return {?proto.grpc.service.webservice.ListAPIKeysRequest.Options}
 */
proto.grpc.service.webservice.ListAPIKeysRequest.prototype.getOptions = function() {
  return /** @type{?proto.grpc.service.webservice.ListAPIKeysRequest.Options} */ (
    jspb.Message.getWrapperField(this, proto.grpc.service.webservice.ListAPIKeysRequest.Options, 2));
};


/**
 * @param {?proto.grpc.service.webservice.ListAPIKeysRequest.Options|undefined} value
 * @return {!proto.grpc.service.webservice.ListAPIKeysRequest} returns this
*/
proto.grpc.service.webservice.ListAPIKeysRequest.prototype.setOptions = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.ListAPIKeysRequest} returns this
 */
proto.grpc.service.webservice.ListAPIKeysRequest.prototype.clearOptions = function() {
  return this.setOptions(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.ListAPIKeysRequest.prototype.hasOptions = function() {
  return jspb.Message.getField(this, 2) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.grpc.service.webservice.ListAPIKeysResponse.repeatedFields_ = [1];



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
proto.grpc.service.webservice.ListAPIKeysResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListAPIKeysResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListAPIKeysResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListAPIKeysResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    keysList: jspb.Message.toObjectList(msg.getKeysList(),
    pkg_model_apikey_pb.APIKey.toObject, includeInstance)
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
 * @return {!proto.grpc.service.webservice.ListAPIKeysResponse}
 */
proto.grpc.service.webservice.ListAPIKeysResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListAPIKeysResponse;
  return proto.grpc.service.webservice.ListAPIKeysResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListAPIKeysResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListAPIKeysResponse}
 */
proto.grpc.service.webservice.ListAPIKeysResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new pkg_model_apikey_pb.APIKey;
      reader.readMessage(value,pkg_model_apikey_pb.APIKey.deserializeBinaryFromReader);
      msg.addKeys(value);
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
proto.grpc.service.webservice.ListAPIKeysResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListAPIKeysResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListAPIKeysResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListAPIKeysResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getKeysList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      pkg_model_apikey_pb.APIKey.serializeBinaryToWriter
    );
  }
};


/**
 * repeated model.APIKey keys = 1;
 * @return {!Array<!proto.model.APIKey>}
 */
proto.grpc.service.webservice.ListAPIKeysResponse.prototype.getKeysList = function() {
  return /** @type{!Array<!proto.model.APIKey>} */ (
    jspb.Message.getRepeatedWrapperField(this, pkg_model_apikey_pb.APIKey, 1));
};


/**
 * @param {!Array<!proto.model.APIKey>} value
 * @return {!proto.grpc.service.webservice.ListAPIKeysResponse} returns this
*/
proto.grpc.service.webservice.ListAPIKeysResponse.prototype.setKeysList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.model.APIKey=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.APIKey}
 */
proto.grpc.service.webservice.ListAPIKeysResponse.prototype.addKeys = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.model.APIKey, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.ListAPIKeysResponse} returns this
 */
proto.grpc.service.webservice.ListAPIKeysResponse.prototype.clearKeysList = function() {
  return this.setKeysList([]);
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
proto.grpc.service.webservice.GetInsightDataRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetInsightDataRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetInsightDataRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetInsightDataRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    metricsKind: jspb.Message.getFieldWithDefault(msg, 1, 0),
    rangeFrom: jspb.Message.getFieldWithDefault(msg, 2, 0),
    rangeTo: jspb.Message.getFieldWithDefault(msg, 3, 0),
    resolution: jspb.Message.getFieldWithDefault(msg, 4, 0),
    applicationId: jspb.Message.getFieldWithDefault(msg, 10, ""),
    labelsMap: (f = msg.getLabelsMap()) ? f.toObject(includeInstance, undefined) : []
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
 * @return {!proto.grpc.service.webservice.GetInsightDataRequest}
 */
proto.grpc.service.webservice.GetInsightDataRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetInsightDataRequest;
  return proto.grpc.service.webservice.GetInsightDataRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetInsightDataRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetInsightDataRequest}
 */
proto.grpc.service.webservice.GetInsightDataRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!proto.model.InsightMetricsKind} */ (reader.readEnum());
      msg.setMetricsKind(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setRangeFrom(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setRangeTo(value);
      break;
    case 4:
      var value = /** @type {!proto.model.InsightResolution} */ (reader.readEnum());
      msg.setResolution(value);
      break;
    case 10:
      var value = /** @type {string} */ (reader.readString());
      msg.setApplicationId(value);
      break;
    case 11:
      var value = msg.getLabelsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "", "");
         });
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
proto.grpc.service.webservice.GetInsightDataRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetInsightDataRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetInsightDataRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetInsightDataRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getMetricsKind();
  if (f !== 0.0) {
    writer.writeEnum(
      1,
      f
    );
  }
  f = message.getRangeFrom();
  if (f !== 0) {
    writer.writeInt64(
      2,
      f
    );
  }
  f = message.getRangeTo();
  if (f !== 0) {
    writer.writeInt64(
      3,
      f
    );
  }
  f = message.getResolution();
  if (f !== 0.0) {
    writer.writeEnum(
      4,
      f
    );
  }
  f = message.getApplicationId();
  if (f.length > 0) {
    writer.writeString(
      10,
      f
    );
  }
  f = message.getLabelsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(11, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
};


/**
 * optional model.InsightMetricsKind metrics_kind = 1;
 * @return {!proto.model.InsightMetricsKind}
 */
proto.grpc.service.webservice.GetInsightDataRequest.prototype.getMetricsKind = function() {
  return /** @type {!proto.model.InsightMetricsKind} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {!proto.model.InsightMetricsKind} value
 * @return {!proto.grpc.service.webservice.GetInsightDataRequest} returns this
 */
proto.grpc.service.webservice.GetInsightDataRequest.prototype.setMetricsKind = function(value) {
  return jspb.Message.setProto3EnumField(this, 1, value);
};


/**
 * optional int64 range_from = 2;
 * @return {number}
 */
proto.grpc.service.webservice.GetInsightDataRequest.prototype.getRangeFrom = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.grpc.service.webservice.GetInsightDataRequest} returns this
 */
proto.grpc.service.webservice.GetInsightDataRequest.prototype.setRangeFrom = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional int64 range_to = 3;
 * @return {number}
 */
proto.grpc.service.webservice.GetInsightDataRequest.prototype.getRangeTo = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.grpc.service.webservice.GetInsightDataRequest} returns this
 */
proto.grpc.service.webservice.GetInsightDataRequest.prototype.setRangeTo = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * optional model.InsightResolution resolution = 4;
 * @return {!proto.model.InsightResolution}
 */
proto.grpc.service.webservice.GetInsightDataRequest.prototype.getResolution = function() {
  return /** @type {!proto.model.InsightResolution} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {!proto.model.InsightResolution} value
 * @return {!proto.grpc.service.webservice.GetInsightDataRequest} returns this
 */
proto.grpc.service.webservice.GetInsightDataRequest.prototype.setResolution = function(value) {
  return jspb.Message.setProto3EnumField(this, 4, value);
};


/**
 * optional string application_id = 10;
 * @return {string}
 */
proto.grpc.service.webservice.GetInsightDataRequest.prototype.getApplicationId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 10, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.GetInsightDataRequest} returns this
 */
proto.grpc.service.webservice.GetInsightDataRequest.prototype.setApplicationId = function(value) {
  return jspb.Message.setProto3StringField(this, 10, value);
};


/**
 * map<string, string> labels = 11;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.grpc.service.webservice.GetInsightDataRequest.prototype.getLabelsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 11, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.grpc.service.webservice.GetInsightDataRequest} returns this
 */
proto.grpc.service.webservice.GetInsightDataRequest.prototype.clearLabelsMap = function() {
  this.getLabelsMap().clear();
  return this;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.grpc.service.webservice.GetInsightDataResponse.repeatedFields_ = [3,4];



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
proto.grpc.service.webservice.GetInsightDataResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetInsightDataResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetInsightDataResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetInsightDataResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    updatedAt: jspb.Message.getFieldWithDefault(msg, 1, 0),
    type: jspb.Message.getFieldWithDefault(msg, 2, 0),
    vectorList: jspb.Message.toObjectList(msg.getVectorList(),
    pkg_model_insight_pb.InsightSample.toObject, includeInstance),
    matrixList: jspb.Message.toObjectList(msg.getMatrixList(),
    pkg_model_insight_pb.InsightSampleStream.toObject, includeInstance)
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
 * @return {!proto.grpc.service.webservice.GetInsightDataResponse}
 */
proto.grpc.service.webservice.GetInsightDataResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetInsightDataResponse;
  return proto.grpc.service.webservice.GetInsightDataResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetInsightDataResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetInsightDataResponse}
 */
proto.grpc.service.webservice.GetInsightDataResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setUpdatedAt(value);
      break;
    case 2:
      var value = /** @type {!proto.model.InsightResultType} */ (reader.readEnum());
      msg.setType(value);
      break;
    case 3:
      var value = new pkg_model_insight_pb.InsightSample;
      reader.readMessage(value,pkg_model_insight_pb.InsightSample.deserializeBinaryFromReader);
      msg.addVector(value);
      break;
    case 4:
      var value = new pkg_model_insight_pb.InsightSampleStream;
      reader.readMessage(value,pkg_model_insight_pb.InsightSampleStream.deserializeBinaryFromReader);
      msg.addMatrix(value);
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
proto.grpc.service.webservice.GetInsightDataResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetInsightDataResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetInsightDataResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetInsightDataResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      1,
      f
    );
  }
  f = message.getType();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getVectorList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      3,
      f,
      pkg_model_insight_pb.InsightSample.serializeBinaryToWriter
    );
  }
  f = message.getMatrixList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      pkg_model_insight_pb.InsightSampleStream.serializeBinaryToWriter
    );
  }
};


/**
 * optional int64 updated_at = 1;
 * @return {number}
 */
proto.grpc.service.webservice.GetInsightDataResponse.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.grpc.service.webservice.GetInsightDataResponse} returns this
 */
proto.grpc.service.webservice.GetInsightDataResponse.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional model.InsightResultType type = 2;
 * @return {!proto.model.InsightResultType}
 */
proto.grpc.service.webservice.GetInsightDataResponse.prototype.getType = function() {
  return /** @type {!proto.model.InsightResultType} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.model.InsightResultType} value
 * @return {!proto.grpc.service.webservice.GetInsightDataResponse} returns this
 */
proto.grpc.service.webservice.GetInsightDataResponse.prototype.setType = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * repeated model.InsightSample vector = 3;
 * @return {!Array<!proto.model.InsightSample>}
 */
proto.grpc.service.webservice.GetInsightDataResponse.prototype.getVectorList = function() {
  return /** @type{!Array<!proto.model.InsightSample>} */ (
    jspb.Message.getRepeatedWrapperField(this, pkg_model_insight_pb.InsightSample, 3));
};


/**
 * @param {!Array<!proto.model.InsightSample>} value
 * @return {!proto.grpc.service.webservice.GetInsightDataResponse} returns this
*/
proto.grpc.service.webservice.GetInsightDataResponse.prototype.setVectorList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 3, value);
};


/**
 * @param {!proto.model.InsightSample=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.InsightSample}
 */
proto.grpc.service.webservice.GetInsightDataResponse.prototype.addVector = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 3, opt_value, proto.model.InsightSample, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.GetInsightDataResponse} returns this
 */
proto.grpc.service.webservice.GetInsightDataResponse.prototype.clearVectorList = function() {
  return this.setVectorList([]);
};


/**
 * repeated model.InsightSampleStream matrix = 4;
 * @return {!Array<!proto.model.InsightSampleStream>}
 */
proto.grpc.service.webservice.GetInsightDataResponse.prototype.getMatrixList = function() {
  return /** @type{!Array<!proto.model.InsightSampleStream>} */ (
    jspb.Message.getRepeatedWrapperField(this, pkg_model_insight_pb.InsightSampleStream, 4));
};


/**
 * @param {!Array<!proto.model.InsightSampleStream>} value
 * @return {!proto.grpc.service.webservice.GetInsightDataResponse} returns this
*/
proto.grpc.service.webservice.GetInsightDataResponse.prototype.setMatrixList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.model.InsightSampleStream=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.InsightSampleStream}
 */
proto.grpc.service.webservice.GetInsightDataResponse.prototype.addMatrix = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.model.InsightSampleStream, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.GetInsightDataResponse} returns this
 */
proto.grpc.service.webservice.GetInsightDataResponse.prototype.clearMatrixList = function() {
  return this.setMatrixList([]);
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
proto.grpc.service.webservice.GetInsightApplicationCountRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetInsightApplicationCountRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetInsightApplicationCountRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetInsightApplicationCountRequest.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.GetInsightApplicationCountRequest}
 */
proto.grpc.service.webservice.GetInsightApplicationCountRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetInsightApplicationCountRequest;
  return proto.grpc.service.webservice.GetInsightApplicationCountRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetInsightApplicationCountRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetInsightApplicationCountRequest}
 */
proto.grpc.service.webservice.GetInsightApplicationCountRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.GetInsightApplicationCountRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetInsightApplicationCountRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetInsightApplicationCountRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetInsightApplicationCountRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.grpc.service.webservice.GetInsightApplicationCountResponse.repeatedFields_ = [2];



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
proto.grpc.service.webservice.GetInsightApplicationCountResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetInsightApplicationCountResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetInsightApplicationCountResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetInsightApplicationCountResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    updatedAt: jspb.Message.getFieldWithDefault(msg, 1, 0),
    countsList: jspb.Message.toObjectList(msg.getCountsList(),
    pkg_model_insight_pb.InsightApplicationCount.toObject, includeInstance)
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
 * @return {!proto.grpc.service.webservice.GetInsightApplicationCountResponse}
 */
proto.grpc.service.webservice.GetInsightApplicationCountResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetInsightApplicationCountResponse;
  return proto.grpc.service.webservice.GetInsightApplicationCountResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetInsightApplicationCountResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetInsightApplicationCountResponse}
 */
proto.grpc.service.webservice.GetInsightApplicationCountResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setUpdatedAt(value);
      break;
    case 2:
      var value = new pkg_model_insight_pb.InsightApplicationCount;
      reader.readMessage(value,pkg_model_insight_pb.InsightApplicationCount.deserializeBinaryFromReader);
      msg.addCounts(value);
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
proto.grpc.service.webservice.GetInsightApplicationCountResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetInsightApplicationCountResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetInsightApplicationCountResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetInsightApplicationCountResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      1,
      f
    );
  }
  f = message.getCountsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      pkg_model_insight_pb.InsightApplicationCount.serializeBinaryToWriter
    );
  }
};


/**
 * optional int64 updated_at = 1;
 * @return {number}
 */
proto.grpc.service.webservice.GetInsightApplicationCountResponse.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.grpc.service.webservice.GetInsightApplicationCountResponse} returns this
 */
proto.grpc.service.webservice.GetInsightApplicationCountResponse.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * repeated model.InsightApplicationCount counts = 2;
 * @return {!Array<!proto.model.InsightApplicationCount>}
 */
proto.grpc.service.webservice.GetInsightApplicationCountResponse.prototype.getCountsList = function() {
  return /** @type{!Array<!proto.model.InsightApplicationCount>} */ (
    jspb.Message.getRepeatedWrapperField(this, pkg_model_insight_pb.InsightApplicationCount, 2));
};


/**
 * @param {!Array<!proto.model.InsightApplicationCount>} value
 * @return {!proto.grpc.service.webservice.GetInsightApplicationCountResponse} returns this
*/
proto.grpc.service.webservice.GetInsightApplicationCountResponse.prototype.setCountsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.model.InsightApplicationCount=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.InsightApplicationCount}
 */
proto.grpc.service.webservice.GetInsightApplicationCountResponse.prototype.addCounts = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.model.InsightApplicationCount, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.GetInsightApplicationCountResponse} returns this
 */
proto.grpc.service.webservice.GetInsightApplicationCountResponse.prototype.clearCountsList = function() {
  return this.setCountsList([]);
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
proto.grpc.service.webservice.ListDeploymentChainsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListDeploymentChainsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListDeploymentChainsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListDeploymentChainsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    options: (f = msg.getOptions()) && proto.grpc.service.webservice.ListDeploymentChainsRequest.Options.toObject(includeInstance, f),
    pageSize: jspb.Message.getFieldWithDefault(msg, 2, 0),
    cursor: jspb.Message.getFieldWithDefault(msg, 3, ""),
    pageMinUpdatedAt: jspb.Message.getFieldWithDefault(msg, 4, 0)
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
 * @return {!proto.grpc.service.webservice.ListDeploymentChainsRequest}
 */
proto.grpc.service.webservice.ListDeploymentChainsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListDeploymentChainsRequest;
  return proto.grpc.service.webservice.ListDeploymentChainsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListDeploymentChainsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListDeploymentChainsRequest}
 */
proto.grpc.service.webservice.ListDeploymentChainsRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.grpc.service.webservice.ListDeploymentChainsRequest.Options;
      reader.readMessage(value,proto.grpc.service.webservice.ListDeploymentChainsRequest.Options.deserializeBinaryFromReader);
      msg.setOptions(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setPageSize(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setCursor(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setPageMinUpdatedAt(value);
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
proto.grpc.service.webservice.ListDeploymentChainsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListDeploymentChainsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListDeploymentChainsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListDeploymentChainsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOptions();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.grpc.service.webservice.ListDeploymentChainsRequest.Options.serializeBinaryToWriter
    );
  }
  f = message.getPageSize();
  if (f !== 0) {
    writer.writeInt32(
      2,
      f
    );
  }
  f = message.getCursor();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getPageMinUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      4,
      f
    );
  }
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
proto.grpc.service.webservice.ListDeploymentChainsRequest.Options.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListDeploymentChainsRequest.Options.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListDeploymentChainsRequest.Options} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListDeploymentChainsRequest.Options.toObject = function(includeInstance, msg) {
  var f, obj = {

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
 * @return {!proto.grpc.service.webservice.ListDeploymentChainsRequest.Options}
 */
proto.grpc.service.webservice.ListDeploymentChainsRequest.Options.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListDeploymentChainsRequest.Options;
  return proto.grpc.service.webservice.ListDeploymentChainsRequest.Options.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListDeploymentChainsRequest.Options} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListDeploymentChainsRequest.Options}
 */
proto.grpc.service.webservice.ListDeploymentChainsRequest.Options.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.grpc.service.webservice.ListDeploymentChainsRequest.Options.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListDeploymentChainsRequest.Options.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListDeploymentChainsRequest.Options} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListDeploymentChainsRequest.Options.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};


/**
 * optional Options options = 1;
 * @return {?proto.grpc.service.webservice.ListDeploymentChainsRequest.Options}
 */
proto.grpc.service.webservice.ListDeploymentChainsRequest.prototype.getOptions = function() {
  return /** @type{?proto.grpc.service.webservice.ListDeploymentChainsRequest.Options} */ (
    jspb.Message.getWrapperField(this, proto.grpc.service.webservice.ListDeploymentChainsRequest.Options, 1));
};


/**
 * @param {?proto.grpc.service.webservice.ListDeploymentChainsRequest.Options|undefined} value
 * @return {!proto.grpc.service.webservice.ListDeploymentChainsRequest} returns this
*/
proto.grpc.service.webservice.ListDeploymentChainsRequest.prototype.setOptions = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.ListDeploymentChainsRequest} returns this
 */
proto.grpc.service.webservice.ListDeploymentChainsRequest.prototype.clearOptions = function() {
  return this.setOptions(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.ListDeploymentChainsRequest.prototype.hasOptions = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional int32 page_size = 2;
 * @return {number}
 */
proto.grpc.service.webservice.ListDeploymentChainsRequest.prototype.getPageSize = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.grpc.service.webservice.ListDeploymentChainsRequest} returns this
 */
proto.grpc.service.webservice.ListDeploymentChainsRequest.prototype.setPageSize = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional string cursor = 3;
 * @return {string}
 */
proto.grpc.service.webservice.ListDeploymentChainsRequest.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.ListDeploymentChainsRequest} returns this
 */
proto.grpc.service.webservice.ListDeploymentChainsRequest.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional int64 page_min_updated_at = 4;
 * @return {number}
 */
proto.grpc.service.webservice.ListDeploymentChainsRequest.prototype.getPageMinUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.grpc.service.webservice.ListDeploymentChainsRequest} returns this
 */
proto.grpc.service.webservice.ListDeploymentChainsRequest.prototype.setPageMinUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.grpc.service.webservice.ListDeploymentChainsResponse.repeatedFields_ = [1];



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
proto.grpc.service.webservice.ListDeploymentChainsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListDeploymentChainsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListDeploymentChainsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListDeploymentChainsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    deploymentChainsList: jspb.Message.toObjectList(msg.getDeploymentChainsList(),
    pkg_model_deployment_chain_pb.DeploymentChain.toObject, includeInstance),
    cursor: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.grpc.service.webservice.ListDeploymentChainsResponse}
 */
proto.grpc.service.webservice.ListDeploymentChainsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListDeploymentChainsResponse;
  return proto.grpc.service.webservice.ListDeploymentChainsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListDeploymentChainsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListDeploymentChainsResponse}
 */
proto.grpc.service.webservice.ListDeploymentChainsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new pkg_model_deployment_chain_pb.DeploymentChain;
      reader.readMessage(value,pkg_model_deployment_chain_pb.DeploymentChain.deserializeBinaryFromReader);
      msg.addDeploymentChains(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setCursor(value);
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
proto.grpc.service.webservice.ListDeploymentChainsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListDeploymentChainsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListDeploymentChainsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListDeploymentChainsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeploymentChainsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      pkg_model_deployment_chain_pb.DeploymentChain.serializeBinaryToWriter
    );
  }
  f = message.getCursor();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * repeated model.DeploymentChain deployment_chains = 1;
 * @return {!Array<!proto.model.DeploymentChain>}
 */
proto.grpc.service.webservice.ListDeploymentChainsResponse.prototype.getDeploymentChainsList = function() {
  return /** @type{!Array<!proto.model.DeploymentChain>} */ (
    jspb.Message.getRepeatedWrapperField(this, pkg_model_deployment_chain_pb.DeploymentChain, 1));
};


/**
 * @param {!Array<!proto.model.DeploymentChain>} value
 * @return {!proto.grpc.service.webservice.ListDeploymentChainsResponse} returns this
*/
proto.grpc.service.webservice.ListDeploymentChainsResponse.prototype.setDeploymentChainsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.model.DeploymentChain=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.DeploymentChain}
 */
proto.grpc.service.webservice.ListDeploymentChainsResponse.prototype.addDeploymentChains = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.model.DeploymentChain, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.ListDeploymentChainsResponse} returns this
 */
proto.grpc.service.webservice.ListDeploymentChainsResponse.prototype.clearDeploymentChainsList = function() {
  return this.setDeploymentChainsList([]);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.grpc.service.webservice.ListDeploymentChainsResponse.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.ListDeploymentChainsResponse} returns this
 */
proto.grpc.service.webservice.ListDeploymentChainsResponse.prototype.setCursor = function(value) {
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
proto.grpc.service.webservice.GetDeploymentChainRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetDeploymentChainRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetDeploymentChainRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetDeploymentChainRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    deploymentChainId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.grpc.service.webservice.GetDeploymentChainRequest}
 */
proto.grpc.service.webservice.GetDeploymentChainRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetDeploymentChainRequest;
  return proto.grpc.service.webservice.GetDeploymentChainRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetDeploymentChainRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetDeploymentChainRequest}
 */
proto.grpc.service.webservice.GetDeploymentChainRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setDeploymentChainId(value);
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
proto.grpc.service.webservice.GetDeploymentChainRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetDeploymentChainRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetDeploymentChainRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetDeploymentChainRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeploymentChainId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string deployment_chain_id = 1;
 * @return {string}
 */
proto.grpc.service.webservice.GetDeploymentChainRequest.prototype.getDeploymentChainId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.GetDeploymentChainRequest} returns this
 */
proto.grpc.service.webservice.GetDeploymentChainRequest.prototype.setDeploymentChainId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
proto.grpc.service.webservice.GetDeploymentChainResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.GetDeploymentChainResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.GetDeploymentChainResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetDeploymentChainResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    deploymentChain: (f = msg.getDeploymentChain()) && pkg_model_deployment_chain_pb.DeploymentChain.toObject(includeInstance, f)
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
 * @return {!proto.grpc.service.webservice.GetDeploymentChainResponse}
 */
proto.grpc.service.webservice.GetDeploymentChainResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.GetDeploymentChainResponse;
  return proto.grpc.service.webservice.GetDeploymentChainResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.GetDeploymentChainResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.GetDeploymentChainResponse}
 */
proto.grpc.service.webservice.GetDeploymentChainResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new pkg_model_deployment_chain_pb.DeploymentChain;
      reader.readMessage(value,pkg_model_deployment_chain_pb.DeploymentChain.deserializeBinaryFromReader);
      msg.setDeploymentChain(value);
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
proto.grpc.service.webservice.GetDeploymentChainResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.GetDeploymentChainResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.GetDeploymentChainResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.GetDeploymentChainResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeploymentChain();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      pkg_model_deployment_chain_pb.DeploymentChain.serializeBinaryToWriter
    );
  }
};


/**
 * optional model.DeploymentChain deployment_chain = 1;
 * @return {?proto.model.DeploymentChain}
 */
proto.grpc.service.webservice.GetDeploymentChainResponse.prototype.getDeploymentChain = function() {
  return /** @type{?proto.model.DeploymentChain} */ (
    jspb.Message.getWrapperField(this, pkg_model_deployment_chain_pb.DeploymentChain, 1));
};


/**
 * @param {?proto.model.DeploymentChain|undefined} value
 * @return {!proto.grpc.service.webservice.GetDeploymentChainResponse} returns this
*/
proto.grpc.service.webservice.GetDeploymentChainResponse.prototype.setDeploymentChain = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.GetDeploymentChainResponse} returns this
 */
proto.grpc.service.webservice.GetDeploymentChainResponse.prototype.clearDeploymentChain = function() {
  return this.setDeploymentChain(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.GetDeploymentChainResponse.prototype.hasDeploymentChain = function() {
  return jspb.Message.getField(this, 1) != null;
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
proto.grpc.service.webservice.ListEventsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListEventsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListEventsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListEventsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    options: (f = msg.getOptions()) && proto.grpc.service.webservice.ListEventsRequest.Options.toObject(includeInstance, f),
    pageSize: jspb.Message.getFieldWithDefault(msg, 2, 0),
    cursor: jspb.Message.getFieldWithDefault(msg, 3, ""),
    pageMinUpdatedAt: jspb.Message.getFieldWithDefault(msg, 4, 0)
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
 * @return {!proto.grpc.service.webservice.ListEventsRequest}
 */
proto.grpc.service.webservice.ListEventsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListEventsRequest;
  return proto.grpc.service.webservice.ListEventsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListEventsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListEventsRequest}
 */
proto.grpc.service.webservice.ListEventsRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.grpc.service.webservice.ListEventsRequest.Options;
      reader.readMessage(value,proto.grpc.service.webservice.ListEventsRequest.Options.deserializeBinaryFromReader);
      msg.setOptions(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setPageSize(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setCursor(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setPageMinUpdatedAt(value);
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
proto.grpc.service.webservice.ListEventsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListEventsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListEventsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListEventsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOptions();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.grpc.service.webservice.ListEventsRequest.Options.serializeBinaryToWriter
    );
  }
  f = message.getPageSize();
  if (f !== 0) {
    writer.writeInt32(
      2,
      f
    );
  }
  f = message.getCursor();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getPageMinUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      4,
      f
    );
  }
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.grpc.service.webservice.ListEventsRequest.Options.repeatedFields_ = [2];



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
proto.grpc.service.webservice.ListEventsRequest.Options.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListEventsRequest.Options.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListEventsRequest.Options} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListEventsRequest.Options.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    statusesList: (f = jspb.Message.getRepeatedField(msg, 2)) == null ? undefined : f,
    labelsMap: (f = msg.getLabelsMap()) ? f.toObject(includeInstance, undefined) : []
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
 * @return {!proto.grpc.service.webservice.ListEventsRequest.Options}
 */
proto.grpc.service.webservice.ListEventsRequest.Options.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListEventsRequest.Options;
  return proto.grpc.service.webservice.ListEventsRequest.Options.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListEventsRequest.Options} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListEventsRequest.Options}
 */
proto.grpc.service.webservice.ListEventsRequest.Options.deserializeBinaryFromReader = function(msg, reader) {
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
      var values = /** @type {!Array<!proto.model.EventStatus>} */ (reader.isDelimited() ? reader.readPackedEnum() : [reader.readEnum()]);
      for (var i = 0; i < values.length; i++) {
        msg.addStatuses(values[i]);
      }
      break;
    case 3:
      var value = msg.getLabelsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "", "");
         });
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
proto.grpc.service.webservice.ListEventsRequest.Options.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListEventsRequest.Options.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListEventsRequest.Options} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListEventsRequest.Options.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStatusesList();
  if (f.length > 0) {
    writer.writePackedEnum(
      2,
      f
    );
  }
  f = message.getLabelsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(3, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.grpc.service.webservice.ListEventsRequest.Options.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.ListEventsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListEventsRequest.Options.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * repeated model.EventStatus statuses = 2;
 * @return {!Array<!proto.model.EventStatus>}
 */
proto.grpc.service.webservice.ListEventsRequest.Options.prototype.getStatusesList = function() {
  return /** @type {!Array<!proto.model.EventStatus>} */ (jspb.Message.getRepeatedField(this, 2));
};


/**
 * @param {!Array<!proto.model.EventStatus>} value
 * @return {!proto.grpc.service.webservice.ListEventsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListEventsRequest.Options.prototype.setStatusesList = function(value) {
  return jspb.Message.setField(this, 2, value || []);
};


/**
 * @param {!proto.model.EventStatus} value
 * @param {number=} opt_index
 * @return {!proto.grpc.service.webservice.ListEventsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListEventsRequest.Options.prototype.addStatuses = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 2, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.ListEventsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListEventsRequest.Options.prototype.clearStatusesList = function() {
  return this.setStatusesList([]);
};


/**
 * map<string, string> labels = 3;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.grpc.service.webservice.ListEventsRequest.Options.prototype.getLabelsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 3, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.grpc.service.webservice.ListEventsRequest.Options} returns this
 */
proto.grpc.service.webservice.ListEventsRequest.Options.prototype.clearLabelsMap = function() {
  this.getLabelsMap().clear();
  return this;
};


/**
 * optional Options options = 1;
 * @return {?proto.grpc.service.webservice.ListEventsRequest.Options}
 */
proto.grpc.service.webservice.ListEventsRequest.prototype.getOptions = function() {
  return /** @type{?proto.grpc.service.webservice.ListEventsRequest.Options} */ (
    jspb.Message.getWrapperField(this, proto.grpc.service.webservice.ListEventsRequest.Options, 1));
};


/**
 * @param {?proto.grpc.service.webservice.ListEventsRequest.Options|undefined} value
 * @return {!proto.grpc.service.webservice.ListEventsRequest} returns this
*/
proto.grpc.service.webservice.ListEventsRequest.prototype.setOptions = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.grpc.service.webservice.ListEventsRequest} returns this
 */
proto.grpc.service.webservice.ListEventsRequest.prototype.clearOptions = function() {
  return this.setOptions(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.grpc.service.webservice.ListEventsRequest.prototype.hasOptions = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional int32 page_size = 2;
 * @return {number}
 */
proto.grpc.service.webservice.ListEventsRequest.prototype.getPageSize = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.grpc.service.webservice.ListEventsRequest} returns this
 */
proto.grpc.service.webservice.ListEventsRequest.prototype.setPageSize = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional string cursor = 3;
 * @return {string}
 */
proto.grpc.service.webservice.ListEventsRequest.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.ListEventsRequest} returns this
 */
proto.grpc.service.webservice.ListEventsRequest.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional int64 page_min_updated_at = 4;
 * @return {number}
 */
proto.grpc.service.webservice.ListEventsRequest.prototype.getPageMinUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.grpc.service.webservice.ListEventsRequest} returns this
 */
proto.grpc.service.webservice.ListEventsRequest.prototype.setPageMinUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.grpc.service.webservice.ListEventsResponse.repeatedFields_ = [1];



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
proto.grpc.service.webservice.ListEventsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.grpc.service.webservice.ListEventsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.grpc.service.webservice.ListEventsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListEventsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    eventsList: jspb.Message.toObjectList(msg.getEventsList(),
    pkg_model_event_pb.Event.toObject, includeInstance),
    cursor: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.grpc.service.webservice.ListEventsResponse}
 */
proto.grpc.service.webservice.ListEventsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.grpc.service.webservice.ListEventsResponse;
  return proto.grpc.service.webservice.ListEventsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.grpc.service.webservice.ListEventsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.grpc.service.webservice.ListEventsResponse}
 */
proto.grpc.service.webservice.ListEventsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new pkg_model_event_pb.Event;
      reader.readMessage(value,pkg_model_event_pb.Event.deserializeBinaryFromReader);
      msg.addEvents(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setCursor(value);
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
proto.grpc.service.webservice.ListEventsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.grpc.service.webservice.ListEventsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.grpc.service.webservice.ListEventsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.grpc.service.webservice.ListEventsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEventsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      pkg_model_event_pb.Event.serializeBinaryToWriter
    );
  }
  f = message.getCursor();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * repeated model.Event events = 1;
 * @return {!Array<!proto.model.Event>}
 */
proto.grpc.service.webservice.ListEventsResponse.prototype.getEventsList = function() {
  return /** @type{!Array<!proto.model.Event>} */ (
    jspb.Message.getRepeatedWrapperField(this, pkg_model_event_pb.Event, 1));
};


/**
 * @param {!Array<!proto.model.Event>} value
 * @return {!proto.grpc.service.webservice.ListEventsResponse} returns this
*/
proto.grpc.service.webservice.ListEventsResponse.prototype.setEventsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.model.Event=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.Event}
 */
proto.grpc.service.webservice.ListEventsResponse.prototype.addEvents = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.model.Event, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.grpc.service.webservice.ListEventsResponse} returns this
 */
proto.grpc.service.webservice.ListEventsResponse.prototype.clearEventsList = function() {
  return this.setEventsList([]);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.grpc.service.webservice.ListEventsResponse.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.grpc.service.webservice.ListEventsResponse} returns this
 */
proto.grpc.service.webservice.ListEventsResponse.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


goog.object.extend(exports, proto.grpc.service.webservice);
