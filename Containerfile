FROM golang:1.23-alpine3.20 AS build-env

RUN apk add --no-cache git ca-certificates nodejs npm

ENV GOPROXY=https://proxy.golang.org,direct \
  GOSUMDB=off \
  GO111MODULE=on \
  PATH=/usr/local/lib/node_modules/.bin:$PATH

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download && go mod verify

RUN go install github.com/swaggo/swag/cmd/swag@latest && \
  npm install -g swagger2openapi

COPY . .

RUN swag init -g cmd/api/main.go --output docs && \
  swagger2openapi -o docs/v3/openapi.json docs/swagger.json && \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-s -w -extldflags '-static'" \
  -a -installsuffix cgo \
  -o app ./cmd/api

FROM scratch

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /src/app /app
COPY --from=build-env /src/docs /docs
COPY --from=build-env /src/migrations /migrations

USER 65534:65534

EXPOSE 8080

ENTRYPOINT ["/app"]
