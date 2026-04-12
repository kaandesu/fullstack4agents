FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api-server cmd/api/main.go

FROM alpine:3.20.1 AS api
WORKDIR /app
COPY --from=build /app/api-server /app/main
EXPOSE ${PORT}
EXPOSE 8090
CMD ["./main"]
