FROM rust:1.60.0 as rustbuilder

ARG TARGETARCH

ENV CAVIF_VERSION=1.3.4

RUN apt-get update -y \
  && apt-get install -y nasm \
  && wget https://github.com/kornelski/cavif-rs/archive/refs/tags/v${CAVIF_VERSION}.tar.gz \
  && tar -xzvf v${CAVIF_VERSION}.tar.gz \
  && mv cavif-rs-${CAVIF_VERSION} cavif-rs \
  && cd cavif-rs \
  && cargo build --release

# FROM golang:1.18 as builder

# ARG TARGETARCH
# ADD . /tiny

# ENV PNGQUANT_VERSION=2.17.0
# ENV MOZJPEG_VERSION=4.0.3

# RUN apt-get update \
#   && apt-get install -y git cmake libpng-dev autoconf automake libtool nasm make wget \
#   && git clone -b "$PNGQUANT_VERSION" --depth=1 https://github.com/kornelski/pngquant.git /pngquant \
#   && cd /pngquant \
#   && ./configure --extra-ldflags=-static --disable-sse && make install \
#   && git clone -b "v$MOZJPEG_VERSION" --depth=1 https://github.com/mozilla/mozjpeg.git /mozjpeg \
#   && cd /mozjpeg \
#   && mkdir build \
#   && cd build \
#   && cmake -G"Unix Makefiles" ../ \
#   && make install \
#   && cp /mozjpeg/build/cjpeg /bin/ \
#   && cd /tiny \
#   && make test \
#   && make build

# FROM ubuntu

# EXPOSE 7001
# EXPOSE 7002

# COPY --from=builder /usr/local/bin/pngquant /usr/local/bin/
# COPY --from=builder /mozjpeg/build/cjpeg-static /usr/local/bin/cjpeg

# COPY --from=builder /tiny/tiny-server /usr/local/bin/tiny-server

# COPY --from=rustbuilder /cavif-rs/target/release/cavif /usr/local/bin/cavif

# RUN apt-get update \
#   && apt-get install -y ca-certificates netcat \
#   && apt-get clean \
#   && rm -rf /var/lib/apt/lists/*

# HEALTHCHECK --interval=10s --timeout=3s \
#   CMD nc -w 1 127.0.0.1 7002

# CMD [ "tiny-server" ]