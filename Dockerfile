FROM golang:1.17 as builder

RUN mkdir /build

WORKDIR /build
COPY . .
RUN set -x && \
    go get -d -v . && \
    CGO_ENABLED=0 GOOS=linux go build -a -o mailer

FROM scratch
COPY --from=builder /build/mailer .

ENTRYPOINT ["./mailer"]
