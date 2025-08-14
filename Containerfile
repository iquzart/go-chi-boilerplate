FROM golang:1.23-alpine3.20 AS build-env

RUN apk add --no-cache git ca-certificates

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download && go mod verify

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3

COPY . .

RUN swag init -g cmd/api/main.go -o docs && \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-s -w -extldflags '-static'" \
  -a -installsuffix cgo \
  -o app ./cmd/api

FROM scratch

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build-env /src/app /app
COPY --from=build-env /src/docs /docs

USER 65534:65534

EXPOSE 8080

ENTRYPOINT ["/app"]
