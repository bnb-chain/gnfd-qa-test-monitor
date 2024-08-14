# FROM alpine
FROM gcr.io/distroless/static-debian11

WORKDIR /build

COPY monitor .

CMD ["./monitor"]