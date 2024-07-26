FROM alpine

WORKDIR /build

COPY monitor .

CMD ["./monitor"]