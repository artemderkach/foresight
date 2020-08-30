FROM golang:1.15.0-alpine3.12 AS builder

WORKDIR /srv/foresight

COPY . .

# RUN CGO_ENABLED=0 GO111MODULE=on go test -mod=vendor ./...
RUN GO111MODULE=on go build -mod=vendor

FROM alpine:3.12

COPY --from=builder /srv/foresight/foresight /srv/foresight

ENTRYPOINT ["/srv/foresight"]
