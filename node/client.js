const grpc = require('grpc');
const fs = require('fs');

const compress = require('./compress/compress_pb');
const services = require('./compress/compress_grpc_pb');

const buf = fs.readFileSync('../assets/lodash.min.js');

const client = new services.CompressClient('127.0.0.1:50051', grpc.credentials.createInsecure());


const request = new compress.CompressRequest();
request.setType(compress.Type.BROTLI);
request.setData(new Uint8Array(buf));
client.do(request, (err, res) => {
  console.dir(err);
  console.dir(res.getData().length);
});