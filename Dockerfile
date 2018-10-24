FROM alpine

MAINTAINER Dmitry Mozzherin

ENV LAST_FULL_REBUILD 2018-10-23

WORKDIR /bin

COPY ./gnfinder/gnfinder /bin

CMD ["gnfinder", "server"]


