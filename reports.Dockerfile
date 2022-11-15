FROM golang:1.19.1

WORKDIR /app

COPY . .

RUN cd reports && go build -o kservice

CMD ["./reports/kservice"]
