FROM golang:1.26 AS builder

WORKDIR /go/src/antsankov/go-live

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o /go/bin/go-live .

FROM alpine:3.21 AS production

COPY --from=builder /go/bin/go-live /usr/bin/go-live

LABEL version="1.3.0"

WORKDIR /workdir

ENTRYPOINT ["go-live"]
