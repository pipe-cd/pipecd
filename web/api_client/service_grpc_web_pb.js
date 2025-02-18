/**
 * @fileoverview gRPC-Web generated client stub for grpc.service.webservice
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');




var pkg_model_common_pb = require('pipecd/web/model/common_pb.js')

var pkg_model_insight_pb = require('pipecd/web/model/insight_pb.js')

var pkg_model_application_pb = require('pipecd/web/model/application_pb.js')

var pkg_model_application_live_state_pb = require('pipecd/web/model/application_live_state_pb.js')

var pkg_model_command_pb = require('pipecd/web/model/command_pb.js')

var pkg_model_deployment_pb = require('pipecd/web/model/deployment_pb.js')

var pkg_model_deployment_trace_pb = require('pipecd/web/model/deployment_trace_pb.js')

var pkg_model_deployment_chain_pb = require('pipecd/web/model/deployment_chain_pb.js')

var pkg_model_logblock_pb = require('pipecd/web/model/logblock_pb.js')

var pkg_model_piped_pb = require('pipecd/web/model/piped_pb.js')

var pkg_model_rbac_pb = require('pipecd/web/model/rbac_pb.js')

var pkg_model_project_pb = require('pipecd/web/model/project_pb.js')

var pkg_model_apikey_pb = require('pipecd/web/model/apikey_pb.js')

var pkg_model_event_pb = require('pipecd/web/model/event_pb.js')

var google_protobuf_wrappers_pb = require('google-protobuf/google/protobuf/wrappers_pb.js')

var google_protobuf_descriptor_pb = require('google-protobuf/google/protobuf/descriptor_pb.js')
const proto = {};
proto.grpc = {};
proto.grpc.service = {};
proto.grpc.service.webservice = require('./service_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.grpc.service.webservice.WebServiceClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'binary';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.grpc.service.webservice.WebServicePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'binary';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.RegisterPipedRequest,
 *   !proto.grpc.service.webservice.RegisterPipedResponse>}
 */
const methodDescriptor_WebService_RegisterPiped = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/RegisterPiped',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.RegisterPipedRequest,
  proto.grpc.service.webservice.RegisterPipedResponse,
  /**
   * @param {!proto.grpc.service.webservice.RegisterPipedRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.RegisterPipedResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.RegisterPipedRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.RegisterPipedResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.RegisterPipedResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.registerPiped =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/RegisterPiped',
      request,
      metadata || {},
      methodDescriptor_WebService_RegisterPiped,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.RegisterPipedRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.RegisterPipedResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.registerPiped =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/RegisterPiped',
      request,
      metadata || {},
      methodDescriptor_WebService_RegisterPiped);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.UpdatePipedRequest,
 *   !proto.grpc.service.webservice.UpdatePipedResponse>}
 */
const methodDescriptor_WebService_UpdatePiped = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/UpdatePiped',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.UpdatePipedRequest,
  proto.grpc.service.webservice.UpdatePipedResponse,
  /**
   * @param {!proto.grpc.service.webservice.UpdatePipedRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.UpdatePipedResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.UpdatePipedRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.UpdatePipedResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.UpdatePipedResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.updatePiped =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/UpdatePiped',
      request,
      metadata || {},
      methodDescriptor_WebService_UpdatePiped,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.UpdatePipedRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.UpdatePipedResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.updatePiped =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/UpdatePiped',
      request,
      metadata || {},
      methodDescriptor_WebService_UpdatePiped);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.RecreatePipedKeyRequest,
 *   !proto.grpc.service.webservice.RecreatePipedKeyResponse>}
 */
const methodDescriptor_WebService_RecreatePipedKey = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/RecreatePipedKey',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.RecreatePipedKeyRequest,
  proto.grpc.service.webservice.RecreatePipedKeyResponse,
  /**
   * @param {!proto.grpc.service.webservice.RecreatePipedKeyRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.RecreatePipedKeyResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.RecreatePipedKeyRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.RecreatePipedKeyResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.RecreatePipedKeyResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.recreatePipedKey =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/RecreatePipedKey',
      request,
      metadata || {},
      methodDescriptor_WebService_RecreatePipedKey,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.RecreatePipedKeyRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.RecreatePipedKeyResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.recreatePipedKey =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/RecreatePipedKey',
      request,
      metadata || {},
      methodDescriptor_WebService_RecreatePipedKey);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.DeleteOldPipedKeysRequest,
 *   !proto.grpc.service.webservice.DeleteOldPipedKeysResponse>}
 */
const methodDescriptor_WebService_DeleteOldPipedKeys = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/DeleteOldPipedKeys',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.DeleteOldPipedKeysRequest,
  proto.grpc.service.webservice.DeleteOldPipedKeysResponse,
  /**
   * @param {!proto.grpc.service.webservice.DeleteOldPipedKeysRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.DeleteOldPipedKeysResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.DeleteOldPipedKeysRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.DeleteOldPipedKeysResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.DeleteOldPipedKeysResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.deleteOldPipedKeys =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/DeleteOldPipedKeys',
      request,
      metadata || {},
      methodDescriptor_WebService_DeleteOldPipedKeys,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.DeleteOldPipedKeysRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.DeleteOldPipedKeysResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.deleteOldPipedKeys =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/DeleteOldPipedKeys',
      request,
      metadata || {},
      methodDescriptor_WebService_DeleteOldPipedKeys);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.EnablePipedRequest,
 *   !proto.grpc.service.webservice.EnablePipedResponse>}
 */
const methodDescriptor_WebService_EnablePiped = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/EnablePiped',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.EnablePipedRequest,
  proto.grpc.service.webservice.EnablePipedResponse,
  /**
   * @param {!proto.grpc.service.webservice.EnablePipedRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.EnablePipedResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.EnablePipedRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.EnablePipedResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.EnablePipedResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.enablePiped =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/EnablePiped',
      request,
      metadata || {},
      methodDescriptor_WebService_EnablePiped,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.EnablePipedRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.EnablePipedResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.enablePiped =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/EnablePiped',
      request,
      metadata || {},
      methodDescriptor_WebService_EnablePiped);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.DisablePipedRequest,
 *   !proto.grpc.service.webservice.DisablePipedResponse>}
 */
const methodDescriptor_WebService_DisablePiped = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/DisablePiped',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.DisablePipedRequest,
  proto.grpc.service.webservice.DisablePipedResponse,
  /**
   * @param {!proto.grpc.service.webservice.DisablePipedRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.DisablePipedResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.DisablePipedRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.DisablePipedResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.DisablePipedResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.disablePiped =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/DisablePiped',
      request,
      metadata || {},
      methodDescriptor_WebService_DisablePiped,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.DisablePipedRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.DisablePipedResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.disablePiped =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/DisablePiped',
      request,
      metadata || {},
      methodDescriptor_WebService_DisablePiped);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.ListPipedsRequest,
 *   !proto.grpc.service.webservice.ListPipedsResponse>}
 */
const methodDescriptor_WebService_ListPipeds = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/ListPipeds',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.ListPipedsRequest,
  proto.grpc.service.webservice.ListPipedsResponse,
  /**
   * @param {!proto.grpc.service.webservice.ListPipedsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.ListPipedsResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.ListPipedsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.ListPipedsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.ListPipedsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.listPipeds =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListPipeds',
      request,
      metadata || {},
      methodDescriptor_WebService_ListPipeds,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.ListPipedsRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.ListPipedsResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.listPipeds =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListPipeds',
      request,
      metadata || {},
      methodDescriptor_WebService_ListPipeds);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.GetPipedRequest,
 *   !proto.grpc.service.webservice.GetPipedResponse>}
 */
const methodDescriptor_WebService_GetPiped = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/GetPiped',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.GetPipedRequest,
  proto.grpc.service.webservice.GetPipedResponse,
  /**
   * @param {!proto.grpc.service.webservice.GetPipedRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.GetPipedResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.GetPipedRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.GetPipedResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.GetPipedResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.getPiped =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetPiped',
      request,
      metadata || {},
      methodDescriptor_WebService_GetPiped,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.GetPipedRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.GetPipedResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.getPiped =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetPiped',
      request,
      metadata || {},
      methodDescriptor_WebService_GetPiped);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest,
 *   !proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse>}
 */
const methodDescriptor_WebService_UpdatePipedDesiredVersion = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/UpdatePipedDesiredVersion',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest,
  proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse,
  /**
   * @param {!proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.updatePipedDesiredVersion =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/UpdatePipedDesiredVersion',
      request,
      metadata || {},
      methodDescriptor_WebService_UpdatePipedDesiredVersion,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.UpdatePipedDesiredVersionRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.UpdatePipedDesiredVersionResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.updatePipedDesiredVersion =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/UpdatePipedDesiredVersion',
      request,
      metadata || {},
      methodDescriptor_WebService_UpdatePipedDesiredVersion);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.RestartPipedRequest,
 *   !proto.grpc.service.webservice.RestartPipedResponse>}
 */
const methodDescriptor_WebService_RestartPiped = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/RestartPiped',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.RestartPipedRequest,
  proto.grpc.service.webservice.RestartPipedResponse,
  /**
   * @param {!proto.grpc.service.webservice.RestartPipedRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.RestartPipedResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.RestartPipedRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.RestartPipedResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.RestartPipedResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.restartPiped =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/RestartPiped',
      request,
      metadata || {},
      methodDescriptor_WebService_RestartPiped,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.RestartPipedRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.RestartPipedResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.restartPiped =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/RestartPiped',
      request,
      metadata || {},
      methodDescriptor_WebService_RestartPiped);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.ListReleasedVersionsRequest,
 *   !proto.grpc.service.webservice.ListReleasedVersionsResponse>}
 */
const methodDescriptor_WebService_ListReleasedVersions = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/ListReleasedVersions',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.ListReleasedVersionsRequest,
  proto.grpc.service.webservice.ListReleasedVersionsResponse,
  /**
   * @param {!proto.grpc.service.webservice.ListReleasedVersionsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.ListReleasedVersionsResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.ListReleasedVersionsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.ListReleasedVersionsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.ListReleasedVersionsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.listReleasedVersions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListReleasedVersions',
      request,
      metadata || {},
      methodDescriptor_WebService_ListReleasedVersions,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.ListReleasedVersionsRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.ListReleasedVersionsResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.listReleasedVersions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListReleasedVersions',
      request,
      metadata || {},
      methodDescriptor_WebService_ListReleasedVersions);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.ListDeprecatedNotesRequest,
 *   !proto.grpc.service.webservice.ListDeprecatedNotesResponse>}
 */
const methodDescriptor_WebService_ListDeprecatedNotes = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/ListDeprecatedNotes',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.ListDeprecatedNotesRequest,
  proto.grpc.service.webservice.ListDeprecatedNotesResponse,
  /**
   * @param {!proto.grpc.service.webservice.ListDeprecatedNotesRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.ListDeprecatedNotesResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.ListDeprecatedNotesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.ListDeprecatedNotesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.ListDeprecatedNotesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.listDeprecatedNotes =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListDeprecatedNotes',
      request,
      metadata || {},
      methodDescriptor_WebService_ListDeprecatedNotes,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.ListDeprecatedNotesRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.ListDeprecatedNotesResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.listDeprecatedNotes =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListDeprecatedNotes',
      request,
      metadata || {},
      methodDescriptor_WebService_ListDeprecatedNotes);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.AddApplicationRequest,
 *   !proto.grpc.service.webservice.AddApplicationResponse>}
 */
const methodDescriptor_WebService_AddApplication = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/AddApplication',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.AddApplicationRequest,
  proto.grpc.service.webservice.AddApplicationResponse,
  /**
   * @param {!proto.grpc.service.webservice.AddApplicationRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.AddApplicationResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.AddApplicationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.AddApplicationResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.AddApplicationResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.addApplication =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/AddApplication',
      request,
      metadata || {},
      methodDescriptor_WebService_AddApplication,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.AddApplicationRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.AddApplicationResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.addApplication =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/AddApplication',
      request,
      metadata || {},
      methodDescriptor_WebService_AddApplication);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.UpdateApplicationRequest,
 *   !proto.grpc.service.webservice.UpdateApplicationResponse>}
 */
const methodDescriptor_WebService_UpdateApplication = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/UpdateApplication',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.UpdateApplicationRequest,
  proto.grpc.service.webservice.UpdateApplicationResponse,
  /**
   * @param {!proto.grpc.service.webservice.UpdateApplicationRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.UpdateApplicationResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.UpdateApplicationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.UpdateApplicationResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.UpdateApplicationResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.updateApplication =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/UpdateApplication',
      request,
      metadata || {},
      methodDescriptor_WebService_UpdateApplication,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.UpdateApplicationRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.UpdateApplicationResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.updateApplication =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/UpdateApplication',
      request,
      metadata || {},
      methodDescriptor_WebService_UpdateApplication);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.EnableApplicationRequest,
 *   !proto.grpc.service.webservice.EnableApplicationResponse>}
 */
const methodDescriptor_WebService_EnableApplication = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/EnableApplication',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.EnableApplicationRequest,
  proto.grpc.service.webservice.EnableApplicationResponse,
  /**
   * @param {!proto.grpc.service.webservice.EnableApplicationRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.EnableApplicationResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.EnableApplicationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.EnableApplicationResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.EnableApplicationResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.enableApplication =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/EnableApplication',
      request,
      metadata || {},
      methodDescriptor_WebService_EnableApplication,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.EnableApplicationRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.EnableApplicationResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.enableApplication =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/EnableApplication',
      request,
      metadata || {},
      methodDescriptor_WebService_EnableApplication);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.DisableApplicationRequest,
 *   !proto.grpc.service.webservice.DisableApplicationResponse>}
 */
const methodDescriptor_WebService_DisableApplication = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/DisableApplication',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.DisableApplicationRequest,
  proto.grpc.service.webservice.DisableApplicationResponse,
  /**
   * @param {!proto.grpc.service.webservice.DisableApplicationRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.DisableApplicationResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.DisableApplicationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.DisableApplicationResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.DisableApplicationResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.disableApplication =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/DisableApplication',
      request,
      metadata || {},
      methodDescriptor_WebService_DisableApplication,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.DisableApplicationRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.DisableApplicationResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.disableApplication =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/DisableApplication',
      request,
      metadata || {},
      methodDescriptor_WebService_DisableApplication);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.DeleteApplicationRequest,
 *   !proto.grpc.service.webservice.DeleteApplicationResponse>}
 */
const methodDescriptor_WebService_DeleteApplication = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/DeleteApplication',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.DeleteApplicationRequest,
  proto.grpc.service.webservice.DeleteApplicationResponse,
  /**
   * @param {!proto.grpc.service.webservice.DeleteApplicationRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.DeleteApplicationResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.DeleteApplicationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.DeleteApplicationResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.DeleteApplicationResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.deleteApplication =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/DeleteApplication',
      request,
      metadata || {},
      methodDescriptor_WebService_DeleteApplication,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.DeleteApplicationRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.DeleteApplicationResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.deleteApplication =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/DeleteApplication',
      request,
      metadata || {},
      methodDescriptor_WebService_DeleteApplication);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.ListApplicationsRequest,
 *   !proto.grpc.service.webservice.ListApplicationsResponse>}
 */
const methodDescriptor_WebService_ListApplications = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/ListApplications',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.ListApplicationsRequest,
  proto.grpc.service.webservice.ListApplicationsResponse,
  /**
   * @param {!proto.grpc.service.webservice.ListApplicationsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.ListApplicationsResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.ListApplicationsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.ListApplicationsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.ListApplicationsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.listApplications =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListApplications',
      request,
      metadata || {},
      methodDescriptor_WebService_ListApplications,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.ListApplicationsRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.ListApplicationsResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.listApplications =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListApplications',
      request,
      metadata || {},
      methodDescriptor_WebService_ListApplications);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.SyncApplicationRequest,
 *   !proto.grpc.service.webservice.SyncApplicationResponse>}
 */
const methodDescriptor_WebService_SyncApplication = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/SyncApplication',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.SyncApplicationRequest,
  proto.grpc.service.webservice.SyncApplicationResponse,
  /**
   * @param {!proto.grpc.service.webservice.SyncApplicationRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.SyncApplicationResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.SyncApplicationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.SyncApplicationResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.SyncApplicationResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.syncApplication =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/SyncApplication',
      request,
      metadata || {},
      methodDescriptor_WebService_SyncApplication,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.SyncApplicationRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.SyncApplicationResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.syncApplication =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/SyncApplication',
      request,
      metadata || {},
      methodDescriptor_WebService_SyncApplication);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.GetApplicationRequest,
 *   !proto.grpc.service.webservice.GetApplicationResponse>}
 */
const methodDescriptor_WebService_GetApplication = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/GetApplication',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.GetApplicationRequest,
  proto.grpc.service.webservice.GetApplicationResponse,
  /**
   * @param {!proto.grpc.service.webservice.GetApplicationRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.GetApplicationResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.GetApplicationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.GetApplicationResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.GetApplicationResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.getApplication =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetApplication',
      request,
      metadata || {},
      methodDescriptor_WebService_GetApplication,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.GetApplicationRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.GetApplicationResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.getApplication =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetApplication',
      request,
      metadata || {},
      methodDescriptor_WebService_GetApplication);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest,
 *   !proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse>}
 */
const methodDescriptor_WebService_GenerateApplicationSealedSecret = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/GenerateApplicationSealedSecret',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest,
  proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse,
  /**
   * @param {!proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.generateApplicationSealedSecret =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GenerateApplicationSealedSecret',
      request,
      metadata || {},
      methodDescriptor_WebService_GenerateApplicationSealedSecret,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.GenerateApplicationSealedSecretRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.GenerateApplicationSealedSecretResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.generateApplicationSealedSecret =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GenerateApplicationSealedSecret',
      request,
      metadata || {},
      methodDescriptor_WebService_GenerateApplicationSealedSecret);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.ListUnregisteredApplicationsRequest,
 *   !proto.grpc.service.webservice.ListUnregisteredApplicationsResponse>}
 */
const methodDescriptor_WebService_ListUnregisteredApplications = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/ListUnregisteredApplications',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.ListUnregisteredApplicationsRequest,
  proto.grpc.service.webservice.ListUnregisteredApplicationsResponse,
  /**
   * @param {!proto.grpc.service.webservice.ListUnregisteredApplicationsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.ListUnregisteredApplicationsResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.ListUnregisteredApplicationsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.ListUnregisteredApplicationsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.ListUnregisteredApplicationsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.listUnregisteredApplications =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListUnregisteredApplications',
      request,
      metadata || {},
      methodDescriptor_WebService_ListUnregisteredApplications,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.ListUnregisteredApplicationsRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.ListUnregisteredApplicationsResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.listUnregisteredApplications =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListUnregisteredApplications',
      request,
      metadata || {},
      methodDescriptor_WebService_ListUnregisteredApplications);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.ListDeploymentsRequest,
 *   !proto.grpc.service.webservice.ListDeploymentsResponse>}
 */
const methodDescriptor_WebService_ListDeployments = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/ListDeployments',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.ListDeploymentsRequest,
  proto.grpc.service.webservice.ListDeploymentsResponse,
  /**
   * @param {!proto.grpc.service.webservice.ListDeploymentsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.ListDeploymentsResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.ListDeploymentsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.ListDeploymentsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.ListDeploymentsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.listDeployments =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListDeployments',
      request,
      metadata || {},
      methodDescriptor_WebService_ListDeployments,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.ListDeploymentsRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.ListDeploymentsResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.listDeployments =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListDeployments',
      request,
      metadata || {},
      methodDescriptor_WebService_ListDeployments);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.GetDeploymentRequest,
 *   !proto.grpc.service.webservice.GetDeploymentResponse>}
 */
const methodDescriptor_WebService_GetDeployment = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/GetDeployment',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.GetDeploymentRequest,
  proto.grpc.service.webservice.GetDeploymentResponse,
  /**
   * @param {!proto.grpc.service.webservice.GetDeploymentRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.GetDeploymentResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.GetDeploymentRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.GetDeploymentResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.GetDeploymentResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.getDeployment =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetDeployment',
      request,
      metadata || {},
      methodDescriptor_WebService_GetDeployment,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.GetDeploymentRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.GetDeploymentResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.getDeployment =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetDeployment',
      request,
      metadata || {},
      methodDescriptor_WebService_GetDeployment);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.GetStageLogRequest,
 *   !proto.grpc.service.webservice.GetStageLogResponse>}
 */
const methodDescriptor_WebService_GetStageLog = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/GetStageLog',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.GetStageLogRequest,
  proto.grpc.service.webservice.GetStageLogResponse,
  /**
   * @param {!proto.grpc.service.webservice.GetStageLogRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.GetStageLogResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.GetStageLogRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.GetStageLogResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.GetStageLogResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.getStageLog =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetStageLog',
      request,
      metadata || {},
      methodDescriptor_WebService_GetStageLog,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.GetStageLogRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.GetStageLogResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.getStageLog =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetStageLog',
      request,
      metadata || {},
      methodDescriptor_WebService_GetStageLog);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.CancelDeploymentRequest,
 *   !proto.grpc.service.webservice.CancelDeploymentResponse>}
 */
const methodDescriptor_WebService_CancelDeployment = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/CancelDeployment',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.CancelDeploymentRequest,
  proto.grpc.service.webservice.CancelDeploymentResponse,
  /**
   * @param {!proto.grpc.service.webservice.CancelDeploymentRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.CancelDeploymentResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.CancelDeploymentRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.CancelDeploymentResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.CancelDeploymentResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.cancelDeployment =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/CancelDeployment',
      request,
      metadata || {},
      methodDescriptor_WebService_CancelDeployment,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.CancelDeploymentRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.CancelDeploymentResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.cancelDeployment =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/CancelDeployment',
      request,
      metadata || {},
      methodDescriptor_WebService_CancelDeployment);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.SkipStageRequest,
 *   !proto.grpc.service.webservice.SkipStageResponse>}
 */
const methodDescriptor_WebService_SkipStage = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/SkipStage',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.SkipStageRequest,
  proto.grpc.service.webservice.SkipStageResponse,
  /**
   * @param {!proto.grpc.service.webservice.SkipStageRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.SkipStageResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.SkipStageRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.SkipStageResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.SkipStageResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.skipStage =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/SkipStage',
      request,
      metadata || {},
      methodDescriptor_WebService_SkipStage,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.SkipStageRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.SkipStageResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.skipStage =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/SkipStage',
      request,
      metadata || {},
      methodDescriptor_WebService_SkipStage);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.ApproveStageRequest,
 *   !proto.grpc.service.webservice.ApproveStageResponse>}
 */
const methodDescriptor_WebService_ApproveStage = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/ApproveStage',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.ApproveStageRequest,
  proto.grpc.service.webservice.ApproveStageResponse,
  /**
   * @param {!proto.grpc.service.webservice.ApproveStageRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.ApproveStageResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.ApproveStageRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.ApproveStageResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.ApproveStageResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.approveStage =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ApproveStage',
      request,
      metadata || {},
      methodDescriptor_WebService_ApproveStage,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.ApproveStageRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.ApproveStageResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.approveStage =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ApproveStage',
      request,
      metadata || {},
      methodDescriptor_WebService_ApproveStage);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.ListDeploymentTracesRequest,
 *   !proto.grpc.service.webservice.ListDeploymentTracesResponse>}
 */
const methodDescriptor_WebService_ListDeploymentTraces = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/ListDeploymentTraces',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.ListDeploymentTracesRequest,
  proto.grpc.service.webservice.ListDeploymentTracesResponse,
  /**
   * @param {!proto.grpc.service.webservice.ListDeploymentTracesRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.ListDeploymentTracesResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.ListDeploymentTracesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.ListDeploymentTracesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.ListDeploymentTracesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.listDeploymentTraces =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListDeploymentTraces',
      request,
      metadata || {},
      methodDescriptor_WebService_ListDeploymentTraces,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.ListDeploymentTracesRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.ListDeploymentTracesResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.listDeploymentTraces =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListDeploymentTraces',
      request,
      metadata || {},
      methodDescriptor_WebService_ListDeploymentTraces);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.GetApplicationLiveStateRequest,
 *   !proto.grpc.service.webservice.GetApplicationLiveStateResponse>}
 */
const methodDescriptor_WebService_GetApplicationLiveState = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/GetApplicationLiveState',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.GetApplicationLiveStateRequest,
  proto.grpc.service.webservice.GetApplicationLiveStateResponse,
  /**
   * @param {!proto.grpc.service.webservice.GetApplicationLiveStateRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.GetApplicationLiveStateResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.GetApplicationLiveStateRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.GetApplicationLiveStateResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.GetApplicationLiveStateResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.getApplicationLiveState =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetApplicationLiveState',
      request,
      metadata || {},
      methodDescriptor_WebService_GetApplicationLiveState,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.GetApplicationLiveStateRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.GetApplicationLiveStateResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.getApplicationLiveState =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetApplicationLiveState',
      request,
      metadata || {},
      methodDescriptor_WebService_GetApplicationLiveState);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.GetProjectRequest,
 *   !proto.grpc.service.webservice.GetProjectResponse>}
 */
const methodDescriptor_WebService_GetProject = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/GetProject',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.GetProjectRequest,
  proto.grpc.service.webservice.GetProjectResponse,
  /**
   * @param {!proto.grpc.service.webservice.GetProjectRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.GetProjectResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.GetProjectRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.GetProjectResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.GetProjectResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.getProject =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetProject',
      request,
      metadata || {},
      methodDescriptor_WebService_GetProject,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.GetProjectRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.GetProjectResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.getProject =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetProject',
      request,
      metadata || {},
      methodDescriptor_WebService_GetProject);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.UpdateProjectStaticAdminRequest,
 *   !proto.grpc.service.webservice.UpdateProjectStaticAdminResponse>}
 */
const methodDescriptor_WebService_UpdateProjectStaticAdmin = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/UpdateProjectStaticAdmin',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.UpdateProjectStaticAdminRequest,
  proto.grpc.service.webservice.UpdateProjectStaticAdminResponse,
  /**
   * @param {!proto.grpc.service.webservice.UpdateProjectStaticAdminRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.UpdateProjectStaticAdminResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.UpdateProjectStaticAdminRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.UpdateProjectStaticAdminResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.UpdateProjectStaticAdminResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.updateProjectStaticAdmin =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/UpdateProjectStaticAdmin',
      request,
      metadata || {},
      methodDescriptor_WebService_UpdateProjectStaticAdmin,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.UpdateProjectStaticAdminRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.UpdateProjectStaticAdminResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.updateProjectStaticAdmin =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/UpdateProjectStaticAdmin',
      request,
      metadata || {},
      methodDescriptor_WebService_UpdateProjectStaticAdmin);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.EnableStaticAdminRequest,
 *   !proto.grpc.service.webservice.EnableStaticAdminResponse>}
 */
const methodDescriptor_WebService_EnableStaticAdmin = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/EnableStaticAdmin',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.EnableStaticAdminRequest,
  proto.grpc.service.webservice.EnableStaticAdminResponse,
  /**
   * @param {!proto.grpc.service.webservice.EnableStaticAdminRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.EnableStaticAdminResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.EnableStaticAdminRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.EnableStaticAdminResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.EnableStaticAdminResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.enableStaticAdmin =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/EnableStaticAdmin',
      request,
      metadata || {},
      methodDescriptor_WebService_EnableStaticAdmin,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.EnableStaticAdminRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.EnableStaticAdminResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.enableStaticAdmin =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/EnableStaticAdmin',
      request,
      metadata || {},
      methodDescriptor_WebService_EnableStaticAdmin);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.DisableStaticAdminRequest,
 *   !proto.grpc.service.webservice.DisableStaticAdminResponse>}
 */
const methodDescriptor_WebService_DisableStaticAdmin = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/DisableStaticAdmin',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.DisableStaticAdminRequest,
  proto.grpc.service.webservice.DisableStaticAdminResponse,
  /**
   * @param {!proto.grpc.service.webservice.DisableStaticAdminRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.DisableStaticAdminResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.DisableStaticAdminRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.DisableStaticAdminResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.DisableStaticAdminResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.disableStaticAdmin =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/DisableStaticAdmin',
      request,
      metadata || {},
      methodDescriptor_WebService_DisableStaticAdmin,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.DisableStaticAdminRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.DisableStaticAdminResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.disableStaticAdmin =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/DisableStaticAdmin',
      request,
      metadata || {},
      methodDescriptor_WebService_DisableStaticAdmin);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.UpdateProjectSSOConfigRequest,
 *   !proto.grpc.service.webservice.UpdateProjectSSOConfigResponse>}
 */
const methodDescriptor_WebService_UpdateProjectSSOConfig = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/UpdateProjectSSOConfig',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.UpdateProjectSSOConfigRequest,
  proto.grpc.service.webservice.UpdateProjectSSOConfigResponse,
  /**
   * @param {!proto.grpc.service.webservice.UpdateProjectSSOConfigRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.UpdateProjectSSOConfigResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.UpdateProjectSSOConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.UpdateProjectSSOConfigResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.UpdateProjectSSOConfigResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.updateProjectSSOConfig =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/UpdateProjectSSOConfig',
      request,
      metadata || {},
      methodDescriptor_WebService_UpdateProjectSSOConfig,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.UpdateProjectSSOConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.UpdateProjectSSOConfigResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.updateProjectSSOConfig =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/UpdateProjectSSOConfig',
      request,
      metadata || {},
      methodDescriptor_WebService_UpdateProjectSSOConfig);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.UpdateProjectRBACConfigRequest,
 *   !proto.grpc.service.webservice.UpdateProjectRBACConfigResponse>}
 */
const methodDescriptor_WebService_UpdateProjectRBACConfig = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/UpdateProjectRBACConfig',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.UpdateProjectRBACConfigRequest,
  proto.grpc.service.webservice.UpdateProjectRBACConfigResponse,
  /**
   * @param {!proto.grpc.service.webservice.UpdateProjectRBACConfigRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.UpdateProjectRBACConfigResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.UpdateProjectRBACConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.UpdateProjectRBACConfigResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.UpdateProjectRBACConfigResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.updateProjectRBACConfig =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/UpdateProjectRBACConfig',
      request,
      metadata || {},
      methodDescriptor_WebService_UpdateProjectRBACConfig,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.UpdateProjectRBACConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.UpdateProjectRBACConfigResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.updateProjectRBACConfig =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/UpdateProjectRBACConfig',
      request,
      metadata || {},
      methodDescriptor_WebService_UpdateProjectRBACConfig);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.GetMeRequest,
 *   !proto.grpc.service.webservice.GetMeResponse>}
 */
const methodDescriptor_WebService_GetMe = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/GetMe',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.GetMeRequest,
  proto.grpc.service.webservice.GetMeResponse,
  /**
   * @param {!proto.grpc.service.webservice.GetMeRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.GetMeResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.GetMeRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.GetMeResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.GetMeResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.getMe =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetMe',
      request,
      metadata || {},
      methodDescriptor_WebService_GetMe,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.GetMeRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.GetMeResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.getMe =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetMe',
      request,
      metadata || {},
      methodDescriptor_WebService_GetMe);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.AddProjectRBACRoleRequest,
 *   !proto.grpc.service.webservice.AddProjectRBACRoleResponse>}
 */
const methodDescriptor_WebService_AddProjectRBACRole = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/AddProjectRBACRole',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.AddProjectRBACRoleRequest,
  proto.grpc.service.webservice.AddProjectRBACRoleResponse,
  /**
   * @param {!proto.grpc.service.webservice.AddProjectRBACRoleRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.AddProjectRBACRoleResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.AddProjectRBACRoleRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.AddProjectRBACRoleResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.AddProjectRBACRoleResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.addProjectRBACRole =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/AddProjectRBACRole',
      request,
      metadata || {},
      methodDescriptor_WebService_AddProjectRBACRole,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.AddProjectRBACRoleRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.AddProjectRBACRoleResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.addProjectRBACRole =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/AddProjectRBACRole',
      request,
      metadata || {},
      methodDescriptor_WebService_AddProjectRBACRole);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.UpdateProjectRBACRoleRequest,
 *   !proto.grpc.service.webservice.UpdateProjectRBACRoleResponse>}
 */
const methodDescriptor_WebService_UpdateProjectRBACRole = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/UpdateProjectRBACRole',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.UpdateProjectRBACRoleRequest,
  proto.grpc.service.webservice.UpdateProjectRBACRoleResponse,
  /**
   * @param {!proto.grpc.service.webservice.UpdateProjectRBACRoleRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.UpdateProjectRBACRoleResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.UpdateProjectRBACRoleRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.UpdateProjectRBACRoleResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.UpdateProjectRBACRoleResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.updateProjectRBACRole =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/UpdateProjectRBACRole',
      request,
      metadata || {},
      methodDescriptor_WebService_UpdateProjectRBACRole,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.UpdateProjectRBACRoleRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.UpdateProjectRBACRoleResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.updateProjectRBACRole =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/UpdateProjectRBACRole',
      request,
      metadata || {},
      methodDescriptor_WebService_UpdateProjectRBACRole);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.DeleteProjectRBACRoleRequest,
 *   !proto.grpc.service.webservice.DeleteProjectRBACRoleResponse>}
 */
const methodDescriptor_WebService_DeleteProjectRBACRole = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/DeleteProjectRBACRole',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.DeleteProjectRBACRoleRequest,
  proto.grpc.service.webservice.DeleteProjectRBACRoleResponse,
  /**
   * @param {!proto.grpc.service.webservice.DeleteProjectRBACRoleRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.DeleteProjectRBACRoleResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.DeleteProjectRBACRoleRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.DeleteProjectRBACRoleResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.DeleteProjectRBACRoleResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.deleteProjectRBACRole =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/DeleteProjectRBACRole',
      request,
      metadata || {},
      methodDescriptor_WebService_DeleteProjectRBACRole,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.DeleteProjectRBACRoleRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.DeleteProjectRBACRoleResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.deleteProjectRBACRole =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/DeleteProjectRBACRole',
      request,
      metadata || {},
      methodDescriptor_WebService_DeleteProjectRBACRole);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.AddProjectUserGroupRequest,
 *   !proto.grpc.service.webservice.AddProjectUserGroupResponse>}
 */
const methodDescriptor_WebService_AddProjectUserGroup = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/AddProjectUserGroup',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.AddProjectUserGroupRequest,
  proto.grpc.service.webservice.AddProjectUserGroupResponse,
  /**
   * @param {!proto.grpc.service.webservice.AddProjectUserGroupRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.AddProjectUserGroupResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.AddProjectUserGroupRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.AddProjectUserGroupResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.AddProjectUserGroupResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.addProjectUserGroup =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/AddProjectUserGroup',
      request,
      metadata || {},
      methodDescriptor_WebService_AddProjectUserGroup,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.AddProjectUserGroupRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.AddProjectUserGroupResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.addProjectUserGroup =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/AddProjectUserGroup',
      request,
      metadata || {},
      methodDescriptor_WebService_AddProjectUserGroup);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.DeleteProjectUserGroupRequest,
 *   !proto.grpc.service.webservice.DeleteProjectUserGroupResponse>}
 */
const methodDescriptor_WebService_DeleteProjectUserGroup = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/DeleteProjectUserGroup',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.DeleteProjectUserGroupRequest,
  proto.grpc.service.webservice.DeleteProjectUserGroupResponse,
  /**
   * @param {!proto.grpc.service.webservice.DeleteProjectUserGroupRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.DeleteProjectUserGroupResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.DeleteProjectUserGroupRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.DeleteProjectUserGroupResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.DeleteProjectUserGroupResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.deleteProjectUserGroup =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/DeleteProjectUserGroup',
      request,
      metadata || {},
      methodDescriptor_WebService_DeleteProjectUserGroup,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.DeleteProjectUserGroupRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.DeleteProjectUserGroupResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.deleteProjectUserGroup =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/DeleteProjectUserGroup',
      request,
      metadata || {},
      methodDescriptor_WebService_DeleteProjectUserGroup);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.GetCommandRequest,
 *   !proto.grpc.service.webservice.GetCommandResponse>}
 */
const methodDescriptor_WebService_GetCommand = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/GetCommand',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.GetCommandRequest,
  proto.grpc.service.webservice.GetCommandResponse,
  /**
   * @param {!proto.grpc.service.webservice.GetCommandRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.GetCommandResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.GetCommandRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.GetCommandResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.GetCommandResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.getCommand =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetCommand',
      request,
      metadata || {},
      methodDescriptor_WebService_GetCommand,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.GetCommandRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.GetCommandResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.getCommand =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetCommand',
      request,
      metadata || {},
      methodDescriptor_WebService_GetCommand);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.GenerateAPIKeyRequest,
 *   !proto.grpc.service.webservice.GenerateAPIKeyResponse>}
 */
const methodDescriptor_WebService_GenerateAPIKey = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/GenerateAPIKey',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.GenerateAPIKeyRequest,
  proto.grpc.service.webservice.GenerateAPIKeyResponse,
  /**
   * @param {!proto.grpc.service.webservice.GenerateAPIKeyRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.GenerateAPIKeyResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.GenerateAPIKeyRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.GenerateAPIKeyResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.GenerateAPIKeyResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.generateAPIKey =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GenerateAPIKey',
      request,
      metadata || {},
      methodDescriptor_WebService_GenerateAPIKey,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.GenerateAPIKeyRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.GenerateAPIKeyResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.generateAPIKey =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GenerateAPIKey',
      request,
      metadata || {},
      methodDescriptor_WebService_GenerateAPIKey);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.DisableAPIKeyRequest,
 *   !proto.grpc.service.webservice.DisableAPIKeyResponse>}
 */
const methodDescriptor_WebService_DisableAPIKey = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/DisableAPIKey',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.DisableAPIKeyRequest,
  proto.grpc.service.webservice.DisableAPIKeyResponse,
  /**
   * @param {!proto.grpc.service.webservice.DisableAPIKeyRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.DisableAPIKeyResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.DisableAPIKeyRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.DisableAPIKeyResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.DisableAPIKeyResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.disableAPIKey =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/DisableAPIKey',
      request,
      metadata || {},
      methodDescriptor_WebService_DisableAPIKey,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.DisableAPIKeyRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.DisableAPIKeyResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.disableAPIKey =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/DisableAPIKey',
      request,
      metadata || {},
      methodDescriptor_WebService_DisableAPIKey);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.ListAPIKeysRequest,
 *   !proto.grpc.service.webservice.ListAPIKeysResponse>}
 */
const methodDescriptor_WebService_ListAPIKeys = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/ListAPIKeys',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.ListAPIKeysRequest,
  proto.grpc.service.webservice.ListAPIKeysResponse,
  /**
   * @param {!proto.grpc.service.webservice.ListAPIKeysRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.ListAPIKeysResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.ListAPIKeysRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.ListAPIKeysResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.ListAPIKeysResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.listAPIKeys =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListAPIKeys',
      request,
      metadata || {},
      methodDescriptor_WebService_ListAPIKeys,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.ListAPIKeysRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.ListAPIKeysResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.listAPIKeys =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListAPIKeys',
      request,
      metadata || {},
      methodDescriptor_WebService_ListAPIKeys);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.GetInsightDataRequest,
 *   !proto.grpc.service.webservice.GetInsightDataResponse>}
 */
const methodDescriptor_WebService_GetInsightData = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/GetInsightData',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.GetInsightDataRequest,
  proto.grpc.service.webservice.GetInsightDataResponse,
  /**
   * @param {!proto.grpc.service.webservice.GetInsightDataRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.GetInsightDataResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.GetInsightDataRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.GetInsightDataResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.GetInsightDataResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.getInsightData =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetInsightData',
      request,
      metadata || {},
      methodDescriptor_WebService_GetInsightData,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.GetInsightDataRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.GetInsightDataResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.getInsightData =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetInsightData',
      request,
      metadata || {},
      methodDescriptor_WebService_GetInsightData);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.GetInsightApplicationCountRequest,
 *   !proto.grpc.service.webservice.GetInsightApplicationCountResponse>}
 */
const methodDescriptor_WebService_GetInsightApplicationCount = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/GetInsightApplicationCount',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.GetInsightApplicationCountRequest,
  proto.grpc.service.webservice.GetInsightApplicationCountResponse,
  /**
   * @param {!proto.grpc.service.webservice.GetInsightApplicationCountRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.GetInsightApplicationCountResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.GetInsightApplicationCountRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.GetInsightApplicationCountResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.GetInsightApplicationCountResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.getInsightApplicationCount =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetInsightApplicationCount',
      request,
      metadata || {},
      methodDescriptor_WebService_GetInsightApplicationCount,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.GetInsightApplicationCountRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.GetInsightApplicationCountResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.getInsightApplicationCount =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetInsightApplicationCount',
      request,
      metadata || {},
      methodDescriptor_WebService_GetInsightApplicationCount);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.ListDeploymentChainsRequest,
 *   !proto.grpc.service.webservice.ListDeploymentChainsResponse>}
 */
const methodDescriptor_WebService_ListDeploymentChains = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/ListDeploymentChains',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.ListDeploymentChainsRequest,
  proto.grpc.service.webservice.ListDeploymentChainsResponse,
  /**
   * @param {!proto.grpc.service.webservice.ListDeploymentChainsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.ListDeploymentChainsResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.ListDeploymentChainsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.ListDeploymentChainsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.ListDeploymentChainsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.listDeploymentChains =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListDeploymentChains',
      request,
      metadata || {},
      methodDescriptor_WebService_ListDeploymentChains,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.ListDeploymentChainsRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.ListDeploymentChainsResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.listDeploymentChains =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListDeploymentChains',
      request,
      metadata || {},
      methodDescriptor_WebService_ListDeploymentChains);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.GetDeploymentChainRequest,
 *   !proto.grpc.service.webservice.GetDeploymentChainResponse>}
 */
const methodDescriptor_WebService_GetDeploymentChain = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/GetDeploymentChain',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.GetDeploymentChainRequest,
  proto.grpc.service.webservice.GetDeploymentChainResponse,
  /**
   * @param {!proto.grpc.service.webservice.GetDeploymentChainRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.GetDeploymentChainResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.GetDeploymentChainRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.GetDeploymentChainResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.GetDeploymentChainResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.getDeploymentChain =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetDeploymentChain',
      request,
      metadata || {},
      methodDescriptor_WebService_GetDeploymentChain,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.GetDeploymentChainRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.GetDeploymentChainResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.getDeploymentChain =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/GetDeploymentChain',
      request,
      metadata || {},
      methodDescriptor_WebService_GetDeploymentChain);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.grpc.service.webservice.ListEventsRequest,
 *   !proto.grpc.service.webservice.ListEventsResponse>}
 */
const methodDescriptor_WebService_ListEvents = new grpc.web.MethodDescriptor(
  '/grpc.service.webservice.WebService/ListEvents',
  grpc.web.MethodType.UNARY,
  proto.grpc.service.webservice.ListEventsRequest,
  proto.grpc.service.webservice.ListEventsResponse,
  /**
   * @param {!proto.grpc.service.webservice.ListEventsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.grpc.service.webservice.ListEventsResponse.deserializeBinary
);


/**
 * @param {!proto.grpc.service.webservice.ListEventsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.grpc.service.webservice.ListEventsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.grpc.service.webservice.ListEventsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.grpc.service.webservice.WebServiceClient.prototype.listEvents =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListEvents',
      request,
      metadata || {},
      methodDescriptor_WebService_ListEvents,
      callback);
};


/**
 * @param {!proto.grpc.service.webservice.ListEventsRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.grpc.service.webservice.ListEventsResponse>}
 *     Promise that resolves to the response
 */
proto.grpc.service.webservice.WebServicePromiseClient.prototype.listEvents =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/grpc.service.webservice.WebService/ListEvents',
      request,
      metadata || {},
      methodDescriptor_WebService_ListEvents);
};


module.exports = proto.grpc.service.webservice;

