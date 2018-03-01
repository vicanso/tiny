FROM golang as builder

ADD ./ /go/src/github.com/vicanso/tiny

RUN apt-get update \
  && apt-get install cmake -y \
  && go get -u github.com/golang/dep/cmd/dep \
  && cd /go/src/github.com/vicanso/tiny \
  && dep ensure \
  && cd vendor/github.com/google/brotli/ \
  && ./configure-cmake \
  && make \
  && make test \
  && make install \
  && cd /go/src/github.com/vicanso/tiny \
  GOOS=linux go build -o tiny main.go


FROM ubuntu

EXPOSE 3015
EXPOSE 3016

COPY --from=builder /go/src/github.com/vicanso/tiny/tiny /
COPY --from=builder /usr/local/lib/libbrotlicommon.* /usr/local/lib/
COPY --from=builder /usr/local/lib/libbrotlienc.* /usr/local/lib/
COPY --from=builder /usr/local/lib/libbrotlidec.* /usr/local/lib/

ENV LD_LIBRARY_PATH /usr/local/lib

RUN apt-get update \
  && apt-get install -y ca-certificates 

HEALTHCHECK --interval=30s --timeout=3s \
  CMD /tiny check || exit 1

CMD ["/tiny"]
