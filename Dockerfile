FROM golang:1.22 as builder

ADD . /app/
WORKDIR /app

RUN make build-alpine
RUN chmod +x ./bin/binary

FROM alpine:3.13.5

COPY --from=builder /app/bin/binary ./app/main
COPY --from=builder /app/config/config.yml ./app/config/config.yml
WORKDIR /app
EXPOSE 8080

CMD ["/app/main"]