FROM golang:1.19.3-alpine3.16 as builder

WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN go build -o raspisensor ./main.go

FROM alpine:3.16.3

WORKDIR /app

COPY --from=builder /app/raspisensor .

ENTRYPOINT /app/raspisensor
