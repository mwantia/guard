FROM golang:alpine AS backend

RUN apk add git make

RUN git clone https://github.com/coredns/coredns.git /go/coredns

COPY ./plugin.cfg /go/coredns/plugin.cfg

RUN cd /go/coredns && go get github.com/coredns/records
RUN cd /go/coredns && go get github.com/miekg/redis
RUN cd /go/coredns && go get github.com/mwantia/guard

RUN cd /go/coredns && go generate
RUN cd /go/coredns && make

FROM alpine:latest AS frontend
COPY --from=backend /go/coredns/coredns /usr/bin/coredns

EXPOSE 53 53/udp
VOLUME ["/coredns"]

ENTRYPOINT ["/usr/bin/coredns"]
CMD ["-conf", "/coredns/Corefile"]