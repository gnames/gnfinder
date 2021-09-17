FROM alpine:3.14

MAINTAINER Dmitry Mozzherin

ENV LAST_FULL_REBUILD 2019-02-18

WORKDIR /bin

COPY ./gnfinder/gnfinder /bin

ENTRYPOINT [ "gnfinder" ]

CMD ["-p", "8999"]
