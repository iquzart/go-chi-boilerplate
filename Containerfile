
FROM golang:1.23-alpine as build-env

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app ./cmd/api

FROM gcr.io/distroless/base:nonroot

WORKDIR /app

COPY --from=build-env /src/app .

EXPOSE 8080

USER nonroot:nonroot

CMD ["./app"]

