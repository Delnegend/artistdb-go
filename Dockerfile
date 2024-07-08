FROM golang:alpine AS builder

WORKDIR /app
COPY ./src ./src
COPY main.go go.mod go.sum ./

RUN go mod tidy
RUN go build -o main .

FROM alpine:latest

WORKDIR /app
COPY ./frontend ./frontend
COPY --from=builder /app/main .

CMD ["./main"]