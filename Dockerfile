from golang:1.22.2-alpine as builder

WORKDIR /app
COPY ./src ./src
COPY main.go go.mod go.sum .

RUN go mod tidy
RUN go build -o main .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/src/frontend ./src/frontend
COPY --from=builder /app/main .

CMD ["./main"]