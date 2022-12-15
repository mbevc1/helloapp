# syntax=docker/dockerfile:1
# vim:ft=dockerfile
ARG VER=dev
FROM golang:1.19.4
ARG VER
WORKDIR /go/src/github.com/mbevc1/helloapp/
ENV CGO_ENABLED=0
COPY . .
RUN go mod tidy && GOOS=linux go build -ldflags "-X main.Version=${VER}" .

FROM scratch
COPY --from=0 /go/src/github.com/mbevc1/helloapp/helloapp .
ENTRYPOINT ["/helloapp"]
