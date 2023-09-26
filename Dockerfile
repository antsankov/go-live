FROM golang:1.21 as builder

WORKDIR /go/src/antsankov/go-live

# cache devtools
# COPY ./scripts/devtools.sh /go/src/mikefarah/yq/scripts/devtools.sh
# RUN ./scripts/devtools.sh

COPY . /go/src/antsankov/go-live

RUN CGO_ENABLED=0 make build

# Choose alpine as a base image to make this useful for CI, as many
# CI tools expect an interactive shell inside the container
FROM alpine:3.18 as production

COPY --from=builder /go/src/antsankov/go-live/bin/go-live /usr/bin/go-live
RUN chmod +x /usr/bin/go-live

ARG VERSION=none
LABEL version=1.2.0

WORKDIR /workdir

