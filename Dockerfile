FROM golang:1.20 as builder

RUN apt-get update && apt-get -uy upgrade
RUN apt-get -y install ca-certificates && update-ca-certificates

RUN mkdir /build

WORKDIR /build
COPY . .
RUN go get -v .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o mailer

FROM scratch

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /build/mailer .

ENTRYPOINT ["./mailer"]
