apt-get update \
  && apt-get install cmake -y \
  && cd /tmp/ \
  && wget https://github.com/google/brotli/archive/v1.0.2.tar.gz -O brotli.tar.gz \
  && tar -xzvf brotli.tar.gz \
  && mv brotli-1.0.2 brotli \
  && cd brotli && ./configure-cmake \
  && make \
  && make test \
  && make install \
  && rm -rf /tiny/lib/* \
  && cp -r /usr/local/lib/libbrotli* /tiny/lib \
  && cd /tiny \
  && go get github.com/buger/jsonparser \
  && go get golang.org/x/net/context \
  && go get google.golang.org/grpc \
  && go get google.golang.org/grpc/reflection \
  && go get github.com/google/brotli/go/cbrotli \
  && go get github.com/nfnt/resize \
  && go get github.com/chai2010/webp \
  && go build
  