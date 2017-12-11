FROM ubuntu 

EXPOSE 50051

ENV LD_LIBRARY_PATH /brotli

ADD ./lib /brotli
ADD ./tiny /

CMD ["/tiny"]

