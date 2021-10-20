ARG VER=v1.0.1
FROM golang:1.16.8
ARG VER
WORKDIR /go/src/github.com/mbevc1/helloapp/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=${VER}" .

FROM scratch
COPY --from=0 /go/src/github.com/mbevc1/helloapp/helloapp .
ENTRYPOINT ["/helloapp"]
