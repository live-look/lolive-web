FROM golang:1.12 AS compile

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/camforchat

FROM alpine:latest

RUN apk add tini

ENTRYPOINT ["tini", "--"]

WORKDIR /app
COPY --from=compile /app/web /app/web
COPY --from=compile /app/server /app/

EXPOSE 3001

CMD ["/app/server"]

