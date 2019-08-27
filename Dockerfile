FROM golang:1.12 as builder

ADD . /tiny

RUN apt-get update \
  && apt-get install -y git cmake libpng-dev autoconf automake libtool nasm make \
  && git clone --depth=1 https://github.com/kornelski/pngquant.git /pngquant \
  && cd /pngquant \
  && make && make install \
  && git clone --depth=1 https://github.com/mozilla/mozjpeg.git /mozjpeg \
  && cd /mozjpeg \
  && mkdir build \
  && cd build \
  && cmake -G"Unix Makefiles" ../ \
  && make install \
  && cd /tiny \
  && make test \
  && make build-linux

FROM ubuntu

EXPOSE 7001
EXPOSE 7002

COPY --from=builder /usr/local/bin/pngquant /bin/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libpng16.so.16 /usr/local/lib/
COPY --from=builder /mozjpeg/build/cjpeg /bin/ 
COPY --from=builder /mozjpeg/build/libjpeg.so.62 /usr/local/lib/

COPY --from=builder /tiny/tiny-linux /usr/local/bin/tiny

ENV LD_LIBRARY_PATH /usr/local/lib

RUN apt-get update \
  && apt-get install -y ca-certificates

CMD [ "tiny" ]