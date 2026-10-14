FROM alpine:3.17

MAINTAINER Dmitry Mozzherin

ENV LAST_FULL_REBUILD 2023-02-22

# send logs to `docker logs` instead of a file inside the container
ENV GNF_LOG_DESTINATION=stderr

WORKDIR /bin

COPY ./gnfinder /bin

ENTRYPOINT [ "gnfinder" ]

CMD ["-p", "8999"]
