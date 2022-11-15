FROM golang:1.19.1

WORKDIR /src/telegram-bot

COPY cmd/ ./cmd
COPY data/ ./data
COPY internal/ ./internal
COPY go.sum .
COPY go.mod .

RUN go mod download

COPY Makefile .
RUN make build

CMD ["./bin/bot"]
