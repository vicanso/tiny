// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var compress_compress_pb = require('../compress/compress_pb.js');

function serialize_compress_CompressReply(arg) {
  if (!(arg instanceof compress_compress_pb.CompressReply)) {
    throw new Error('Expected argument of type compress.CompressReply');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_compress_CompressReply(buffer_arg) {
  return compress_compress_pb.CompressReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_compress_CompressRequest(arg) {
  if (!(arg instanceof compress_compress_pb.CompressRequest)) {
    throw new Error('Expected argument of type compress.CompressRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_compress_CompressRequest(buffer_arg) {
  return compress_compress_pb.CompressRequest.deserializeBinary(new Uint8Array(buffer_arg));
}


// The compress service definition.
var CompressService = exports.CompressService = {
  // dom compress 
  do: {
    path: '/compress.Compress/Do',
    requestStream: false,
    responseStream: false,
    requestType: compress_compress_pb.CompressRequest,
    responseType: compress_compress_pb.CompressReply,
    requestSerialize: serialize_compress_CompressRequest,
    requestDeserialize: deserialize_compress_CompressRequest,
    responseSerialize: serialize_compress_CompressReply,
    responseDeserialize: deserialize_compress_CompressReply,
  },
};

exports.CompressClient = grpc.makeGenericClientConstructor(CompressService);
