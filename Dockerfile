# syntax=docker/dockerfile:1
# vim:ft=dockerfile
ARG VER=v1.0.2
FROM golang:1.18.5
ARG VER
WORKDIR /go/src/github.com/mbevc1/helloapp/
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=${VER}" .

FROM scratch
COPY --from=0 /go/src/github.com/mbevc1/helloapp/helloapp .
ENTRYPOINT ["/helloapp"]
