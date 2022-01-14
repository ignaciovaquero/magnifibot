ARG ARCH
FROM ${ARCH}/golang:1.17.6-alpine3.15 as builder

WORKDIR /go/src/github.com/igvaquero18/magnifibot
COPY . .

RUN go mod tidy && \
  GOOS=linux GOARCH=${ARCH} go build -o magnifibot

FROM ${ARCH}/alpine:3.15.0
COPY --from=builder /go/src/github.com/igvaquero18/magnifibot/magnifibot /magnifibot
EXPOSE 80
ENTRYPOINT ["/magnifibot"]
