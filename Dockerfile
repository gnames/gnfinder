FROM alpine

MAINTAINER Dmitry Mozzherin

ENV LAST_FULL_REBUILD 2019-02-18

WORKDIR /bin

COPY ./gnfinder/gnfinder /bin

CMD ["gnfinder", "grpc"]
