FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

# Строим приложение
RUN go build -o todo-app ./cmd/app/main.go

FROM alpine:3.15

RUN apk add --no-cache postgresql-client

WORKDIR /app

COPY --from=builder /app/todo-app /app/

COPY wait-for-postgres.sh /app/

RUN chmod +x /app/wait-for-postgres.sh

CMD ["./todo-app"]
