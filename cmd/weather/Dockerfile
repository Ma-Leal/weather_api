FROM golang:1.24 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o weather ./cmd/weather

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=build /app/weather .
COPY cmd/weather/.env /app/.env

EXPOSE 8080

ENTRYPOINT [ "./weather" ]
