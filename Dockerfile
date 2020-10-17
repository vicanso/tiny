FROM golang:1.15 as builder

ADD . /tiny

RUN apt-get update \
  && apt-get install -y git cmake libpng-dev autoconf automake libtool nasm make \
  && git clone -b 2.12.6 --depth=1 https://github.com/kornelski/pngquant.git /pngquant \
  && cd /pngquant \
  && make && make install \
  && git clone -b v4.0.1-rc2 --depth=1 https://github.com/mozilla/mozjpeg.git /mozjpeg \
  && cd /mozjpeg \
  && mkdir build \
  && cd build \
  && cmake -G"Unix Makefiles" ../ \
  && make install \
  && cp /mozjpeg/build/cjpeg /bin/ \
  && cd /tiny \
  && make test \
  && make build

FROM ubuntu

EXPOSE 7001
EXPOSE 7002

COPY --from=builder /usr/local/bin/pngquant /usr/local/bin/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libpng16.so.16 /usr/local/lib/
COPY --from=builder /mozjpeg/build/cjpeg /usr/local/bin/ 
COPY --from=builder /mozjpeg/build/libjpeg.so.62 /usr/local/lib/

COPY --from=builder /tiny/tiny-server /usr/local/bin/tiny-server

ENV LD_LIBRARY_PATH /usr/local/lib

RUN apt-get update \
  && apt-get install -y ca-certificates netcat \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

HEALTHCHECK --interval=10s --timeout=3s \
  CMD nc -w 1 127.0.0.1 7002

CMD [ "tiny-server" ]