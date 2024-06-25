FROM golang:1.22.4-alpine as builder

WORKDIR /app
COPY ./src ./src
COPY main.go go.mod go.sum ./

RUN go mod tidy
RUN go build -o main .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/frontend ./frontend
COPY --from=builder /app/main .

CMD ["./main"]