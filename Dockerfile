FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o todo-server main.go


FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/todo-server .
COPY --from=builder /app/web ./web

RUN mkdir -p /app/data

# Create env
ENV TODO_PORT=7540
ENV TODO_DBFILE=data/scheduler.db
ENV TODO_PASSWORD=12345
ENV TODO_SECRET=my_super_secret_key

EXPOSE 7540

CMD ["./todo-server"]
