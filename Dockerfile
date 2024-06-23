FROM golang:1.19.0-alpine3.16 AS builder

WORKDIR /usr/src/app


RUN apk update \
  && apk add git curl make build-base bash \
  && apk add --no-cache libc6-compat

COPY . .

RUN make build

FROM alpine:3.11

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /usr/local/bin

COPY --from=builder /usr/src/app/bin/autentia-api .
COPY --from=builder /usr/src/app/pkg/mail/templates ./pkg/mail/templates

CMD ["./ghh-api"]
