FROM golang:latest

WORKDIR /usr/app/
COPY . /usr/app/

#RUN go mod download && go get -u ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

ENV TZ Europe/Moscow

CMD ["./main"]
