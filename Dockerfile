FROM golang:1.11.6-alpine3.9 as builder

LABEL maintainer="Oleg Ozimok oleg.ozimok@corp.kismia.com"

COPY . /go/src/github.com/kismia/pushd

WORKDIR /go/src/github.com/kismia/pushd

RUN go build -o /pushd ./cmd/pushd

FROM alpine:3.9

COPY --from=builder /pushd /usr/bin/pushd

EXPOSE 6379 9100

STOPSIGNAL SIGTERM

ENTRYPOINT ["/usr/bin/pushd"]