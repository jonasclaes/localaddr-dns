FROM golang:1.17
ARG app_name

ENV PORT 8090
ENV TTL 84000
ENV BASE_DOMAIN localaddr.net

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD [ "sh", "-c", "localaddr-dns -port ${PORT} -ttl ${TTL} -base_domain ${BASE_DOMAIN}" ]