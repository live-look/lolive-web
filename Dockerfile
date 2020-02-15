# Build
FROM golang:1.12 AS compile

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server

# Production
FROM alpine:latest

RUN set -x; \
  apk add ca-certificates

WORKDIR /app
COPY --from=compile /app/web/template /app/templates
# TODO: ARG APP_ENV
COPY --from=compile /app/configs/.env.development /app/
COPY --from=compile /app/server /app/

EXPOSE 3001

CMD ["/app/server"]

