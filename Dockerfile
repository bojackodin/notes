FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk update
RUN apk add make

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/ go mod download

COPY . ./
RUN --mount=type=cache,target=/go/pkg/mod/ make build 

FROM alpine:latest
COPY --from=builder /app/bin/app /bin/
COPY --from=builder /app/deployment/etc/config.yml /etc/app/config.yml
ENTRYPOINT [ "/bin/app", "-config", "/etc/app/config.yml" ]
